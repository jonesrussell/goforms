package mocksubscription

import (
	"context"
)

// mockCall represents a single method call
type mockCall struct {
	method string
	args   []interface{}
	ret    []interface{}
}

// matchArgs compares two argument slices
func matchArgs(exp, got []interface{}) bool {
	if len(exp) != len(got) {
		return false
	}
	for i := range exp {
		// For context, just check if both are contexts
		if _, expIsCtx := exp[i].(context.Context); expIsCtx {
			_, gotIsCtx := got[i].(context.Context)
			return gotIsCtx
		}
		if exp[i] != got[i] {
			return false
		}
	}
	return true
}
