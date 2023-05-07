package internal

import (
	"bytes"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"testing"
)

func testNewApp(t *testing.T) *cli.App {
	t.Helper()
	app := NewApp(testNewRepository(t))
	app.Writer = &bytes.Buffer{}
	app.ErrWriter = &bytes.Buffer{}
	return app
}

func testNewRepository(t *testing.T) *Repository {
	t.Helper()
	return NewRepository(testNewPathResolver(t))
}

func testNewPathResolver(t *testing.T) *PathResolver {
	t.Helper()
	dir, err := os.MkdirTemp("", "gptx_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return NewPathResolver(dir)
}

func testStoreManager(t *testing.T) *StoreManager {
	t.Helper()

	tempFile, err := os.CreateTemp("", "gptx_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})

	sm := &StoreManager{
		DBPath: tempFile.Name(),
	}
	t.Cleanup(func() {
		_ = sm.Close()
	})

	store, err := sm.Open()
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Init(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	return sm
}

func testCacheManager(t *testing.T) *CacheManager {
	t.Helper()

	tempFile, err := os.CreateTemp("", "gptx_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})

	//c := (tempFile.Name())

	cm := &CacheManager{
		DBPath:    tempFile.Name(),
		MaxLength: 0,
	}
	t.Cleanup(func() {
		_ = cm.Close()
	})

	c, err := cm.Open()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	return cm
}

func testTempFile(t *testing.T, content []byte) *os.File {
	t.Helper()
	tempFile, err := os.CreateTemp("", "gptx_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempFile.Name())
	})

	_, err = tempFile.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	return tempFile
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testHttpClient(t *testing.T, fn RoundTripFunc) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: fn,
	}
}
