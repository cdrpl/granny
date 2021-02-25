package main

import (
	"testing"
)

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = generateToken(16)
	}
}
