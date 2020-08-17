package optionsvalidator

import (
	"errors"

	"github.com/hytromo/faulty-crane/internal/configuration"
)

// Validate ensures that the application options are valid and returns an error otherwise
func Validate(options configuration.AppOptions) error {
	if options.Configure.SubcommandEnabled {
		if options.Configure.Config == "" {
			return errors.New("Please specify a configuration file to save your answers to")
		}
	} else if options.Clean.SubcommandEnabled {
		if options.Clean.ContainerRegistry.Link == "" || options.Clean.ContainerRegistry.Access == "" {
			return errors.New("Please specify a valid container registry")
		}
	}

	return nil
}
