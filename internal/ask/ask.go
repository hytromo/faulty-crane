package ask

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
)

// Question struct defines a standarized way to ask questions to the end user, by default reads from stdin
type Question struct {
	Description     string
	PossibleAnswers []string
	DefaultValue    string
	ReadDevice      io.Reader
}

// Str asks a question to the user expecting a string answer
func Str(question Question) string {
	readDevice := question.ReadDevice

	if readDevice == nil {
		readDevice = os.Stdin
	}

	finalText := question.Description

	hasSpecificPossibleAnswers := len(question.PossibleAnswers) != 0

	if hasSpecificPossibleAnswers {
		finalText = finalText + " (" + strings.Join(question.PossibleAnswers, "/") + ")"
	}

	if question.DefaultValue != "" {
		finalText = finalText + " [" + question.DefaultValue + "]"
	}

	finalText = finalText + ": "

	answer := ""
	scanner := bufio.NewScanner(readDevice)

	for {
		fmt.Print(finalText)
		scanner.Scan()

		answer = scanner.Text()

		if answer == "" {
			answer = question.DefaultValue
		}

		if !hasSpecificPossibleAnswers {
			return answer
		}

		if hasSpecificPossibleAnswers && stringutil.StrInSlice(answer, question.PossibleAnswers) {
			return answer
		}

		fmt.Printf("Please give a valid input, '%v' is not a valid value\n", answer)
	}
}
