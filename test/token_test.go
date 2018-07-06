package test

import (
	"github.com/byxor/qbd/token"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExtractingTokens(t *testing.T) {
	type Entry struct {
		input  []byte
		output token.Token
	}

	entries := []Entry{
		{[]byte{}, token.None},

		{[]byte{0x00}, token.EndOfFile},
		{[]byte{0x01}, token.EndOfLine},
		{[]byte{0x03}, token.StartOfStruct},
		{[]byte{0x04}, token.EndOfStruct},
		{[]byte{0x05}, token.StartOfArray},
		{[]byte{0x06}, token.EndOfArray},

		{[]byte{0x16, 0x00, 0x00, 0x00, 0x00}, token.Name},
		{[]byte{0x16, 0xBB, 0xEE, 0xEE, 0xFF}, token.Name},

		// Invalid names, not enough bytes
		{[]byte{0x16, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x16, 0x11, 0x22}, token.Invalid},
		{[]byte{0x16, 0x33}, token.Invalid},
		{[]byte{0x16}, token.Invalid},

		{[]byte{0x17, 0x00, 0x00, 0x00, 0x00}, token.Integer},
		{[]byte{0x17, 0xBA, 0x5E, 0xBA, 0x11}, token.Integer},

		// Invalid integers, not enough bytes
		{[]byte{0x17, 0x00, 0x00, 0x00}, token.Invalid},
		{[]byte{0x17, 0x11, 0x22}, token.Invalid},
		{[]byte{0x17, 0x33}, token.Invalid},
		{[]byte{0x17}, token.Invalid},

		{[]byte{0x1A, 0x00, 0x00, 0x00, 0x00}, token.Float},
		{[]byte{0x1A, 0x12, 0x34, 0x56, 0x78}, token.Float},

		{[]byte{0x23}, token.StartOfFunction},
		{[]byte{0x24}, token.EndOfFunction},
		{[]byte{0x29}, token.Return},

		{[]byte{0x25}, token.StartOfIf},
		{[]byte{0x26}, token.Else},
		{[]byte{0x27}, token.ElseIf},
		{[]byte{0x28}, token.EndOfIf},
	}

	for _, entry := range entries {
		tokens := make(chan token.Token)
		go token.GetTokens(tokens, entry.input)

		select {
		case token := <-tokens:
			assert.Equal(t, entry.output, token)
		case <-time.After(1 * time.Second):
			assert.Equal(t, "timeout", "!!!")
		}
	}
}
