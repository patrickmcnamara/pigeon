package main

import (
	"strings"

	"github.com/patrickmcnamara/bird/seed"
)

func seedToText(sr *seed.Reader) string {
	var sb strings.Builder
	sw := seed.NewWriter(&sb)
	seed.Copy(sw, sr)
	return sb.String()
}
