package ask

import (
	"bytes"
	"testing"
)

func TestStr(t *testing.T) {
	var readBuffer bytes.Buffer

	question := Question{
		Description:     "Give me value",
		DefaultValue:    "yes",
		PossibleAnswers: []string{"yes", "no"},
		ReadDevice:      &readBuffer,
	}

	// test the default value
	readBuffer.Write([]byte("\n"))

	result := Str(question)

	if result != "yes" {
		t.Error("Default value was not accepted")
	}

	// test a non-default value is being accepted
	readBuffer.Write([]byte("no\n"))

	result = Str(question)

	if result != "no" {
		t.Error("Non-default value was not accepted")
	}

	// test giving non-acceptable answers to begin with and then non-default value
	readBuffer.Write([]byte("yas\ny\nn\nno"))

	result = Str(question)

	if result != "no" {
		t.Error("Non-possible answers were possibly accepted")
	}

	// test giving non-acceptable answers to begin with and then default value
	readBuffer.Write([]byte("yas\ny\nn\nnope\nyes"))

	result = Str(question)

	if result != "yes" {
		t.Error("Non-possible answers were possibly accepted")
	}

	// remove the default value now
	question = Question{
		Description:     "Give me value",
		PossibleAnswers: []string{"yes", "no"},
		ReadDevice:      &readBuffer,
	}

	// test giving newlines should make the last value the accepted one
	readBuffer.Write([]byte("\n\n\ny\nno"))
	result = Str(question)

	if result != "no" {
		t.Error("Last valid value was not the one received")
	}

	// remove the possible answers, leave the default
	question = Question{
		Description:  "Give me value",
		DefaultValue: "alex",
		ReadDevice:   &readBuffer,
	}

	readBuffer.Write([]byte("\n"))
	result = Str(question)

	if result != "alex" {
		t.Error("Result is not the default value")
	}

	readBuffer.Write([]byte("kostantine\n"))
	result = Str(question)

	if result != "kostantine" {
		t.Error("Inputted value not received as result")
	}

	// no the possible answers, no default
	question = Question{
		Description: "Give me value",
		ReadDevice:  &readBuffer,
	}

	readBuffer.Write([]byte("\n"))
	result = Str(question)

	if result != "" {
		t.Error("Inputted value not received as result")
	}

	readBuffer.Write([]byte("value\n"))
	result = Str(question)

	if result != "value" {
		t.Error("Inputted value not received as result")
	}
}
