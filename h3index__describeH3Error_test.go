// Tests ported from testDescribeH3Error.c
package h3

import (
	"testing"
)

func TestDescribeH3Error_NoError(t *testing.T) {
	t.Parallel()
	err := eSuccess
	result := describeH3Error(err)
	expected := "Success"
	if result != expected {
		t.Errorf("describeH3Error(%v) = %q, want %q", err, result, expected)
	}
}

func TestDescribeH3Error_InvalidCell(t *testing.T) {
	t.Parallel()
	err := eCellInvalid
	result := describeH3Error(err)
	expected := "Cell argument was not valid"
	if result != expected {
		t.Errorf("describeH3Error(%v) = %q, want %q", err, result, expected)
	}
}

func TestDescribeH3Error_InvalidH3Error(t *testing.T) {
	t.Parallel()
	err := h3Error(9001) // Will probably never hit this
	result := describeH3Error(err)
	expected := "Invalid error code"
	if result != expected {
		t.Errorf("describeH3Error(%v) = %q, want %q", err, result, expected)
	}
}
