package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dentych/taskeroo/internal/database"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const telegramUrl = "https://api.telegram.org"

type Telegram struct {
	repo    *database.TelegramRepo
	client  *http.Client
	token   string
	baseUrl string

	lastID  int
	context context.Context
}

func NewTelegram(repo *database.TelegramRepo, token string) *Telegram {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &Telegram{repo: repo, client: client, token: token, baseUrl: fmt.Sprintf("%s/bot%s", telegramUrl, token)}
}

func (t *Telegram) Start() error {
	t.context = context.Background()
	go t.loop()
	return nil
}

func (t *Telegram) loop() {
	for {
		resp, err := t.client.Get(fmt.Sprintf("%s/getUpdates?offset=%d&timeout=25", t.baseUrl, t.lastID+1))
		if err != nil {
			log.Printf("Failed to get updates from Telegram API: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode >= 300 {
			log.Printf("Telegram API answered with non-successful HTTP status code: %d\n", resp.StatusCode)
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read Telegram API response body: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var updateResponse UpdateResponse
		err = json.Unmarshal(body, &updateResponse)
		if err != nil {
			log.Printf("Failed to unmarshal Telegram API response body: %s\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for _, update := range updateResponse.Result {
			err = t.handleUpdate(update)
			if err != nil {
				log.Printf("There was an error handling Telegram update: %s\n", err)
			}
		}
	}
}

func (t *Telegram) handleCommand(msg Message) error {
	switch msg.Text {
	case "/start":
		return t.sendMessage(msg.From.ID, "Hej! Skriv /connect for at forbinde din Taskeroo konto med denne bot.")
	case "/connect":
		return t.handleConnect(msg)
	default:
		return t.sendMessage(msg.From.ID, "Kommando ikke forstået. Prøv en anden!")
	}
}

func (t *Telegram) sendMessage(telegramUserID int, text string) error {
	body := SendMessageParameters{
		ChatID: strconv.Itoa(telegramUserID),
		Text:   text,
	}
	formatted, err := json.Marshal(&body)
	if err != nil {
		return err
	}
	resp, err := t.client.Post(fmt.Sprintf("%s/sendMessage", t.baseUrl), "application/json", bytes.NewReader(formatted))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		log.Printf("Telegram send message returned non OK status code: %d\n", resp.StatusCode)
		return nil
	}
	return nil
}

func (t *Telegram) handleUpdate(update Update) error {
	if update.UpdateID > t.lastID {
		t.lastID = update.UpdateID
	}

	if update.Message == nil {
		return nil
	}

	if strings.HasPrefix(update.Message.Text, "/") {
		return t.handleCommand(*update.Message)
	}

	log.Printf("Failed to handle update, as it wasn't a known type. Message was: %s\n", update.Message.Text)
	return nil
}

func (t *Telegram) handleConnect(msg Message) error {
	connectID := uuid.NewString()
	err := t.repo.DeleteAllByTelegramUserID(t.context, msg.From.ID)
	if err != nil {
		return err
	}
	err = t.repo.Create(t.context, database.NewTelegram{
		ID:             connectID,
		TelegramUserID: msg.From.ID,
	})
	if err != nil {
		return err
	}

	err = t.sendMessage(msg.From.ID, "For at forbinde din Telegram konto med din Taskeroo konto, "+
		"skal du gøre følgende:\n\n"+
		"1. Log ind på https://taskeroo.tychsen.me, hvis du ikke er i forvejen\n"+
		"2. Følg dette link: "+fmt.Sprintf("https://taskeroo.tychsen.me/telegram-connect/%s", connectID))
	if err != nil {
		return err
	}
	return nil
}

func (t *Telegram) SendMessage(ctx context.Context, telegramUserID int, msg string) error {
	return t.sendMessage(telegramUserID, msg)
}
