package config

import (
	"errors"
	"project/internal/types"

	"github.com/dracory/env"
)

// LoadV2 loads configuration using the metadata-driven variable definitions.
func LoadV2() (types.ConfigInterface, error) {
	cfg := types.Config{}

	vars := variableDefinitions()

	var err error

	for _, variable := range vars {
		switch variable.variableType {
		case variableTypeInt:
			err = processVariableInt(&cfg, variable)
		case variableTypeFloat:
			err = processVariableFloat(&cfg, variable)
		case variableTypeBool:
			err = processVariableBool(&cfg, variable)
		default:
			err = processVariableString(&cfg, variable)
		}

		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func processVariableBool(cfg types.ConfigInterface, variable variable) error {
	if variable.assignBool == nil {
		return errors.New(`variable ` + variable.key + ` is does not have an assignBool function`)
	}

	v := env.GetBool(variable.key)

	return variable.assignBool(cfg, v)
}

func processVariableInt(cfg types.ConfigInterface, variable variable) error {
	if variable.assignInt == nil {
		return errors.New(`variable ` + variable.key + ` is does not have an assignInt function`)
	}

	v := env.GetInt(variable.key)

	return variable.assignInt(cfg, v)
}

func processVariableFloat(cfg types.ConfigInterface, variable variable) error {
	if variable.assignFloat == nil {
		return errors.New(`variable ` + variable.key + ` is does not have an assignFloat function`)
	}

	v := env.GetFloat(variable.key)

	return variable.assignFloat(cfg, v)
}

func processVariableString(cfg types.ConfigInterface, variable variable) error {
	if variable.assignString == nil {
		return errors.New(`variable ` + variable.key + ` is does not have an assignString function`)
	}

	v := env.GetString(variable.key)

	return variable.assignString(cfg, v)
}
