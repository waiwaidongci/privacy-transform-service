package cacheprobe

import (
	"context"
	"testing"
	"time"

	"github.com/ali/go-0821/privacy-transform-service/internal/adapter/cache"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

type revisionCache interface {
	Set(context.Context, string, domain.TransformRevision)
	Get(context.Context, string) (*domain.TransformRevision, bool)
}

func revision() domain.TransformRevision {
	pct := 25
	return domain.TransformRevision{Number: 1, Value: map[string]any{"scope": []any{"email"}}, Rules: []domain.Rule{{ID: "r", Tags: map[string]string{"region": "jp"}, Percentage: &pct, Value: map[string]any{"mask": "full"}}}}
}

func assertInputDetached(t *testing.T, c revisionCache) {
	t.Helper()
	v := revision()
	c.Set(context.Background(), "p", v)
	v.Rules[0].Tags["region"] = "us"
	*v.Rules[0].Percentage = 90
	v.Value.(map[string]any)["scope"].([]any)[0] = "phone"
	got, _ := c.Get(context.Background(), "p")
	if got.Rules[0].Tags["region"] != "jp" || *got.Rules[0].Percentage != 25 || got.Value.(map[string]any)["scope"].([]any)[0] != "email" {
		t.Fatalf("cache retained caller aliases: %#v", got)
	}
}

func assertOutputDetached(t *testing.T, c revisionCache) {
	t.Helper()
	c.Set(context.Background(), "p", revision())
	got, _ := c.Get(context.Background(), "p")
	got.Rules[0].Tags["region"] = "us"
	*got.Rules[0].Percentage = 90
	got.Value.(map[string]any)["scope"].([]any)[0] = "phone"
	again, _ := c.Get(context.Background(), "p")
	if again.Rules[0].Tags["region"] != "jp" || *again.Rules[0].Percentage != 25 || again.Value.(map[string]any)["scope"].([]any)[0] != "email" {
		t.Fatalf("returned value mutated cache state: %#v", again)
	}
}

func TestMemoryCacheDetachesInputRevision(t *testing.T) { assertInputDetached(t, cache.NewMemory()) }
func TestMemoryCacheDetachesReturnedRevision(t *testing.T) { assertOutputDetached(t, cache.NewMemory()) }
func TestTTLCacheDetachesInputRevision(t *testing.T) { assertInputDetached(t, cache.NewTTL(time.Minute)) }
func TestTTLCacheDetachesReturnedRevision(t *testing.T) { assertOutputDetached(t, cache.NewTTL(time.Minute)) }
