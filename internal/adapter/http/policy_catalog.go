package httpadapter

import (
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"net/http"
)

func (s *Server) privacyWorkspaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.PolicyWorkspace
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.policy.CreatePolicyWorkspace(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) processingPurposes(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.policy.ListProcessingPurposes(r.Context(), r.URL.Query().Get("workspace_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.ProcessingPurpose
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.policy.CreateProcessingPurpose(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) transformRuleSets(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.policy.ListTransformRuleSets(r.Context(), r.URL.Query().Get("workspace_id"), r.URL.Query().Get("purpose_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.TransformRuleSet
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.policy.CreateTransformRuleSet(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
