package main

import "github.com/patrickmcnamara/bird/seed"

type birdResponse struct {
	sr    *seed.Reader
	close func() (err error)
}
