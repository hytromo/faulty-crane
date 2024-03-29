package optionsvalidator

import (
	"errors"

	"github.com/hytromo/faulty-crane/internal/configuration"
)

// Validate ensures that the application options are valid and returns an error otherwise
func Validate(options configuration.AppOptions) error {
	if options.Configure.SubcommandEnabled {
		if options.Configure.Config == "" {
			return errors.New("please specify a configuration file to save your answers to")
		}
	} else if options.Apply.SubcommandEnabled {
		if configuration.IsGCR(&options) {
			if options.ApplyPlanCommon.GoogleContainerRegistry.Host == "" || options.ApplyPlanCommon.GoogleContainerRegistry.Token == "" {
				return errors.New("please specify a valid container registry and access key for GCR")
			}
		} else if configuration.IsDockerhub(&options) {
			if options.ApplyPlanCommon.DockerhubContainerRegistry.Namespace == "" || options.ApplyPlanCommon.DockerhubContainerRegistry.Password == "" || options.ApplyPlanCommon.DockerhubContainerRegistry.Username == "" {
				return errors.New("please specify a valid namespace, username and password for Dockerhub")
			}
		}
	}

	return nil
}
