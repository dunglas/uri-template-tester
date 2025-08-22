package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yosida95/uritemplate"
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
		Values uritemplate.Values
	}{true, match}, "", "  ")
	if _, err := w.Write(payload); err != nil {
		slog.Info("failed to write response", "remote_addr", r.RemoteAddr, "error", err)
	}
}
