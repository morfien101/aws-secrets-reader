package main

import (
	"strings"
	"testing"
)

func TestPostProcessUpper(t *testing.T) {
	input := `{"testKey1":"value", "testkey2":"value"}`

	output, err := postProcess(input, "", true)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	if !strings.Contains(output, "TESTKEY1") && !strings.Contains(output, "TESTKEY2") {
		t.Logf("Did not find the uppercased keys expected. Got: %s", output)
		t.Fail()
	}
}

func TestPostProcessPrepend(t *testing.T) {
	input := `{"one":"value", "two":"value"}`

	output, err := postProcess(input, "gopher_", false)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	if !strings.Contains(output, "gopher_one") && !strings.Contains(output, "gopher_two") {
		t.Logf("Did not find the gopher keys expected. Got: %s", output)
		t.Fail()
	}
}

func TestPostProcessPrependAndUpper(t *testing.T) {
	input := `{"one":"value", "two":"value"}`
	*flagPrependKeys = "gopher_"

	output, err := postProcess(input, "gopher_", true)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	if !strings.Contains(output, "GOPHER_ONE") && !strings.Contains(output, "GOPHER_TWO") {
		t.Logf("Did not find the gopher keys expected. Got: %s", output)
		t.Fail()
	}
}
