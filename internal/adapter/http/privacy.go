// Package implementation for privacy transformation and sensitive-value protection.
package httpadapter

import (
	"encoding/json"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"net/http"
	"strings"
)

func (s *Server) classifications(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		writeJSON(w, 200, s.privacy.ListClassifications())
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var c domain.Classification
	if err := decode(r, &c); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.privacy.PutClassification(r.Context(), c)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, o)
}
func (s *Server) privacyPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var p domain.TransformPolicy
		if err := decode(r, &p); err != nil {
			writeErr(w, err)
			return
		}
		o, err := s.privacy.PutPolicy(r.Context(), p)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, 201, o)
		return
	}
	w.WriteHeader(405)
}
func (s *Server) publishPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		ID string `json:"id"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	p, err := s.privacy.PublishPolicy(in.ID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, p)
}
func (s *Server) rollbackPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		ID string `json:"id"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	p, _ := s.privacy.RollbackPolicy(in.ID)
	writeJSON(w, 200, p)
}
func (s *Server) simulatePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		PolicyID string         `json:"policy_id"`
		Data     map[string]any `json:"data"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.privacy.Simulate(r.Context(), in.PolicyID, in.Data)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) privacyResults(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.privacy.ListResults())
}
func (s *Server) transform(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		RequestID string         `json:"request_id"`
		PolicyID  string         `json:"policy_id"`
		Data      map[string]any `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, err)
		return
	}
	if in.RequestID == "" {
		in.RequestID = "req-transform"
	}
	out, err := s.privacy.Transform(r.Context(), in.RequestID, in.PolicyID, in.Data)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, out)
}
func (s *Server) transformBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		PolicyID string           `json:"policy_id"`
		Records  []map[string]any `json:"records"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	out := make([]domain.TransformResult, 0, len(in.Records))
	for i, d := range in.Records {
		res, err := s.privacy.Transform(r.Context(), strings.Join([]string{"batch", string(rune(i + 48))}, "-"), in.PolicyID, d)
		if err != nil {
			writeErr(w, err)
			return
		}
		out = append(out, res)
	}
	writeJSON(w, 200, map[string]any{"results": out, "count": len(out)})
}

func (s *Server) revokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct {
		Token string `json:"token"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.privacy.RevokeToken(in.Token); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}
