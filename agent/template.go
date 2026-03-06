package agent

import (
	"strconv"
	"strings"
)

// substitute replaces ${N} placeholders in string arguments with the output of a
// previously completed step (1-based). This lets a plan feed one tool's result
// into the next, the way an agent chains tool calls. Unknown references are left
// untouched.
func substitute(args map[string]any, outputs []string) map[string]any {
	if len(args) == 0 || len(outputs) == 0 {
		return args
	}
	out := make(map[string]any, len(args))
	for k, v := range args {
		if s, ok := v.(string); ok {
			out[k] = expand(s, outputs)
		} else {
			out[k] = v
		}
	}
	return out
}

// expand replaces every ${N} in s with outputs[N-1] when N is in range.
func expand(s string, outputs []string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '$' && i+1 < len(s) && s[i+1] == '{' {
			end := strings.IndexByte(s[i:], '}')
			if end > 0 {
				ref := s[i+2 : i+end]
				if n, err := strconv.Atoi(ref); err == nil && n >= 1 && n <= len(outputs) {
					b.WriteString(outputs[n-1])
					i += end + 1
					continue
				}
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
