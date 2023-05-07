package internal

import (
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConversationKey(t *testing.T) {
	tests := []struct {
		value    string
		expected *ConversationKey
	}{
		{
			value: "test",
			expected: &ConversationKey{
				Value:   "test",
				IsId:    false,
				IsEmpty: false,
				IdValue: 0,
			},
		},
		{
			value: "",
			expected: &ConversationKey{
				Value:   "",
				IsId:    false,
				IsEmpty: true,
				IdValue: 0,
			},
		},
		{
			value: "1",
			expected: &ConversationKey{
				Value:   "1",
				IsId:    true,
				IsEmpty: false,
				IdValue: 1,
			},
		},
		{
			value: "1a1",
			expected: &ConversationKey{
				Value:   "1a1",
				IsId:    false,
				IsEmpty: false,
				IdValue: 0,
			},
		},
	}

	for _, tt := range tests {
		ret := NewConversationKey(tt.value)
		assert.Equal(t, tt.expected, ret)
	}
}

func TestConversation_IsNew(t *testing.T) {
	tests := []struct {
		co       *Conversation
		expected bool
	}{
		{
			co: &Conversation{
				Id: 0,
			},
			expected: true,
		},
		{
			co: &Conversation{
				Id: 1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.co.IsNew())
	}
}

func TestConversation_AddMessage(t *testing.T) {
	co := NewConversation()
	co.AddMessage(openai.ChatCompletionMessage{
		Role:    "user",
		Content: "test1",
	})
	assert.Equal(t, 1, len(co.Messages))
}

func TestCheckValidConversationName(t *testing.T) {
	tests := []struct {
		value   string
		isValid bool
	}{
		{
			value:   "test",
			isValid: true,
		},
		{
			value:   "",
			isValid: false,
		},
		{
			value:   "1",
			isValid: false,
		},
		{
			value:   "1a1",
			isValid: true,
		},
	}

	for _, tt := range tests {
		err := checkValidConversationName(tt.value)
		if tt.isValid {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestConversationList_TryAppendConversation(t *testing.T) {
	t.Run("append", func(t *testing.T) {
		l := &ConversationList{
			Conversations: []*Conversation{},
		}
		l.TryAppendConversation(&Conversation{})
		l.TryAppendConversation(&Conversation{})
		l.TryAppendConversation(&Conversation{})

		assert.Equal(t, 3, len(l.Conversations))
	})

	t.Run("limit", func(t *testing.T) {
		l := &ConversationList{
			query: &ListConversationsQuery{
				Limit: 2,
			},
			Conversations: []*Conversation{},
		}
		l.TryAppendConversation(&Conversation{})
		l.TryAppendConversation(&Conversation{})
		l.TryAppendConversation(&Conversation{})

		assert.Equal(t, 2, len(l.Conversations))
	})

	t.Run("label", func(t *testing.T) {
		l := &ConversationList{
			query: &ListConversationsQuery{
				Label: "test",
			},
			Conversations: []*Conversation{},
		}
		l.TryAppendConversation(&Conversation{
			Label: "test",
		})
		l.TryAppendConversation(&Conversation{
			Label: "test",
		})
		l.TryAppendConversation(&Conversation{
			Label: "not-test",
		})

		assert.Equal(t, 2, len(l.Conversations))
	})
}
