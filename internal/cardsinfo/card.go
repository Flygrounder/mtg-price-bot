package cardsinfo

import (
	"strings"
)

type ScgCardPrice struct {
	Price   string
	Edition string
	Link    string
}

type card struct {
	Name   string `json:"name"`
	Layout string `json:"layout"`
}

func (c *card) getName() string {
	if c.Layout == "transform" {
		return strings.Replace(c.Name, "//", "|", 1)
	}
	return c.Name
}
