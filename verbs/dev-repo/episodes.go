package devrepo

import (
	"context"

	"github.com/nesste/phoenix/internal/verb"
)

func recallHandler(recaller EpisodeRecaller) verb.HandlerFunc {
	return func(ctx context.Context, request verb.Request) (any, error) {
		if recaller == nil {
			return nil, verb.NewFailure("recall_unavailable", "episode recall is not available before the episode store starts", nil)
		}
		pointers, err := recaller.Recall(ctx, request.SessionID, request.Args["query"].(string))
		if err != nil {
			return nil, verb.NewFailure("recall_failed", "episode recall failed", nil)
		}
		episodes := make([]any, 0, len(pointers))
		for _, pointer := range pointers {
			episodes = append(episodes, map[string]any{"episode_id": pointer.EpisodeID, "act_id": pointer.ActID})
		}
		return map[string]any{"episodes": episodes}, nil
	}
}
