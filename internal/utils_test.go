package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestGetUint64ValueFromStringFlag(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
		hasError bool
	}{
		{
			input:    "",
			expected: 0,
		},
		{
			input:    "0",
			expected: 0,
		},
		{
			input:    "1",
			expected: 1,
		},
		{
			input:    "a",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		a := cli.NewApp()
		a.Flags = []cli.Flag{
			&cli.StringFlag{
				Name: "testFlag",
			},
		}
		a.Action = func(ctx *cli.Context) error {
			ret, err := getUint64ValueFromStringFlag(ctx, "testFlag")
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, ret)

			return nil
		}

		err := a.Run([]string{"gptx", "--testFlag", tt.input})
		assert.NoError(t, err)
	}
}

func TestTrimLeftSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    " aaaa",
			expected: "aaaa",
		},
		{
			input:    " \naaaa",
			expected: "aaaa",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, trimLeftSpaces(tt.input))
	}
}

func TestTruncateChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		length   int
	}{
		{
			input:    "abcd",
			expected: "abcd",
			length:   5,
		},
		{
			input:    "abcde",
			expected: "abcde",
			length:   5,
		},
		{
			input:    "abcdef",
			expected: "ab...",
			length:   5,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, truncateChars(tt.input, tt.length))
	}
}

func TestEqualStringSlice(t *testing.T) {
	tests := []struct {
		input1   []string
		input2   []string
		expected bool
	}{
		{
			input1:   []string{"a", "b", "c"},
			input2:   []string{"a", "b", "c"},
			expected: true,
		},
		{
			input1:   []string{"a", "b", "c"},
			input2:   []string{"a", "b", "d"},
			expected: false,
		},
		{
			input1:   []string{"a", "b", "c"},
			input2:   []string{"a", "b"},
			expected: false,
		},
		{
			input1:   []string{"a", "b", "c"},
			input2:   []string{"a", "b", "c", "d"},
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, equalStringSlice(tt.input1, tt.input2))
	}
}
