package handlers

import (
	"testing"
)

func TestFindingKeyInURLPath(t *testing.T) {
	testCases := []struct {
		test string
		want string
	}{
		{"/thisorthat", "thisorthat"},
		{"/this-that", "this-that"},
		{"/lol/lol/lol/", ""},
	}

	for _, tc := range testCases {
		got := findKey(tc.test)
		if got != tc.want {
			t.Errorf("want %s, got %s", tc.want, got)
		}
	}
}
