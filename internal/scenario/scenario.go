package scenario

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

const (
	incorrectMessage         = "Некорректная команда"
	cardNotFoundMessage      = "Карта не найдена"
	pricesUnavailableMessage = "Цены временно недоступны, попробуйте позже"
)

type Scenario struct {
	Sender      Sender
	Logger      *log.Logger
	InfoFetcher CardInfoFetcher
	Cache       CardCache
}

type UserMessage struct {
	Body   string
	UserId int64
}

type CardCache interface {
	Init(ctx context.Context) error
	Get(ctx context.Context, cardName string) ([]cardsinfo.ScgCardPrice, error)
	Set(ctx context.Context, cardName string, prices []cardsinfo.ScgCardPrice) error
}

type CardInfoFetcher interface {
	GetNameByCardId(set string, number string) string
	GetOriginalName(name string) string
	GetPrices(name string) ([]cardsinfo.ScgCardPrice, error)
}

type Sender interface {
	Send(userId int64, message string)
	SendPrices(userId int64, cardName string, prices []cardsinfo.ScgCardPrice)
}

func (s *Scenario) HandleSearch(ctx context.Context, msg *UserMessage) {
	cardName, err := s.getCardNameByCommand(msg.Body)
	if err != nil {
		s.Sender.Send(msg.UserId, incorrectMessage)
		s.Logger.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), msg.Body)
	} else if cardName == "" {
		s.Sender.Send(msg.UserId, cardNotFoundMessage)
		s.Logger.Printf("[info] Could not find card. User input: %s", msg.Body)
	} else {
		prices, err := s.Cache.Get(ctx, cardName)
		if err == nil {
			s.Sender.SendPrices(msg.UserId, cardName, prices)
			return
		}
		prices, err = s.InfoFetcher.GetPrices(cardName)
		if err != nil {
			s.Sender.Send(msg.UserId, pricesUnavailableMessage)
			s.Logger.Printf("[error] Could not find SCG prices. Message: %s card name: %s", err.Error(), cardName)
			return
		}
		err = s.Cache.Set(ctx, cardName, prices)
		if err != nil {
			s.Logger.Println(fmt.Errorf("failed add entry in cache: %w", err))
		}
		s.Sender.SendPrices(msg.UserId, cardName, prices)
	}
}

func (s *Scenario) getCardNameByCommand(command string) (string, error) {
	var name string
	switch {
	case strings.HasPrefix(command, "!s"):
		split := strings.Split(command, " ")
		if len(split) < 3 {
			return "", errors.New("wrong command")
		}
		set := split[1]
		number := split[2]
		name = s.InfoFetcher.GetNameByCardId(set, number)
	default:
		name = s.InfoFetcher.GetOriginalName(command)
	}
	return name, nil
}
