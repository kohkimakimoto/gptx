package internal

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kohkimakimoto/gptx/internal/builtin"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"sync"
	"time"
)

// Repository holds the application context data.
type Repository struct {
	PathResolver *PathResolver
	Config       *Config
	ClientConfig openai.ClientConfig
	CacheManager *CacheManager
	StoreManager *StoreManager
	// The following parameters are used internally of this object.
	inited bool
	lock   sync.RWMutex
}

// repositoryAwareAction wraps the action function to inject the Repository
// and ensure that the repository is closed after the action is done.
func repositoryAwareAction(action func(c *cli.Context, r *Repository) error) func(c *cli.Context) error {
	return func(c *cli.Context) (rErr error) {
		r, err := getRepository(c.App)
		if err != nil {
			return err
		}

		// ensure that the repository is closed when the action is done
		defer func() {
			if err := r.Close(); err != nil {
				if rErr != nil {
					rErr = fmt.Errorf("%s: %s", rErr.Error(), err.Error())
				} else {
					rErr = err
				}
			}
		}()

		return action(c, r)
	}
}

type HomeDirNotFoundError struct {
	Dir string
}

func (e *HomeDirNotFoundError) Error() string {
	return fmt.Sprintf("The gptx configuration has not been initialized yet in the directory '%s'. Please run 'gptx init'.", e.Dir)
}

func NewRepository(pr *PathResolver) *Repository {
	return &Repository{
		PathResolver: pr,
		Config:       NewConfig(),
	}
}

func (r *Repository) Init() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.inited {
		return nil
	}

	if _, err := os.Stat(r.PathResolver.Dir); os.IsNotExist(err) {
		return &HomeDirNotFoundError{Dir: r.PathResolver.Dir}
	}

	if _, err := os.Stat(r.PathResolver.LibExecDir()); os.IsNotExist(err) {
		err = os.MkdirAll(r.PathResolver.LibExecDir(), os.FileMode(0700))
		if err != nil {
			return err
		}
	}

	if err := builtin.InitLibexecFiles(r.PathResolver.LibExecDir()); err != nil {
		return err
	}

	// Resolve config file path (default: ~/.gptx/gptx.toml) and load the config.
	cFilePath := r.PathResolver.ConfigFilePath()
	if _, err := os.Stat(cFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(cFilePath, []byte(initialConfig), os.FileMode(0600)); err != nil {
			return err
		}
	}
	if err := r.Config.LoadFromFile(r.PathResolver.ConfigFilePath()); err != nil {
		return err
	}

	// Load config values from environment variables
	if v := os.Getenv("OPENAI_API_KEY"); v != "" && r.Config.OpenAIAPIKey == "" {
		r.Config.OpenAIAPIKey = v
	}

	r.ClientConfig = openai.DefaultConfig(r.Config.OpenAIAPIKey)

	// init store
	r.StoreManager = &StoreManager{
		DBPath: r.PathResolver.DBFilePath(),
	}
	store, err := r.StoreManager.Open()
	if err != nil {
		return err
	}
	defer store.Close()
	if err := store.Init(); err != nil {
		return err
	}

	// init Cache
	r.CacheManager = &CacheManager{
		DBPath:    r.PathResolver.CacheDBFilePath(),
		MaxLength: r.Config.MaxCacheLength,
	}
	cache, err := r.CacheManager.Open()
	if err != nil {
		return err
	}
	defer cache.Close()
	if err := cache.Init(); err != nil {
		return err
	}

	r.inited = true

	return nil
}

func (r *Repository) NewChatService(w io.Writer) (*ChatService, error) {
	c := &ChatService{}
	c.PathResolver = r.PathResolver
	c.ClientConfig = r.ClientConfig
	c.StoreManager = r.StoreManager
	c.CacheManager = r.CacheManager
	c.HookFactory = &HookFactory{}
	c.Writer = &OutputWriter{
		Writer:         w,
		UseAnimation:   isTerminal(w),
		AnimationSpeed: 10 * time.Millisecond,
		Color:          color.New(color.FgMagenta, color.Bold),
	}

	sp := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(w))
	if err := sp.Color("green"); err != nil {
		return nil, err
	}
	sp.Suffix = color.GreenString(" Waiting for ChatGPT to respond...")
	c.Spinner = sp

	return c, nil
}

func (r *Repository) Close() error {
	if r.StoreManager != nil {
		if err := r.StoreManager.Close(); err != nil {
			return err
		}
	}

	if r.CacheManager != nil {
		if err := r.CacheManager.Close(); err != nil {
			return err
		}
	}

	return nil
}
