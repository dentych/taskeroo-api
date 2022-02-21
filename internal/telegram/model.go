package telegram

type UpdateResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	ID   int    `json:"message_id"`
	From From   `json:"from"`
	Text string `json:"text"`
}

type From struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type SendMessageParameters struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}
