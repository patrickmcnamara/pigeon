package main

import (
	"fmt"
	"strings"

	"github.com/patrickmcnamara/bird/seed"
)

func seedToText(sr *seed.Reader) string {
	quote := false
	code := false
	var sb strings.Builder
	for {
		line, err := sr.ReadLine()
		if err != nil {
			break
		}
		switch l := line.(type) {
		case seed.Text:
			switch {
			case quote:
				sb.WriteString(fmt.Sprintln(">", string(l)))
			case code:
				sb.WriteString(fmt.Sprintln("`", string(l)))
			default:
				sb.WriteString(fmt.Sprintln(string(l)))
			}
		case seed.Header:
			sb.WriteString(fmt.Sprintln(strings.Repeat("#", l.Level), l.Text))
		case seed.Link:
			sb.WriteString(fmt.Sprintf("=> %s (%s)\n", l.Text, l.URL))
		case seed.Quote:
			quote = !quote
		case seed.Code:
			code = !code
		case seed.Break:
			sb.WriteString(fmt.Sprintln())
		}
	}
	return sb.String()
}
