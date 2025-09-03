package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	logger := log.Default()
	handler := requestHandler{
		logger: logger,
		client: http.DefaultClient,
	}
	err := run(&handler)
	if err != nil {
		logger.Fatalf("failed with error: %v", err.Error())
	}
}

func run(handler *requestHandler) error {
	err := handler.setTelegramWebhook()
	if err != nil {
		return fmt.Errorf("failed to set telegram webhook: %w", err)
	}
	handler.logger.Printf("successfully set telegram webhook")
	http.HandleFunc("/tg", handler.handleTelegramRequest)
	http.HandleFunc("/vk", handler.handleVKRequest)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		return fmt.Errorf("server stopped with error: %w", err)
	}
	return nil
}

type requestHandler struct {
	logger *log.Logger
	client *http.Client
}

func (r *requestHandler) setTelegramWebhook() error {
	tgURL := fmt.Sprintf("https://api.telegram.org/bot%v/setWebhook", os.Getenv("TG_TOKEN"))
	resp, err := r.client.PostForm(tgURL, url.Values{
		"url":          []string{os.Getenv("TG_WEBHOOK_URL")},
		"secret_token": []string{os.Getenv("TG_SECRET")},
	})
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram webhook url returned non-ok status: %v", resp.StatusCode)
	}
	return nil
}

func (r *requestHandler) handleTelegramRequest(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Telegram-Bot-Api-Secret-Token") != os.Getenv("TG_SECRET") {
		r.logger.Printf("got request with incorrect secret")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		r.logger.Printf("error reading request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upd := struct {
		Message *struct {
			Text *string `json:"text"`
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	}{}
	err = json.Unmarshal(res, &upd)
	if err != nil {
		r.logger.Printf("error unmarshalling request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if upd.Message == nil || upd.Message.Text == nil {
		return
	}

	msg, chatID := *upd.Message.Text, upd.Message.Chat.ID

	original, err := r.findOriginalName(msg)
	if err != nil {
		err = r.sendTelegramMessage(chatID, "Карта не найдена")
		if err != nil {
			r.logger.Printf("failed to send the response: %v", err.Error())
		}
		return
	}

	cards, _ := r.findCardPrices(original)

	cards = cards[:min(len(cards), 5)]

	text := formatPricesForTelegram(original, cards)

	err = r.sendTelegramMessage(chatID, text)
	if err != nil {
		r.logger.Printf("failed to send the response: %v", err.Error())
	}
}

func (r *requestHandler) handleVKRequest(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		r.logger.Printf("failed to read vk request body: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	event := struct {
		Type   string `json:"type"`
		Object struct {
			FromID int64  `json:"from_id"`
			Text   string `json:"text"`
		} `json:"object"`
		Secret string `json:"secret"`
	}{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		r.logger.Printf("failed to unmarshal vk request body: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Secret != os.Getenv("VK_SECRET") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "confirmation" {
		_, _ = w.Write([]byte(os.Getenv("VK_CONFIRMATION")))
		return
	}

	if event.Type != "message_new" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() { _, _ = w.Write([]byte("ok")) }()

	original, err := r.findOriginalName(event.Object.Text)
	if err != nil {
		err = r.sendVKMessage(event.Object.FromID, "Карта не найдена")
		if err != nil {
			r.logger.Printf("failed to send vk message: %v", err.Error())
		}
		return
	}

	cards, _ := r.findCardPrices(original)

	cards = cards[:min(len(cards), 5)]

	err = r.sendVKMessage(event.Object.FromID, formatPricesForVK(original, cards))
	if err != nil {
		r.logger.Printf("failed to send vk message: %v", err.Error())
	}
}

func (r *requestHandler) findOriginalName(query string) (string, error) {
	parts := strings.Split(query, " ")
	var addr string
	if len(parts) == 3 && parts[0] == "!s" {
		set, num := parts[1], parts[2]
		addr = fmt.Sprintf("https://api.scryfall.com/cards/%v/%v", url.PathEscape(strings.ToLower(set)), url.PathEscape(num))
	} else {
		esc := url.QueryEscape(query)
		addr = fmt.Sprintf("https://api.scryfall.com/cards/named?fuzzy=%v", esc)
	}
	resp, err := r.client.Get(addr)
	if err != nil {
		err := fmt.Errorf("failed to send request to scryfall: %w", err)
		r.logger.Println(err.Error())
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("scryfall returned non-ok status code: %v", resp.StatusCode)
		if resp.StatusCode != http.StatusNotFound {
			r.logger.Println(err.Error())
		}
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("failed to read scryfall response: %w", err)
		r.logger.Println(err.Error())
		return "", err
	}

	cardInfo := struct {
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(body, &cardInfo)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal scryfall response: %w", err)
		r.logger.Println(err.Error())
		return "", err
	}

	return cardInfo.Name, nil
}

type pricedCard struct {
	link  string
	set   string
	price string
}

func (r *requestHandler) findCardPrices(name string) ([]pricedCard, error) {
	req := map[string]any{
		"Keyword":    name,
		"clientguid": os.Getenv("SCG_CLIENT_GUID"),
	}
	bodyRaw, err := json.Marshal(req)
	if err != nil {
		err := fmt.Errorf("failed to marshal scg request: %w", err)
		r.logger.Println(err.Error())
		return nil, err
	}
	reqBody := strings.NewReader(string(bodyRaw))
	resp, err := r.client.Post("https://essearchapi-na.hawksearch.com/api/v2/search", "application/json", reqBody)
	if err != nil {
		err := fmt.Errorf("failed to send request to starcitygames: %w", err)
		r.logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("scg returned non-ok status code: %v", resp.StatusCode)
		r.logger.Println(err.Error())
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("failed to get response from scg: %w", err)
		r.logger.Println(err.Error())
		return nil, err
	}

	priceData := struct {
		Results []struct {
			Document struct {
				PriceRetail      []string `json:"price_retail"`
				Set              []string `json:"set"`
				ChildInformation []string `json:"child_information"`
			} `json:"Document"`
		} `json:"Results"`
	}{}
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal scg response: %w", err)
		r.logger.Println(err.Error())
		return nil, err
	}

	var prices []pricedCard
	for _, res := range priceData.Results {
		doc := res.Document
		if len(doc.Set) == 0 || len(doc.PriceRetail) == 0 || len(doc.ChildInformation) == 0 {
			continue
		}
		childInfo := struct {
			URL []string `json:"url"`
		}{}
		err := json.Unmarshal([]byte(doc.ChildInformation[0]), &childInfo)
		if err != nil || len(childInfo.URL) == 0 {
			continue
		}
		prices = append(prices, pricedCard{
			set:   doc.Set[0],
			price: doc.PriceRetail[0],
			link:  childInfo.URL[0],
		})
	}

	return prices, nil
}

func (r *requestHandler) sendTelegramMessage(chatID int64, message string) error {
	tgURL := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", os.Getenv("TG_TOKEN"))

	type LinkPreviwOptions struct {
		IsDisabled bool `json:"is_disabled"`
	}

	req := struct {
		ChatID             int64             `json:"chat_id"`
		Text               string            `json:"text"`
		ParseMode          string            `json:"parse_mode"`
		LinkPreviewOptions LinkPreviwOptions `json:"link_preview_options"`
	}{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "HTML",
		LinkPreviewOptions: LinkPreviwOptions{
			IsDisabled: true,
		},
	}
	bodyJson, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram request: %w", err)
	}
	body := strings.NewReader(string(bodyJson))
	resp, err := r.client.Post(tgURL, "application/json", body)
	if err != nil {
		return fmt.Errorf("failed to send request to telegram: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram responded with non-ok status code: %v", resp.StatusCode)
	}
	return nil
}

func (r *requestHandler) sendVKMessage(chatID int64, message string) error {
	resp, err := r.client.Get(fmt.Sprintf("https://api.vk.ru/method/messages.send?message=%v&peer_id=%v&access_token=%v&v=5.103&random_id=0", url.QueryEscape(message), chatID, url.QueryEscape(os.Getenv("VK_TOKEN"))))
	if err != nil {
		return fmt.Errorf("failed to send vk message request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("when sending message, vk responded with non-ok status code: %v", resp.StatusCode)
	}
	return nil
}

func formatPricesForVK(original string, cards []pricedCard) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Оригинальное название: %v\n\n", original))
	for i, priced := range cards {
		builder.WriteString(fmt.Sprintf("%v. %v: $%v\nhttps://starcitygames.com%v\n", i+1, priced.set, priced.price, priced.link))
	}
	return builder.String()
}

func formatPricesForTelegram(original string, cards []pricedCard) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Оригинальное название: %v\n\n", original))
	for i, priced := range cards {
		builder.WriteString(fmt.Sprintf("%v. <a href=\"https://starcitygames.com%v\">%v</a>: $%v\n", i+1, priced.link, priced.set, priced.price))
	}
	return builder.String()
}
