package story

import (
	"testing"
)

func TestNewHistory(t *testing.T) {
	NewHistory()
}

func BenchmarkNewHistory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewHistory()
	}
}
