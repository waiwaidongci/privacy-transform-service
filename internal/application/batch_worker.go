package application

import (
	"context"

	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func (s *Service) runBatchEvaluations(ctx context.Context, sets []domain.TransformRuleSet, ec domain.EvaluationContext) ([]ValueResult, error) {
	results := make([]ValueResult, 0, len(sets))
	done := make(chan error, len(sets))
	for _, set := range sets {
		go func(id string) {
			result, err := s.Evaluate(ctx, id, ec)
			if err == nil {
				results = append(results, result)
			}
			done <- err
		}(set.ID)
	}
	for range sets {
		if err := <-done; err != nil {
			return nil, err
		}
	}
	return results, nil
}
