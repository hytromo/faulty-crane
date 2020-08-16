package configurationhelper

import (
	"fmt"

	"github.com/hytromo/faulty-crane/internal/argsparser"
)

func CreateNewConfiguration(args argsparser.ConfigureCliOptions) {
	fmt.Printf("Creating new configuration with %+v\n", args)
}
