package protocol

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())

	parsedCmd, _ := ParseCommand(r)

	assert.Equal(t, cmd, parsedCmd)
}

func TestParseSGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	r := bytes.NewReader(cmd.Bytes())

	parsedCmd, _ := ParseCommand(r)

	assert.Equal(t, cmd, parsedCmd)
}

func BenchmarkParseCommand(b *testing.B) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	r := bytes.NewReader(cmd.Bytes())

	for i := 0; i < b.N; i++ {
		parseSetCommand(r)
	}
}
