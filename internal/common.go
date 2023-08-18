package internal

import "time"

type Meta struct {
	Limit   int     `json:"limit"`
	Cursors Cursors `json:"cursors"`
}

type Cursors struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
	Meta  *Meta       `json:"meta"`
}

type APIRateLimitError struct {
	Name     string
	TryAfter string
}

type RateLimitError struct {
	RetryAfter time.Duration
}
