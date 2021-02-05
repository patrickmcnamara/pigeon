package main

import (
	"strings"

	"github.com/patrickmcnamara/bird/seed"
)

func seedToText(sr *seed.Reader) string {
	var sb strings.Builder
	sw := seed.NewWriter(&sb)
	for {
		line, err := sr.ReadLine()
		if err != nil {
			break
		}
		sw.Line(line)
	}
	return sb.String()
}
