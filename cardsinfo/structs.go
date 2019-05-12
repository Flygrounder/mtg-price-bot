package cardsinfo

import (
	"strings"
)

type CardPrice struct {
	Name    string
	Price   float64
	Link    string
	Edition string
}

type Card struct {
	Name   string `json:"name"`
	Layout string `json:"layout"`
}

func (c *Card) getName() string {
	if c.Layout == "transform" {
		return strings.Replace(c.Name, "//", "|", 1)
	}
	return c.Name
}
