package scenario

import (
	"errors"
	"log"
	"strings"
)

const (
	incorrectMessage         = "Некорректная команда"
	cardNotFoundMessage      = "Карта не найдена"
	pricesUnavailableMessage = "Цены временно недоступны, попробуйте позже"
)

type Scenario struct {
    Sender Sender
    Logger             *log.Logger
    InfoFetcher        CardInfoFetcher
    Cache              CardCache
}

type UserMessage struct {
	Body   string
	UserId int64
}

type CardCache interface {
	Get(cardName string) (string, error)
	Set(cardName string, message string)
}

type CardInfoFetcher interface {
	GetFormattedCardPrices(name string) (string, error)
	GetNameByCardId(set string, number string) string
	GetOriginalName(name string) string
}

type Sender interface {
	Send(userId int64, message string)
}

func (s *Scenario) HandleSearch(msg *UserMessage) {
	cardName, err := s.getCardNameByCommand(msg.Body)
	if err != nil {
		s.Sender.Send(msg.UserId, incorrectMessage)
		s.Logger.Printf("[info] Not correct command. Message: %s user input: %s", err.Error(), msg.Body)
	} else if cardName == "" {
		s.Sender.Send(msg.UserId, cardNotFoundMessage)
		s.Logger.Printf("[info] Could not find card. User input: %s", msg.Body)
	} else {
		message, err := s.getMessage(cardName)
		if err != nil {
			s.Sender.Send(msg.UserId, pricesUnavailableMessage)
			s.Logger.Printf("[error] Could not find SCG prices. Message: %s card name: %s", err.Error(), cardName)
			return
		}
		s.Sender.Send(msg.UserId, message)
	}
}

func (s *Scenario) getMessage(cardName string) (string, error) {
	val, err := s.Cache.Get(cardName)
	if err != nil {
		message, err := s.InfoFetcher.GetFormattedCardPrices(cardName)
		if err != nil {
			return "", err
		}
		s.Cache.Set(cardName, message)
		return message, nil
	}
	return val, nil
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
