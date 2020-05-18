package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yosida95/uritemplate"
)

func match(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "deny")

	query := r.URL.Query()
	template := query.Get("template")

	errors := make([]string, 0, 2)
	if template == "" {
		errors = append(errors, "The \"template\" parameter is mandatory.")
	}

	uri := query.Get("uri")
	if uri == "" {
		errors = append(errors, "The \"uri\" parameter is mandatory.")
	}

	if len(errors) > 0 {
		err, _ := json.MarshalIndent(error{strings.Join(errors, " ")}, "", "  ")
		http.Error(w, string(err), http.StatusBadRequest)
		return
	}

	tpl, err := uritemplate.New(template)
	if nil != err {
		err, _ := json.MarshalIndent(error{fmt.Sprintf("\"%s\" is not a valid URI template (RFC6570).", template)}, "", "  ")
		http.Error(w, string(err), http.StatusBadRequest)
		return
	}

	match := tpl.Match(uri)
	if match == nil {
		err, _ := json.Marshal(struct{ Match bool }{false})
		fmt.Fprintf(w, "%s", string(err))
		return
	}

	json, _ := json.MarshalIndent(struct {
		Match  bool
		Values uritemplate.Values
	}{true, match}, "", "  ")
	if _, err := w.Write(json); err != nil {
		panic(err)
	}
}
