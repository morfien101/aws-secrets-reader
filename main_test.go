package main

import (
	"testing"
)

func TestPostProcessUpper(t *testing.T) {
	input := `{"testKey1":"value", "testkey2":"value"}`

	outputMap, err := postProcess(input, "", true)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}
	expected_keys := []string{"TESTKEY1", "TESTKEY2"}
	for _, key := range expected_keys {
		if value, ok := outputMap[key]; !ok {
			t.Logf("Did not find the uppercased key expected. Got: %s", value)
			t.Fail()
		}
	}
}

func TestPostProcessPrepend(t *testing.T) {
	input := `{"one":"value", "two":"value"}`

	outputMap, err := postProcess(input, "gopher_", false)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	expected_keys := []string{"gopher_one", "gopher_two"}
	for _, key := range expected_keys {
		if value, ok := outputMap[key]; !ok {
			t.Logf("Did not find the prefixed key expected. Got: %s", value)
			t.Fail()
		}
	}
}

func TestPostProcessPrependAndUpper(t *testing.T) {
	input := `{"one":"value", "two":"value"}`
	*flagPrependKeys = "gopher_"

	outputMap, err := postProcess(input, "gopher_", true)

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	expected_keys := []string{"GOPHER_ONE", "GOPHER_TWO"}
	for _, key := range expected_keys {
		if value, ok := outputMap[key]; !ok {
			t.Logf("Did not find the upper-cased and prefixed key expected. Got: %s", value)
			t.Fail()
		}
	}
}

func TestJSONFormatting(t *testing.T) {
	input := `{"one":"value$\"", "two":"value"}`
	*flagPrependKeys = "gopher_"

	outputMap, _ := postProcess(input, "gopher_", true)
	outputString, err := format(outputMap, "json")

	if err != nil {
		t.Logf("postProcess failed. Error: %s", err)
		t.FailNow()
	}

	t.Log("\n", outputString)
}
