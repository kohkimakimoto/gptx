package internal

import (
	"os"
	"path/filepath"
)

type PathResolver struct {
	Dir            string
	pathEnvUpdated bool
}

func NewPathResolver(dir string) *PathResolver {
	return &PathResolver{
		Dir: dir,
	}
}

const AppHomeEnvKey = "GPTX_HOME"

func getAppHomeDir() string {
	appHomeDir := os.Getenv(AppHomeEnvKey)
	if appHomeDir == "" {
		appHomeDir = filepath.Join(os.Getenv("HOME"), ".gptx")
	}
	return appHomeDir
}

func (r *PathResolver) ConfigFilePath() string {
	return filepath.Join(r.Dir, "config.toml")
}

func (r *PathResolver) DBFilePath() string {
	return filepath.Join(r.Dir, "gptx.db")
}

func (r *PathResolver) CacheDBFilePath() string {
	return filepath.Join(r.Dir, "cache.db")
}

func (r *PathResolver) HistoryFilePath() string {
	return filepath.Join(r.Dir, "history.txt")
}

func (r *PathResolver) LibExecDir() string {
	return filepath.Join(r.Dir, "libexec")
}

func (r *PathResolver) LibExecFilePath(name string) string {
	return filepath.Join(r.LibExecDir(), name)
}

func updatePathEnv(pr *PathResolver) error {
	if pr.pathEnvUpdated {
		return nil
	}

	path := os.Getenv("PATH")
	if path == "" {
		path = pr.LibExecDir()
	} else {
		path = pr.LibExecDir() + ":" + path
	}
	if err := os.Setenv("PATH", path); err != nil {
		return err
	}
	pr.pathEnvUpdated = true
	return nil
}
