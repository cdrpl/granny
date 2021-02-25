package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	// It should set the os.environment.
	reader := strings.NewReader("HOST_TEST=123456\nhost2=num2")
	parseEnvFile(reader)

	val := os.Getenv("HOST_TEST")
	val2 := os.Getenv("host2")
	if val != "123456" || val2 != "num2" {
		t.Error("Env vars were not correctly parsed")
	}

	// It should not overwrite existing values
	os.Setenv("HOST_TEST", "already_set")
	reader = strings.NewReader("HOST_TEST=123456")
	parseEnvFile(reader)

	val = os.Getenv("HOST_TEST")
	if val != "already_set" {
		t.Error("Env vars were overridden")
	}
}

func BenchmarkParseEnvFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader := strings.NewReader("HOST_TEST=123456\nhost2=num2")
		b.StartTimer()
		parseEnvFile(reader)
	}
}
