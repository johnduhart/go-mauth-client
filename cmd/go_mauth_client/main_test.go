package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

// Confirm the IsNull function identifies an empty string and vice versa
func TestIsNull(t *testing.T) {
	test := ""
	expected := true
	actual := IsNull(&test)
	if actual != expected {
		t.Error("Failed with null String")
	}
	test = "Some Apples"
	expected = false
	actual = IsNull(&test)
	if actual != expected {
		t.Error("Failed with actual String")
	}
}

// Confirm that the CheckAction function identifies supported Verbs
func TestCheckAction(t *testing.T) {
	tests := map[string]bool{
		"GET":       true,
		"POST":      true,
		"DELETE":    true,
		"PUT":       true,
		"PINEAPPLE": false,
	}
	for verb, expected := range tests {
		actual := CheckAction(&verb)
		if actual != expected {
			t.Error("Failed with ", verb)
		}

	}

}

// Confirm that we can load a configuration file
func TestLoadMAuthConfig(t *testing.T) {
	client, _ := LoadMAuthConfig("/this/does/not/exist.txt")
	if client != nil {
		t.Error("Expected non-existing file to return nil")
	}
	client, _ = LoadMAuthConfig("test/test_config.json")
	if client == nil {
		t.Error("Expected existing file to return not nil")
	}
	if client.AppId != "11111111-2222-4105-b42e-88888888888" {
		t.Error("Incorrect APP ID")
	}
}

// Confirm that once we get a file or string we process it appropriately
func TestProcessConfiguration(t *testing.T) {
	var test_json string

	test_json = "{"
	_, err := ProcessConfiguration([]byte(test_json))
	if err == nil {
		t.Error("Expected failure with invalid JSON")
	}
	test_json = "{\"private_key_file\":\"test/private_key.pem\"}"
	_, err = ProcessConfiguration([]byte(test_json))
	if err == nil {
		t.Error("Expected failure with no app_uuid")
	}
	test_json = "{\"app_uuid\":\"11111111-2222-4105-b42e-88888888888\"}"
	_, err = ProcessConfiguration([]byte(test_json))
	if err == nil {
		t.Error("Expected failure with no private key details")
	}
	test_json = "{\"app_uuid\":\"11111111-2222-4105-b42e-88888888888\",\"private_key_file\":\"test/private_key.pem\"}"
	_, err = ProcessConfiguration([]byte(test_json))
	if err != nil {
		t.Error("Expected success with app_uuid and private_key_file")
	}
	content, _ := ioutil.ReadFile("test/private_key.pem")
	key_text := string(content)
	// escape the newlines
	key_content := strings.Replace(key_text, "\n", "\\n", -1)
	test_json = "{\"app_uuid\":\"11111111-2222-4105-b42e-88888888888\",\"private_key_text\":\"" + key_content + "\"}"
	_, err = ProcessConfiguration([]byte(test_json))
	if err != nil {
		t.Error("Expected success with app_uuid and private_key_text")
	}
}
