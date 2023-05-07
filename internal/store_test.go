package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversationNotFoundError_Error(t *testing.T) {
	err := &ConversationNotFoundError{Key: 1234}
	assert.Equal(t, "conversation '1234' is not found", err.Error())
}

func TestConversationNameDuplicatedError_Error(t *testing.T) {
	err := &ConversationNameDuplicatedError{Name: "test-conversation", Id: 123}
	assert.Equal(t, "conversation name 'test-conversation' is already used by id 123", err.Error())
}

func TestConversationInvalidNameError_Error(t *testing.T) {
	err := &ConversationInvalidNameError{Name: "1234"}
	assert.Equal(t, "conversation name '1234' is invalid (using only numbers for a name is not allowed)", err.Error())
}

func TestStore_CreateConversation(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// get a stored conversation and compare it with the original one
		co2, err := s.GetConversationById(co.Id)
		assert.NoError(t, err)
		assert.Equal(t, co, co2)
	})

	t.Run("create with the same name", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		co2 := NewConversation()
		co2.Name = "test-conversation1"
		err = s.CreateConversation(co2)
		assert.Error(t, err)
		assert.IsType(t, &ConversationNameDuplicatedError{}, err)
		assert.Equal(t, co.Name, err.(*ConversationNameDuplicatedError).Name)
		assert.Equal(t, co.Id, err.(*ConversationNameDuplicatedError).Id)
	})

	t.Run("create with invalid name", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "1234"
		err = s.CreateConversation(co)
		assert.Error(t, err)
		assert.IsType(t, &ConversationInvalidNameError{}, err)
		assert.Equal(t, co.Name, err.(*ConversationInvalidNameError).Name)
	})
}

func TestStore_UpdateConversation(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		co.Name = "test-conversation1-updated"
		err = s.UpdateConversation(co)
		assert.NoError(t, err)

		// get a stored conversation and compare it with the original one
		co2, err := s.GetConversationById(co.Id)
		assert.NoError(t, err)
		assert.Equal(t, co, co2)

		// get a stored conversation by updated name
		co3, err := s.GetConversationByName(co.Name)
		assert.NoError(t, err)
		assert.Equal(t, co, co3)

		// get a stored conversation by old name
		co4, err := s.GetConversationByName("test-conversation1")
		assert.Error(t, err)
		assert.Nil(t, co4)
	})

	t.Run("update with invalid id", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		co2 := NewConversation()
		co2.Id = 11111 // invalid id that does not exist
		co2.Name = "test-conversation1-updated"
		err = s.UpdateConversation(co2)
		assert.Error(t, err)
	})

	t.Run("update with invalid name", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		co.Name = "1111"
		err = s.UpdateConversation(co)
		assert.Error(t, err)
		assert.IsType(t, &ConversationInvalidNameError{}, err)
		assert.Equal(t, co.Name, err.(*ConversationInvalidNameError).Name)
	})
}

func TestStore_DeleteConversationById(t *testing.T) {
	t.Run("delete a conversation by id", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = s.DeleteConversationById(co.Id)
		assert.NoError(t, err)

		// get a stored conversation and check if it is nil
		co2, err := s.GetConversationById(co.Id)
		assert.Error(t, err)
		assert.Nil(t, co2)
	})

	t.Run("delete a conversation by name", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = s.DeleteConversationById(co.Id)
		assert.NoError(t, err)

		// get a stored conversation and check if it is nil
		co2, err := s.GetConversationByName(co.Name)
		assert.Error(t, err)
		assert.Nil(t, co2)
	})

	t.Run("delete a conversation with invalid id", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		err = s.DeleteConversationById(11111)
		assert.Error(t, err)
	})
}

func TestStore_GetConversationByKey(t *testing.T) {
	t.Run("get a conversation by numeric (id) key", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// get a stored conversation by id
		co2, err := s.GetConversationByKey(NewConversationKey("1"))
		assert.NoError(t, err)
		assert.Equal(t, co, co2)
	})

	t.Run("get a conversation by string (name) key", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// get a stored conversation by name
		co2, err := s.GetConversationByKey(NewConversationKey("test-conversation1"))
		assert.NoError(t, err)
		assert.Equal(t, co, co2)
	})
}

func TestStore_RenameConversationByKey(t *testing.T) {
	t.Run("rename a conversation", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// rename a conversation
		err = s.RenameConversationByKey(NewConversationKey("1"), "test-conversation1-renamed")
		assert.NoError(t, err)

		// get a stored conversation and compare it with the original one
		co2, err := s.GetConversationById(co.Id)
		assert.NoError(t, err)
		assert.Equal(t, "test-conversation1-renamed", co2.Name)
	})

	t.Run("rename a conversation with invalid key", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// rename a conversation
		err = s.RenameConversationByKey(NewConversationKey("11111"), "test-conversation1-renamed")
		assert.Error(t, err)
	})

	t.Run("rename a conversation with invalid name", func(t *testing.T) {
		sm := testStoreManager(t)
		s, err := sm.Open()
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "test-conversation1"
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		// rename a conversation but numbers only name is invalid
		err = s.RenameConversationByKey(NewConversationKey("1"), "123")
		assert.Error(t, err)
		assert.IsType(t, &ConversationInvalidNameError{}, err)
		assert.Equal(t, "123", err.(*ConversationInvalidNameError).Name)
	})
}

func TestStore_ListConversions(t *testing.T) {
	sm := testStoreManager(t)
	s, err := sm.Open()
	assert.NoError(t, err)

	// test data
	conversations := make([]*Conversation, 0)
	for i := 1; i <= 30; i++ {
		co := NewConversation()
		co.Name = fmt.Sprintf("test-conversation-%d", i)
		if i <= 10 {
			co.Label = "test-label"
		}

		err := s.CreateConversation(co)
		assert.NoError(t, err)

		conversations = append(conversations, co)
	}

	rConversations := reverseConversationsSlice(conversations)

	tests := []struct {
		testName string
		query    *ListConversationsQuery
		want     *ConversationList
	}{
		// default query
		{
			query: &ListConversationsQuery{},
			want: &ConversationList{
				Conversations: conversations,
				HasNext:       false,
				Next:          nil,
				Count:         30,
			},
		},
		// default query with label
		{
			query: &ListConversationsQuery{
				Label: "test-label",
			},
			want: &ConversationList{
				Conversations: conversations[:10],
				HasNext:       false,
				Next:          nil,
				Count:         10,
			},
		},
		// limit 10
		{
			query: &ListConversationsQuery{
				Limit: 10,
			},
			want: &ConversationList{
				Conversations: conversations[:10],
				HasNext:       true,
				Next:          func() *uint64 { n := uint64(11); return &n }(),
				Count:         10,
			},
		},
		// limit 10, offset 11
		{
			query: &ListConversationsQuery{
				Limit: 10,
				Begin: func() *uint64 { n := uint64(11); return &n }(),
			},
			want: &ConversationList{
				Conversations: conversations[10:20],
				HasNext:       true,
				Next:          func() *uint64 { n := uint64(21); return &n }(),
				Count:         10,
			},
		},
		// limit 10, offset 21
		{
			query: &ListConversationsQuery{
				Limit: 10,
				Begin: func() *uint64 { n := uint64(21); return &n }(),
			},
			want: &ConversationList{
				Conversations: conversations[20:],
				HasNext:       false,
				Next:          nil,
				Count:         10,
			},
		},
		// limit 10 reverse
		{
			query: &ListConversationsQuery{
				Limit:   10,
				Reverse: true,
			},
			want: &ConversationList{
				Conversations: rConversations[:10],
				HasNext:       true,
				Next:          func() *uint64 { n := uint64(20); return &n }(),
				Count:         10,
			},
		},
		// limit 10, offset 20 reverse
		{
			query: &ListConversationsQuery{
				Limit:   10,
				Begin:   func() *uint64 { n := uint64(20); return &n }(),
				Reverse: true,
			},
			want: &ConversationList{
				Conversations: rConversations[10:20],
				HasNext:       true,
				Next:          func() *uint64 { n := uint64(10); return &n }(),
				Count:         10,
			},
		},
		// limit 10, offset 10 reverse
		{
			query: &ListConversationsQuery{
				Limit:   10,
				Begin:   func() *uint64 { n := uint64(10); return &n }(),
				Reverse: true,
			},
			want: &ConversationList{
				Conversations: rConversations[20:],
				HasNext:       false,
				Next:          nil,
				Count:         10,
			},
		},
	}

	for _, tt := range tests {
		tt.want.query = tt.query
		got, err := s.ListConversations(tt.query)
		assert.NoError(t, err)
		assert.Equal(t, tt.want, got)
	}
}

func reverseConversationsSlice(input []*Conversation) []*Conversation {
	result := make([]*Conversation, len(input))

	copy(result, input)

	// Reverse the result slice in place
	for i := 0; i < len(result)/2; i++ {
		j := len(result) - i - 1
		result[i], result[j] = result[j], result[i]
	}

	return result
}
