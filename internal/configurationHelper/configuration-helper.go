package configurationHelper

import (
	"fmt"

	"github.com/hytromo/faulty-crane/internal/argsParser"
)

func CreateNewConfiguration(args argsParser.ConfigureCliOptions) {
	fmt.Printf("Creating new configuration with %+v\n", args)
}
