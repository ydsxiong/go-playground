package utils

import (
	"testing"
)

func TestFib(t *testing.T) {

	tests := []uint64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	funcToCheck := fib
	for index, expected := range tests {
		input := index + 1
		if output := funcToCheck(input); output != expected {
			t.Fatalf("at index %d, expected %d, but got %d.", index, expected, output)
		}
	}
}

func TestFibV2(t *testing.T) {

	tests := []uint64{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	funcToCheck := fib_v2()
	for index, expected := range tests {
		if output := funcToCheck(); output != expected {
			t.Fatalf("at index %d, expected %d, but got %d.", index, expected, output)
		}
	}
}

func TestSumOf10naturalnumbersimperative(t *testing.T) {
	expected := 275
	if output := sumOf10naturalnumbersbyimperativeapproach(); output != expected {
		t.Fatalf("expected %d, but got %d.", expected, output)
	}
}

func TestSumOf10naturalnumbersresursion(t *testing.T) {
	expected := 275
	if output := sumOf10naturalnumbersbyRecursionapproach(); output != expected {
		t.Fatalf("expected %d, but got %d.", expected, output)
	}
}

func TestSumOf10naturalnumbersfunctional(t *testing.T) {
	expected := 275
	if output := sumOf10naturalnumbersbyfunctionalapproach(); output != expected {
		t.Fatalf("expected %d, but got %d.", expected, output)
	}
}

func BenchmarkFibV2(b *testing.B) {
	funcToBenchmark := fib_v2()
	for i := 0; i < b.N; i++ {
		_ = funcToBenchmark()
	}
}
