package internal

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"time"
)

// OutputWriter is an io.Writer wrapper that enables printing text with a typewriter-style animation.
type OutputWriter struct {
	Writer         io.Writer
	UseAnimation   bool
	AnimationSpeed time.Duration
	Color          *color.Color
}

func (w *OutputWriter) Println(text string) {
	if w.UseAnimation {
		for _, c := range text {
			_, _ = w.Color.Fprintf(w.Writer, "%c", c)
			time.Sleep(w.AnimationSpeed)
		}
		_, _ = fmt.Fprintln(w.Writer)
	} else {
		_, _ = w.Color.Fprintln(w.Writer, text)
	}
}
