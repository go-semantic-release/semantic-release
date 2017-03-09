package condition

import (
	"testing"
)

func TestTravisValid(t *testing.T) {
	err := Travis("", "", false)
	if err == nil {
		t.Fail()
	}
}
