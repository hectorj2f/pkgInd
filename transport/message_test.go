package transport

import (
	"log"
	"reflect"
	"testing"
)

func TestMessageFormat(t *testing.T) {
	rawMsg := "INDEX|cloog|gmp,isl,pkg-config\n"
	if err := ValidateMessage(rawMsg); err != nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}

	rawMsg = "INDEX|cloog\n"
	if err := ValidateMessage(rawMsg); err == nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}

	rawMsg = "INDEX|cloog|\n"
	if err := ValidateMessage(rawMsg); err != nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}

	rawMsg = "XXINDEX|cloog|gmp,isl,pkg-config\n"
	if err := ValidateMessage(rawMsg); err == nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}

	rawMsg = "XXINDEX||gmp,isl,pkg-config\n"
	if err := ValidateMessage(rawMsg); err == nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}

	rawMsg = "XXINDEX||gmp,isl,pkg-config"
	if err := ValidateMessage(rawMsg); err == nil {
		log.Fatalf("unable to validate the message %s: %v", rawMsg, err)
	}
}

func TestMessageExtraction(t *testing.T) {
	rawMsg := "INDEX|cloog|gmp,isl,pkg-config\n"

	msg, err := ExtractMessage(rawMsg)
	if err != nil {
		log.Fatalf("Unexpected error when extracting the message: %v", err)
	}
	expected := &MessageRequest{
		Command:      "INDEX",
		Package:      "cloog",
		Dependencies: []string{"gmp", "isl", "pkg-config"},
	}

	if reflect.DeepEqual(msg, expected) {
		log.Fatalf("message is wrong %v expected %v", msg, expected)
	}

	rawMsg = "QUERY|cloog|\n"
	msg, err = ExtractMessage(rawMsg)
	expected = &MessageRequest{
		Command:      "QUERY",
		Package:      "cloog",
		Dependencies: []string{},
	}
	if reflect.DeepEqual(msg, expected) {
		log.Fatalf("message is wrong %v expected %v", msg, expected)
	}

	rawMsg = "INDEX|cloog|gmp\n"
	msg, err = ExtractMessage(rawMsg)
	expected = &MessageRequest{
		Command:      "INDEX",
		Package:      "cloog",
		Dependencies: []string{"gmp"},
	}
	if reflect.DeepEqual(msg, expected) {
		log.Fatalf("message is wrong %v expected %v", msg, expected)
	}
}
