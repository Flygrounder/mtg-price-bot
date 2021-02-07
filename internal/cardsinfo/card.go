package cardsinfo

import (
	"fmt"
	"strings"
)

type scgCardPrice struct {
	price   string
	edition string
	link    string
}

func (s *scgCardPrice) format() string {
	return fmt.Sprintf("%v: %v\n%v\n", s.edition, s.price, s.link)
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
