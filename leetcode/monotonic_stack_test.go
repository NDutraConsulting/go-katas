package leetcode

import (
	"testing"
)

func TestDailyTemp(t *testing.T) {
	tests := []struct {
		temperatures []int
		expect       []int
	}{
		{temperatures: []int{73, 74, 75, 71, 69, 72, 76, 73}, expect: []int{1, 1, 4, 2, 1, 1, 0, 0}},
	}
	for _, tc := range tests {
		if got := dailyTemperatures(tc.temperatures); !equal(got, tc.expect) {
			t.Errorf("dailyTemperatures(%v) = %v; want %v", tc.temperatures, got, tc.expect)
		}
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
