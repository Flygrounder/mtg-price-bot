package cardsinfo

import (
	"strings"
)

type CardPrice struct {
	Name    string
	Price   string
	PriceFoil string
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

type ScgResponse struct {
	Response ScgResponseContainer `json:"response"`
}

type ScgResponseContainer struct {
	Data []ScgConditionContainer `json:"data"`
}

type ScgConditionContainer struct {
	Price        float64        `json:"price"`
	OptionValues []ScgCondition `json:"option_values"`
}

type ScgCondition struct {
	Label string `json:"label"`
}
