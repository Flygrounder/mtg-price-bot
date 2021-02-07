package cardsinfo

import (
	"fmt"
	"strings"
)

type CardPrice interface {
	Format() string
}

type ScgCardPrice struct {
	Price   string
	Edition string
	Link    string
}

func (s *ScgCardPrice) Format() string {
	return fmt.Sprintf("%v: %v\n%v\n", s.Edition, s.Price, s.Link)
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
