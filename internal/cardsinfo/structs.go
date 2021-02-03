package cardsinfo

import (
	"fmt"
	"strings"
)

type CardPrice interface {
	Format() string
}

type TcgCardPrice struct {
	FullArt   bool
	Name      string
	Price     string
	PriceFoil string
	Link      string
	Edition   string
}

func (t *TcgCardPrice) Format() string {
	return fmt.Sprintf("%v\nRegular: %v\nFoil: %v\n%v\n", t.Edition, formatTcgPrice(t.Price), formatTcgPrice(t.PriceFoil), t.Link)
}

func formatTcgPrice(price string) string {
	if price == "" {
		return "-"
	}
	return fmt.Sprintf("$%v", price)
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
