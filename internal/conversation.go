package internal

import (
	"github.com/sashabaranov/go-openai"
	"regexp"
	"strconv"
	"time"
)

type ConversationKey struct {
	Value   string
	IsId    bool
	IsEmpty bool
	IdValue uint64
}

var idRegex = regexp.MustCompile(`^\d+$`)

func NewConversationKey(value string) *ConversationKey {
	k := &ConversationKey{
		Value:   value,
		IsEmpty: false,
	}
	if value == "" {
		k.IsEmpty = true
		return k
	}

	if idRegex.MatchString(value) {
		if uint64v, err := strconv.ParseUint(value, 10, 64); err == nil {
			k.IsId = true
			k.IdValue = uint64v
		}
	}

	return k
}

type Conversation struct {
	Id        uint64                         `json:"id"`              // Conversation ID
	Prompt    string                         `json:"prompt"`          // The initial input text that starts the conversation
	Name      string                         `json:"name,omitempty"`  // The unique name of the conversation
	Label     string                         `json:"label,omitempty"` // A label for categorizing the conversation
	CreatedAt time.Time                      `json:"created_at"`      // When the conversation was created
	Messages  []openai.ChatCompletionMessage `json:"messages"`        // Messages in the conversation
	Hooks     []string                       `json:"hooks,omitempty"` // Registered Hooks for the conversation
}

func NewConversation() *Conversation {
	return &Conversation{
		Id:        0,
		Prompt:    "",
		Name:      "",
		Label:     "",
		CreatedAt: time.Now().UTC(),
		Messages:  []openai.ChatCompletionMessage{},
		Hooks:     []string{},
	}
}

func (c *Conversation) IsNew() bool {
	return c.Id == 0
}

func (c *Conversation) AddMessage(msg openai.ChatCompletionMessage) {
	c.Messages = append(c.Messages, msg)
}

func checkValidConversationName(name string) error {
	k := NewConversationKey(name)
	if !k.IsEmpty && !k.IsId {
		return nil
	} else {
		return &ConversationInvalidNameError{
			Name: name,
		}
	}
}

type ListConversationsQuery struct {
	Begin   *uint64
	Reverse bool
	Limit   int
	Label   string
}

type ConversationList struct {
	query         *ListConversationsQuery
	Conversations []*Conversation `json:"conversations"`
	HasNext       bool            `json:"has_next"`
	Next          *uint64         `json:"next,string,omitempty"`
	Count         int             `json:"count"`
}

func (l *ConversationList) TryAppendConversation(c *Conversation) {
	if l.query != nil && l.query.Label != "" && c.Label != l.query.Label {
		// skip if label is specified and it doesn't match
		return
	}

	if l.IsLimitReached() {
		// skip if limit is specified and it's reached
		return
	}

	l.Conversations = append(l.Conversations, c)
}

func (l *ConversationList) IsLimitReached() bool {
	return l.query != nil && l.query.Limit > 0 && len(l.Conversations) >= l.query.Limit
}
