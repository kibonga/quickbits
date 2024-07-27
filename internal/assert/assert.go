package assert

import "testing"

func Equal[T comparable](t *testing.T, value, expected T) {
	t.Helper()

	if value != expected {
		t.Errorf("got %v; want %v", value, expected)
	}
}
