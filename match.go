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
	w.Header().Set("Content-Type", "application/payload")
	w.Header().Set("X-Content-Type-Options", "nosniff")
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
		payload, _ := json.MarshalIndent(jsonError{strings.Join(errors, " ")}, "", "  ")
		http.Error(w, string(payload), http.StatusBadRequest)

		return
	}

	tpl, err := uritemplate.New(template)
	if nil != err {
		payload, _ := json.MarshalIndent(jsonError{fmt.Sprintf(`"%s" is not a valid URI template (RFC6570).`, template)}, "", "  ")
		http.Error(w, string(payload), http.StatusBadRequest)

		return
	}

	match := tpl.Match(uri)
	if match == nil {
		payload, _ := json.Marshal(struct{ Match bool }{false})
		if _, err := w.Write(payload); err != nil {
			slog.Info("failed to write response", "remote_addr", r.RemoteAddr, "error", err)
		}

		return
	}

	payload, _ := json.MarshalIndent(struct {
		Match  bool
		Values map[string]any
	}{true, flattenValues(match)}, "", "  ")
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
