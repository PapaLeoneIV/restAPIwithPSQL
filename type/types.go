package types

import()

type Album struct {
	Id    string  `json:"id"`
	Title  string  `json:"title"`
	Artist string `json:"artist"`
	Price int     `json:"price"`
}