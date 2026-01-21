package concurrentsensors

import (
	"testing"
)

func TestOilFeild(t *testing.T) {

	result := sim1()
	expect := true
	if result != expect {
		t.Errorf("sim1 result: pass = %t; want pass = %t", result, expect)
	}
}
