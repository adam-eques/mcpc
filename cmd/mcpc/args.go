package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseArgs turns CLI tokens into a tool-arguments map. Two forms are accepted:
//
//	key=value    string value
//	key:=value   raw JSON value (number, boolean, array, object, ...)
//
// For example: expression="2+2", or limit:=5, or tags:='["a","b"]'.
func parseArgs(tokens []string) (map[string]any, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	out := make(map[string]any, len(tokens))
	for _, tok := range tokens {
		if i := strings.Index(tok, ":="); i > 0 {
			key := tok[:i]
			var v any
			if err := json.Unmarshal([]byte(tok[i+2:]), &v); err != nil {
				return nil, fmt.Errorf("invalid JSON for %q: %w", key, err)
			}
			out[key] = v
			continue
		}
		if i := strings.Index(tok, "="); i > 0 {
			out[tok[:i]] = tok[i+1:]
			continue
		}
		return nil, fmt.Errorf("argument %q must be key=value or key:=json", tok)
	}
	return out, nil
}
