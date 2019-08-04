// Package versionscmd contains command parsing utilities for the `versions` CLI
package versionscmd

import (
	"github.com/posener/complete"
)

// PredictSet1 accepts one of a specific set of terms, once
func PredictSet1(options ...string) complete.Predictor {
	m := make(map[string]struct{}, len(options))
	for _, o := range options {
		m[o] = struct{}{}
	}
	return predictSet1{m, options}
}

type predictSet1 struct {
	m map[string]struct{}
	s []string
}

func (p predictSet1) Predict(a complete.Args) []string {
	if _, ok := p.m[a.LastCompleted]; ok {
		return nil
	}
	return p.s
}
