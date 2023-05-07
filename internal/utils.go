package internal

import (
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"strconv"
	"strings"
)

// getUint64ValueFromStringFlag returns the uint64 value of a string flag.
// You can use this function if you want to use uint64 flag but does not want to display the default value.
func getUint64ValueFromStringFlag(c *cli.Context, flagName string) (uint64, error) {
	v := c.String(flagName)
	if v == "" {
		return 0, nil
	}
	uint64v, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64v, nil
}

// isTerminal returns true if the writer is a terminal.
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}

// isPipe returns true if the reader is a pipe.
func isPipe(r io.Reader) bool {
	if f, ok := r.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return stat.Mode()&os.ModeNamedPipe != 0
	}
	return false
}

// trimLeftSpaces trims left spaces and newlines.
func trimLeftSpaces(s string) string {
	return strings.TrimLeft(s, " \n")
}

// truncateChars truncates the string to the specified number of characters.
func truncateChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
