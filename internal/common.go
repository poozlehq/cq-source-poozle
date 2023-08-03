package internal

type Meta struct {
	Limit   int     `json:"limit"`
	Cursors Cursors `json:"cursors"`
}

type Cursors struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}
