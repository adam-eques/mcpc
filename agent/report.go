package agent

// Result records the outcome of a single step.
type Result struct {
	Step   Step
	Output string
	IsErr  bool
	Err    error
}

// Report is the outcome of running a plan.
type Report struct {
	Results []Result
}

// Failed reports whether any step ended in error.
func (r Report) Failed() bool {
	for _, res := range r.Results {
		if res.Err != nil || res.IsErr {
			return true
		}
	}
	return false
}

// Outputs returns the textual output of each completed step in order.
func (r Report) Outputs() []string {
	out := make([]string, len(r.Results))
	for i, res := range r.Results {
		out[i] = res.Output
	}
	return out
}
