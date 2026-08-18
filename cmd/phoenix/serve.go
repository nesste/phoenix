package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/activate"
	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/surface"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
	devrepo "github.com/nesste/phoenix/verbs/dev-repo"
)

const externalOutputLimit = 256 * 1024

type serveOptions struct {
	arm             string
	serverVersion   string
	worldPath       string
	schemaPath      string
	episodePath     string
	rootRefsPath    string
	stateEventsPath string
	worldBuild      string
	warning         io.Writer
}

type servePaths struct {
	world    string
	schema   string
	episodes string
}

type assembledSurface struct {
	server *mcp.Server
	store  *episode.Store
}

func (assembled *assembledSurface) Close() error {
	if assembled == nil || assembled.store == nil {
		return nil
	}
	return assembled.store.Close()
}

func assembleSurface(options serveOptions) (*assembledSurface, error) {
	suppressFrontier, suppressTeaching, flat, err := armPolicy(options.arm)
	if err != nil {
		return nil, err
	}
	definition, err := world.Load(options.schemaPath, options.worldPath)
	if err != nil {
		return nil, err
	}
	rootRefs, err := loadRootReferences(options.rootRefsPath)
	if err != nil {
		return nil, err
	}
	afterAct, err := loadStateEvents(options.stateEventsPath)
	if err != nil {
		return nil, err
	}
	graph, err := configuredGraph(definition, repositoryResolver{}, rootRefs)
	if err != nil {
		return nil, err
	}
	commandRunner, err := externalRunner()
	if err != nil {
		return nil, err
	}
	store := openEpisodeStore(options.episodePath, options.warning)
	var episodeLog surface.EpisodeLog
	var episodeRecaller devrepo.EpisodeRecaller
	if store != nil {
		episodeLog = store
		episodeRecaller = store
	}
	build := options.worldBuild
	if build == "" {
		build = definition.Digest()
	}
	admission, err := configuredAdmission(admissionConfig{
		definition: definition, graph: graph, commandRunner: commandRunner,
		episodeLog: episodeLog, episodeRecaller: episodeRecaller, warning: options.warning,
		afterAct: afterAct, build: build, suppressFrontier: suppressFrontier, suppressTeaching: suppressTeaching,
	})
	if err != nil {
		return nil, closeOnError(store, err)
	}
	if flat {
		adapter, flatErr := surface.NewFlatMCP(options.serverVersion, admission, definition)
		if flatErr != nil {
			return nil, closeOnError(store, flatErr)
		}
		return &assembledSurface{server: adapter.Server(), store: store}, nil
	}
	return &assembledSurface{server: surface.NewMCP(options.serverVersion, admission).Server(), store: store}, nil
}

type admissionConfig struct {
	definition       *world.Definition
	graph            *world.Graph
	commandRunner    verb.CommandRunner
	episodeLog       surface.EpisodeLog
	episodeRecaller  devrepo.EpisodeRecaller
	warning          io.Writer
	afterAct         func(context.Context, int, surface.Input, surface.Envelope) error
	build            string
	suppressFrontier bool
	suppressTeaching bool
}

func configuredAdmission(config admissionConfig) (*surface.Admission, error) {
	registry := verb.NewRegistry()
	if err := devrepo.Register(registry, devrepo.Config{
		Runner: config.commandRunner, GoExecutable: "go", GitExecutable: "git", Episodes: config.episodeRecaller,
	}); err != nil {
		return nil, fmt.Errorf("register dev-repo verbs: %w", err)
	}
	frontierEngine, err := frontier.New(config.definition)
	if err != nil {
		return nil, err
	}
	activationEngine, err := activate.New(config.definition)
	if err != nil {
		return nil, err
	}
	teachingEngine, err := teach.New(config.definition)
	if err != nil {
		return nil, err
	}
	return surface.New(surface.Config{
		WorldBuild: config.build, Graph: config.graph,
		Executor: verb.NewExecutor(registry, verb.Options{}),
		Frontier: frontierEngine, Activation: activationEngine, Teacher: teachingEngine,
		Episodes: config.episodeLog, Warning: config.warning, AfterAct: config.afterAct,
		SuppressFrontier: config.suppressFrontier, SuppressTeaching: config.suppressTeaching,
	})
}

func armPolicy(arm string) (suppressFrontier, suppressTeaching, flat bool, err error) {
	switch arm {
	case "", "C":
		return false, false, false, nil
	case "D":
		return true, false, false, nil
	case "E":
		return true, true, false, nil
	case "A", "B":
		return true, true, true, nil
	default:
		return false, false, false, fmt.Errorf("unsupported experiment arm %q", arm)
	}
}

func configuredGraph(definition *world.Definition, resolver world.Resolver, roots map[string]string) (*world.Graph, error) {
	if roots == nil {
		return world.NewGraph(definition, resolver), nil
	}
	return world.NewGraphWithRootReferences(definition, resolver, roots)
}

func openEpisodeStore(path string, warning io.Writer) *episode.Store {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		fmt.Fprintf(warning, "phoenix warning: episode logging disabled: %v\n", err)
		return nil
	}
	store, err := episode.Open(path, episode.Options{})
	if err != nil {
		fmt.Fprintf(warning, "phoenix warning: episode logging disabled: %v\n", err)
		return nil
	}
	return store
}

func closeOnError(store *episode.Store, err error) error {
	if store != nil {
		_ = store.Close()
	}
	return err
}

func loadRootReferences(path string) (map[string]string, error) {
	if path == "" {
		return nil, nil
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read root references: %w", err)
	}
	var references map[string]string
	decoder := json.NewDecoder(strings.NewReader(string(contents)))
	if err := decoder.Decode(&references); err != nil {
		return nil, fmt.Errorf("decode root references: %w", err)
	}
	return references, nil
}

func externalRunner() (*verb.ExternalRunner, error) {
	executables := make([]verb.Executable, 0, 2)
	for _, name := range []string{"go", "git"} {
		path, err := exec.LookPath(name)
		if err != nil {
			return nil, fmt.Errorf("locate %s executable: %w", name, err)
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("resolve %s executable: %w", name, err)
		}
		digest, err := verb.DigestExecutable(path)
		if err != nil {
			return nil, fmt.Errorf("digest %s executable: %w", name, err)
		}
		executables = append(executables, verb.Executable{ID: name, Path: path, Digest: digest})
	}
	return verb.NewExternalRunner(executables, externalOutputLimit)
}

func defaultServePaths() servePaths {
	root := findRepositoryRoot()
	cache, err := os.UserCacheDir()
	if err != nil || cache == "" {
		cache = os.TempDir()
	}
	return servePaths{
		world:    filepath.Join(root, "worlds", "dev-repo", "world.json"),
		schema:   filepath.Join(root, "spec", "world.schema.json"),
		episodes: filepath.Join(cache, "phoenix", "episodes.db"),
	}
}

func findRepositoryRoot() string {
	current, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "spec", "world.schema.json")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "."
		}
		current = parent
	}
}

type repositoryResolver struct{}

func (repositoryResolver) Resolve(_ context.Context, resource world.Resource) (any, error) {
	if resource.Kind != "path" {
		return map[string]any{"kind": resource.Kind, "value": resource.Value}, nil
	}
	root, err := filepath.Abs(resource.Value)
	if err != nil {
		return nil, err
	}
	entries := []any{}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == ".git" {
			return fs.SkipDir
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(contents)
		entries = append(entries, map[string]any{
			"path": filepath.ToSlash(relative), "sha256": hex.EncodeToString(sum[:]),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("resolve repository state: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].(map[string]any)["path"].(string) < entries[j].(map[string]any)["path"].(string)
	})
	return map[string]any{"files": entries}, nil
}
