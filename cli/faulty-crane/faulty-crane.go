package main

import (
	"fmt"
	"os"

	"github.com/hytromo/faulty-crane/internal/argsParser"
)

func main() {
	fmt.Println("New version v5")
	cliOptions := argsParser.Parse(os.Args)
	fmt.Printf("parsed cli options is: %+v\n", cliOptions)
}
