package devrepo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

func readHandler(maxBytes int64) verb.HandlerFunc {
	return func(_ context.Context, request verb.Request) (any, error) {
		path, relative, err := safeExistingFile(request.Resource, request.Args["path"].(string))
		if err != nil {
			return nil, err
		}
		contents, err := readBounded(path, maxBytes)
		if err != nil {
			return nil, err
		}
		return map[string]any{"path": relative, "content": string(contents), "digest": digest(contents)}, nil
	}
}

func editHandler(maxBytes int64) verb.HandlerFunc {
	return func(_ context.Context, request verb.Request) (any, error) {
		path, relative, err := safeExistingFile(request.Resource, request.Args["path"].(string))
		if err != nil {
			return nil, err
		}
		contents, err := readBounded(path, maxBytes)
		if err != nil {
			return nil, err
		}
		oldText := request.Args["old"].(string)
		if strings.Count(string(contents), oldText) != 1 {
			return nil, verb.NewFailure("edit_conflict", "edit precondition did not match exactly once", nil)
		}
		updated := []byte(strings.Replace(string(contents), oldText, request.Args["new"].(string), 1))
		if int64(len(updated)) > maxBytes {
			return nil, verb.NewFailure("output_limit", "edited file exceeds its size limit", map[string]any{"limit_bytes": maxBytes})
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, verb.NewFailure("edit_failed", "file metadata could not be read", nil)
		}
		if err := os.WriteFile(path, updated, info.Mode().Perm()); err != nil {
			return nil, verb.NewFailure("edit_failed", "file could not be written", nil)
		}
		return map[string]any{"path": relative, "changed": true, "digest": digest(updated)}, nil
	}
}

func findHandler(config Config) verb.HandlerFunc {
	return func(_ context.Context, request verb.Request) (any, error) {
		root, err := repoRoot(request.Resource)
		if err != nil {
			return nil, err
		}
		query := strings.ToLower(request.Args["query"].(string))
		matches := reachableMatches(request.Reachable, query, config.MaxSearchMatches)
		if len(matches) < config.MaxSearchMatches {
			fileMatches, err := searchFiles(root, query, config.MaxReadBytes, config.MaxSearchMatches-len(matches))
			if err != nil {
				return nil, err
			}
			matches = append(matches, fileMatches...)
		}
		return map[string]any{"matches": stringsAsAny(matches)}, nil
	}
}

func reachableMatches(finder verb.ReachableFinder, query string, limit int) []string {
	if finder == nil {
		return []string{}
	}
	matches := finder.FindReachable(query)
	if len(matches) > limit {
		matches = matches[:limit]
	}
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, fmt.Sprintf("%s %s.%s — %s", match.Handle, match.Type, match.Verb, match.Label))
	}
	return result
}

func searchFiles(root, query string, maxBytes int64, limit int) ([]string, error) {
	matches := make([]string, 0, limit)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if len(matches) >= limit {
			return fs.SkipAll
		}
		if entry.IsDir() {
			if entry.Name() == ".git" && path != root {
				return fs.SkipDir
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		relative, _ := filepath.Rel(root, path)
		contents, err := readBounded(path, maxBytes)
		if err != nil {
			return nil
		}
		if strings.Contains(strings.ToLower(filepath.ToSlash(relative)+"\n"+string(contents)), query) {
			matches = append(matches, filepath.ToSlash(relative))
		}
		return nil
	})
	if err != nil {
		return nil, verb.NewFailure("find_failed", "repository search failed", nil)
	}
	return matches, nil
}

func safeExistingFile(resource world.Resource, requested string) (string, string, error) {
	root, err := repoRoot(resource)
	if err != nil {
		return "", "", err
	}
	if filepath.IsAbs(requested) {
		return "", "", verb.NewFailure("path_outside_repo", "path is outside the repository", nil)
	}
	candidate := filepath.Join(root, filepath.Clean(filepath.FromSlash(requested)))
	relative, err := filepath.Rel(root, candidate)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", "", verb.NewFailure("path_outside_repo", "path is outside the repository", nil)
	}
	info, err := os.Lstat(candidate)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return "", "", verb.NewFailure("file_unavailable", "requested file is not a regular repository file", nil)
	}
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", "", verb.NewFailure("file_unavailable", "requested file could not be resolved", nil)
	}
	resolvedRelative, err := filepath.Rel(root, resolved)
	if err != nil || resolvedRelative == ".." || strings.HasPrefix(resolvedRelative, ".."+string(filepath.Separator)) {
		return "", "", verb.NewFailure("path_outside_repo", "path is outside the repository", nil)
	}
	return resolved, filepath.ToSlash(relative), nil
}

func repoRoot(resource world.Resource) (string, error) {
	if resource.Kind != "path" || resource.Value == "" {
		return "", verb.NewFailure("resource_invalid", "verb requires a repository path resource", nil)
	}
	root, err := filepath.Abs(resource.Value)
	if err != nil {
		return "", verb.NewFailure("resource_invalid", "repository path could not be resolved", nil)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", verb.NewFailure("resource_invalid", "repository path could not be resolved", nil)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return "", verb.NewFailure("resource_invalid", "repository path is not a directory", nil)
	}
	return root, nil
}

func readBounded(path string, maxBytes int64) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, verb.NewFailure("file_unavailable", "file metadata could not be read", nil)
	}
	if info.Size() > maxBytes {
		return nil, verb.NewFailure("output_limit", "file exceeds its read limit", map[string]any{"limit_bytes": maxBytes})
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, verb.NewFailure("read_failed", "file could not be read", nil)
	}
	return contents, nil
}

func digest(contents []byte) string {
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func stringsAsAny(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
