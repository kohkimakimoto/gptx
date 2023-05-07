package internal

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/sashabaranov/go-openai"
	"os"
	"os/exec"
)

type ChatService struct {
	ClientConfig openai.ClientConfig
	PathResolver *PathResolver
	StoreManager *StoreManager
	CacheManager *CacheManager
	HookFactory  *HookFactory
	Writer       *OutputWriter
	Spinner      *spinner.Spinner
	Conversation *Conversation
	Hooks        []*Hook
	NoLoading    bool
	NoCache      bool
	OnMemory     bool
	Model        string
	Temperature  float32
	TopP         float32
	HooksEnv     []string
}

func (c *ChatService) DisableOutputAnimation() {
	c.Writer.UseAnimation = false
}

// InitConversation resolves the conversation to use.
// If "resume" is specified, it loads the conversation from the store.
func (c *ChatService) InitConversation(resume string, name string, label string) error {
	if resume != "" && name != "" {
		return fmt.Errorf("resume and name are mutually exclusive")
	}

	if resume != "" && label != "" {
		return fmt.Errorf("resume and label are mutually exclusive")
	}

	if resume != "" && c.OnMemory {
		return fmt.Errorf("resume can not be used with on-memory mode")
	}

	var co *Conversation
	if resume != "" {
		// get an existing conversation if "resume" is specified
		_co, err := c.getConversationByKey(NewConversationKey(resume))
		if err != nil {
			return err
		}
		co = _co
	} else {
		// create a new conversation
		co = NewConversation()
		if name != "" {
			// set name if "name" is specified
			ok, err := c.isExistsConversationByName(name)
			if err != nil {
				return err
			}
			if ok {
				return fmt.Errorf("conversation with name %q already exists", name)
			}
			co.Name = name
		}

		if label != "" {
			// set label if "label" is specified
			co.Label = label
		}
	}
	c.Conversation = co
	return nil
}

func (c *ChatService) LoadHooks(hookNames []string) error {
	// before loading hooks, update PATH environment variable
	if err := updatePathEnv(c.PathResolver); err != nil {
		return err
	}

	if !c.Conversation.IsNew() && len(hookNames) > 0 {
		if !equalStringSlice(c.Conversation.Hooks, hookNames) {
			return fmt.Errorf("hooks can be specified for new conversation only")
		}
	}

	if c.Conversation.IsNew() {
		c.Conversation.Hooks = hookNames
	}

	var hooks = make([]*Hook, 0, len(c.Conversation.Hooks))
	for _, name := range c.Conversation.Hooks {
		h, err := c.HookFactory.NewHook(name)
		if err != nil {
			return err
		}
		hooks = append(hooks, h)
	}
	c.Hooks = hooks
	return nil
}

func (c *ChatService) Chat(prompt string) error {
	if c.Conversation.IsNew() {
		c.Conversation.Prompt = prompt
	}

	// run pre-message hooks
	for _, hook := range c.Hooks {
		_prompt, err := c.runPreMessageHook(hook, prompt)
		if err != nil {
			return err
		}
		prompt = _prompt
	}

	// create new message and add it to the conversation
	m := openai.ChatCompletionMessage{}
	m.Role = openai.ChatMessageRoleUser
	m.Content = prompt
	c.Conversation.AddMessage(m)

	if !c.OnMemory {
		// save conversation
		if err := func() error {
			// this anonymous function is used to open and close the store in the specific scope
			store, err := c.StoreManager.Open()
			if err != nil {
				return err
			}
			defer store.Close()

			if c.Conversation.IsNew() {
				if err := store.CreateConversation(c.Conversation); err != nil {
					return err
				}
			} else {
				if err := store.UpdateConversation(c.Conversation); err != nil {
					return err
				}
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	content, err := c.requestChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       c.Model,
		Temperature: c.Temperature,
		TopP:        c.TopP,
		Messages:    c.Conversation.Messages,
	})
	if err != nil {
		return err
	}

	// run post-message hooks
	for _, hook := range c.Hooks {
		_content, err := c.runPostMessageHook(hook, content)
		if err != nil {
			return err
		}
		content = _content
	}

	// save completion as an assistant message
	m = openai.ChatCompletionMessage{}
	m.Role = openai.ChatMessageRoleAssistant
	m.Content = content
	c.Conversation.AddMessage(m)

	if !c.OnMemory {
		if err := func() error {
			// this anonymous function is used to open and close the store in the specific scope
			store, err := c.StoreManager.Open()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.UpdateConversation(c.Conversation); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	c.Writer.Println(content)

	// run finish hooks
	for _, hook := range c.Hooks {
		if err := c.runFinishHook(hook, content); err != nil {
			return err
		}
	}
	return nil
}

func (c *ChatService) getConversationByKey(key *ConversationKey) (*Conversation, error) {
	store, err := c.StoreManager.Open()
	if err != nil {
		return nil, err
	}
	defer store.Close()

	return store.GetConversationByKey(key)
}

func (c *ChatService) isExistsConversationByName(name string) (bool, error) {
	store, err := c.StoreManager.Open()
	if err != nil {
		return false, err
	}
	defer store.Close()

	_, err = store.GetConversationByName(name)
	if err != nil {
		if _, ok := err.(*ConversationNotFoundError); ok {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (c *ChatService) requestChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (string, error) {
	client := openai.NewClientWithConfig(c.ClientConfig)

	if !c.NoCache && !c.OnMemory {
		cache, err := c.CacheManager.Open()
		if err != nil {
			return "", err
		}
		defer cache.Close()

		b, err := json.Marshal(req)
		if err != nil {
			return "", err
		}
		sha := sha256.Sum256(b)
		key := sha[:]
		item, err := cache.Get(key)
		if err != nil {
			if _, ok := err.(*CacheItemNotFoundError); ok {
				// Cache miss. Request to OpenAI API
				c.spinnerStart()
				resp, err := client.CreateChatCompletion(ctx, req)
				c.spinnerStop()
				if err != nil {
					return "", err
				}
				content := resp.Choices[0].Message.Content
				if err := cache.Set(key, []byte(content)); err != nil {
					return "", err
				}
				return content, nil
			} else {
				return "", err
			}
		} else {
			// Cache hit
			return string(item), nil
		}
	} else {
		c.spinnerStart()
		resp, err := client.CreateChatCompletion(ctx, req)
		c.spinnerStop()
		if err != nil {
			return "", err
		}
		content := resp.Choices[0].Message.Content
		return content, nil
	}
}

func (c *ChatService) spinnerStart() {
	if !c.NoLoading {
		c.Spinner.Start()
	}
}

func (c *ChatService) spinnerStop() {
	if !c.NoLoading {
		c.Spinner.Stop()
	}
}

func (c *ChatService) runPreMessageHook(h *Hook, prompt string) (string, error) {
	file, err := os.CreateTemp("", "gptx-prompt-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())
	// Write the prompt to the file.
	if _, err := file.WriteString(prompt); err != nil {
		return "", err
	}

	cmd := h.Command()

	cmdEnv := os.Environ()
	cmdEnv = append(cmdEnv, c.HooksEnv...)
	cmdEnv = append(cmdEnv,
		fmt.Sprintf("GPTX_HOOK_TYPE=%s", HookTypePreMessage),
		fmt.Sprintf("GPTX_MESSAGE_INDEX=%d", len(c.Conversation.Messages)),
		fmt.Sprintf("GPTX_CONVERSATION_ID=%d", c.Conversation.Id),
		fmt.Sprintf("GPTX_PROMPT_FILE=%s", file.Name()),
	)

	cmd.Env = cmdEnv

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// read the modified prompt from the file after the hook command is finished
	b, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *ChatService) runPostMessageHook(h *Hook, completion string) (string, error) {
	file, err := os.CreateTemp("", "gptx-completion-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())
	// Write the completion to the file.
	if _, err := file.WriteString(completion); err != nil {
		return "", err
	}

	cmd := h.Command()

	cmdEnv := os.Environ()
	cmdEnv = append(cmdEnv, c.HooksEnv...)
	cmdEnv = append(cmdEnv,
		fmt.Sprintf("GPTX_HOOK_TYPE=%s", HookTypePostMessage),
		fmt.Sprintf("GPTX_MESSAGE_INDEX=%d", len(c.Conversation.Messages)-1),
		fmt.Sprintf("GPTX_CONVERSATION_ID=%d", c.Conversation.Id),
		fmt.Sprintf("GPTX_COMPLETION_FILE=%s", file.Name()),
	)
	cmd.Env = cmdEnv

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// read the modified completion from the file after the hook command is finished
	b, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *ChatService) runFinishHook(h *Hook, completion string) error {
	file, err := os.CreateTemp("", "gptx-completion-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	// Write the completion to the file.
	if _, err := file.WriteString(completion); err != nil {
		return err
	}

	cmd := h.Command()

	cmdEnv := os.Environ()
	cmdEnv = append(cmdEnv, c.HooksEnv...)
	cmdEnv = append(cmdEnv,
		fmt.Sprintf("GPTX_HOOK_TYPE=%s", HookTypeFinish),
		fmt.Sprintf("GPTX_MESSAGE_INDEX=%d", len(c.Conversation.Messages)-2),
		fmt.Sprintf("GPTX_CONVERSATION_ID=%d", c.Conversation.Id),
		fmt.Sprintf("GPTX_COMPLETION_FILE=%s", file.Name()),
	)
	cmd.Env = cmdEnv

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

const (
	cancelExitCode = 3
)

func isErrCancel(err error) bool {
	if err == nil {
		return false
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	if exitErr.ExitCode() == cancelExitCode {
		return true
	}
	return false
}
