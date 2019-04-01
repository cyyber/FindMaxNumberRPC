package maxnumber

import (
	"fmt"
	"testing"
)

var testData = []struct {
	value int64
	out1 int64
	out2 bool
} {
	{1, 1, true},
	{5,5, true},
	{3,5, false},
	{6,6, true},
	{2,6, false},
	{20,20, true},
}

func TestNewMaxNumber(t *testing.T) {
	m := NewMaxNumber()
	if m == nil {
		t.Error("Nil was not expected")
	}
}

func TestMaxNumber_FindMaxNumber(t *testing.T) {
	m := NewMaxNumber()
	for _, data := range testData {
		value, isNewMax := m.FindMaxNumber(data.value)
		if value != data.out1 {
			t.Error(fmt.Sprintf("Max value should be %d", data.value))
		}
		if isNewMax != data.out2 {
			t.Error(fmt.Sprintf("isNewMax should have been %t for value %d", data.out2, 10))
		}
	}
}

func TestMaxNumber_GetMaxNumber(t *testing.T) {
	m := NewMaxNumber()
	for _, data := range testData {
		_, _ = m.FindMaxNumber(data.value)
	}
	lastIndex := len(testData) - 1
	value := m.GetMaxNumber()
	if value != testData[lastIndex].out1 {
		t.Error(fmt.Sprintf("Max value should be %d", testData[lastIndex].out1))
	}
}