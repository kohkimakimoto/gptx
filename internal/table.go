package internal

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"io"
	"regexp"
)

// SimpleTableWriter is a simple table Writer that does not draw borders.
type SimpleTableWriter struct {
	table.Writer
	Out io.Writer
}

func NewSimpleTableWriter(out io.Writer) table.Writer {
	// custom borderless table style
	style := table.StyleDefault
	style.Box.PaddingLeft = ""
	style.Box.PaddingRight = ""
	style.Box.MiddleVertical = "   "
	style.Options = table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
	style.Format.Header = text.FormatDefault

	w := table.NewWriter()
	w.SetStyle(style)
	return &SimpleTableWriter{
		Writer: w,
		Out:    out,
	}
}

var reRemoveTrailingSpace = regexp.MustCompile(`\s+\n`)

func (t *SimpleTableWriter) Render() string {
	// Wrap the table.Writer's Render() method to remove trailing spaces.
	outStr := t.Writer.Render()
	outStr = reRemoveTrailingSpace.ReplaceAllString(outStr, "\n")
	_, _ = fmt.Fprintln(t.Out, outStr)
	return outStr
}
