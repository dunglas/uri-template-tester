package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yosida95/uritemplate/v3"
)

type jsonError struct {
	Error string
}

func match(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "deny")

	query := r.URL.Query()
	template := query.Get("template")

	errors := make([]string, 0, 2)
	if template == "" {
		errors = append(errors, `The "template" parameter is mandatory.`)
	}

	uri := query.Get("uri")
	if uri == "" {
		errors = append(errors, `The "uri" parameter is mandatory.`)
	}

	if len(errors) > 0 {
		writeJSONError(w, r, http.StatusBadRequest, strings.Join(errors, " "))
		return
	}

	tpl, err := uritemplate.New(template)
	if nil != err {
		writeJSONError(w, r, http.StatusBadRequest, fmt.Sprintf(`"%s" is not a valid URI template (RFC6570).`, template))
		return
	}

	match := tpl.Match(uri)
	if match == nil {
		writeJSON(w, r, http.StatusOK, struct{ Match bool }{false})
		return
	}

	writeJSON(w, r, http.StatusOK, struct {
		Match  bool
		Values map[string]any
	}{true, flattenValues(match)})
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	payload, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		slog.Error("failed to marshal response", "remote_addr", r.RemoteAddr, "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		slog.Info("failed to write response", "remote_addr", r.RemoteAddr, "error", err)
	}
}

func writeJSONError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	payload, err := json.MarshalIndent(jsonError{msg}, "", "  ")
	if err != nil {
		slog.Error("failed to marshal error response", "remote_addr", r.RemoteAddr, "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		slog.Info("failed to write response", "remote_addr", r.RemoteAddr, "error", err)
	}
}

// flattenValues turns uritemplate/v3's Values into a JSON-friendly shape:
// single-string variables come out as strings, lists as []string, key/value
// arrays as map[string]string. Keeps the rendered output readable in the UI.
func flattenValues(v uritemplate.Values) map[string]any {
	out := make(map[string]any, len(v))
	for k, val := range v {
		switch val.T {
		case uritemplate.ValueTypeList:
			out[k] = val.List()
		case uritemplate.ValueTypeKV:
			kv := val.KV()
			m := make(map[string]string, len(kv)/2)
			for i := 0; i+1 < len(kv); i += 2 {
				m[kv[i]] = kv[i+1]
			}
			out[k] = m
		default:
			out[k] = val.String()
		}
	}
	return out
}
