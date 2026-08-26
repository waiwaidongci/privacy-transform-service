package txmetricsprobe

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/repository"
)

func TestTransactionCommitNilReceiver(t *testing.T) {
	var tx *repository.Transaction
	_ = tx.Commit()
}

func TestTransactionRollbackNilReceiver(t *testing.T) {
	var tx *repository.Transaction
	_ = tx.Rollback()
}

func TestTransactionAddUndoNilReceiver(t *testing.T) {
	var tx *repository.Transaction
	tx.AddUndo(func() {})
}

func TestMetricsIncNilReceiver(t *testing.T) {
	var m *metrics.Metrics
	m.Inc()
}

func TestMetricsHandlerNilReceiver(t *testing.T) {
	var m *metrics.Metrics
	rec := httptest.NewRecorder()
	m.Handler(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
}
