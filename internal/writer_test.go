package internal

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOutputWriter_Println(t *testing.T) {
	t.Run("without animation", func(t *testing.T) {
		var buf bytes.Buffer
		w := &OutputWriter{
			Writer:         &buf,
			UseAnimation:   false,
			AnimationSpeed: 10 * time.Millisecond,
			Color:          color.New(color.FgMagenta, color.Bold),
		}
		w.UseAnimation = false
		w.Println("Hello, world!")
		assert.Equal(t, "Hello, world!\n", buf.String())
	})

	t.Run("with animation", func(t *testing.T) {
		// does not check animation. just check if it outputs text.
		var buf bytes.Buffer
		w := &OutputWriter{
			Writer:         &buf,
			UseAnimation:   true,
			AnimationSpeed: 10 * time.Millisecond,
			Color:          color.New(color.FgMagenta, color.Bold),
		}
		w.UseAnimation = true
		w.Println("Hello, world!")
		assert.Equal(t, "Hello, world!\n", buf.String())
	})
}
