package current_test

import (
	"savemanager/current"
	"testing"
)

var existResult = []bool{}

func BenchmarkProcess(b *testing.B) {
	str := "eldenring.exe"
	for i := 0; i < b.N; i++ {

		current.IsProcessRunning(str)
		// existResult = append(existResult, exists)
	}
}
