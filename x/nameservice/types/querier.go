package types

import "strings"

type QueryResResolve struct {
	Value string 'json:"value"'
}

func (r QueryResResolve) String() string {
	return r.Value
}

type QueryResName []string

func(n QueryResName) String() string {
	return strings.Join(n[:], "\n")
}