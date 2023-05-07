package internal

import (
	"bytes"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSimpleTableWriter_Render(t *testing.T) {
	out := &bytes.Buffer{}
	tw := NewSimpleTableWriter(out)
	tw.AppendHeader(table.Row{
		"AAA",
		"BBB",
		"CCC",
	})
	tw.AppendRow([]interface{}{"a", "b", "c"})
	tw.AppendRow([]interface{}{"aaa", "bbb", "ccc"})
	tw.AppendRow([]interface{}{"aaaaaa", "bbbbbb", "cccccc"})
	tw.Render()

	ret := out.String()
	assert.Equal(t, strings.TrimPrefix(`
AAA      BBB      CCC
a        b        c
aaa      bbb      ccc
aaaaaa   bbbbbb   cccccc
`, "\n"), ret)

}
