package rnd_test

import (
	"testing"

	"github.com/cdrpl/idlemon/pkg/rnd"
)

// It should return a string of (bytes * 2) len.
func TestGenerateToken(t *testing.T) {
	numBytes := 16

	token, err := rnd.GenerateToken(numBytes)
	if err != nil {
		t.Error(err)
	}

	// token len should be bytes * 2
	if len(token) != numBytes*2 {
		t.Error("Token len should be equal to numBytes * 2")
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = rnd.GenerateToken(16)
	}
}
