// Tests ported from H3 v4.4.0: src/apps/testapps/testDescribeH3Error.c.
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

// The following tests are ported from H3 C 4.4.0 testDescribeH3Error.c.

func TestDescribeH3Error_InvalidH3ErrorEnd(t *testing.T) {
	t.Parallel()
	if got := describeH3Error(h3ErrorEnd); got != "Invalid error code" {
		t.Errorf("describeH3Error(h3ErrorEnd) = %q, want invalid-code message", got)
	}
}

func TestDescribeH3Error_InvalidH3ErrorEndPlus(t *testing.T) {
	t.Parallel()
	// Try to catch if someone adds an error code after H3_ERROR_END.
	if got := describeH3Error(h3ErrorEnd + 1); got != "Invalid error code" {
		t.Errorf("describeH3Error(h3ErrorEnd+1) = %q, want invalid-code message", got)
	}
}

func TestDescribeH3Error_ErrorCodesNotValidIndexes(t *testing.T) {
	t.Parallel()
	for e := eSuccess; e < h3ErrorEnd; e++ {
		if isValidIndex(h3Index(e)) {
			t.Errorf("error code %d must not be a valid index", e)
		}
	}
}
