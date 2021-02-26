package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/authelia/authelia/internal/configuration/schema"
	"github.com/authelia/authelia/internal/utils"
)

// ValidateSession validates and update session configuration.
func ValidateSession(configuration *schema.SessionConfiguration, validator *schema.StructValidator) {
	if configuration.Name == "" {
		configuration.Name = schema.DefaultSessionConfiguration.Name
	}

	if configuration.Redis != nil {
		if configuration.Redis.Sentinel != "" {
			validateRedisSentinel(configuration, validator)
		} else {
			validateRedis(configuration, validator)
		}
	}

	if configuration.Expiration == "" {
		configuration.Expiration = schema.DefaultSessionConfiguration.Expiration // 1 hour
	} else if _, err := utils.ParseDurationString(configuration.Expiration); err != nil {
		validator.Push(fmt.Errorf("Error occurred parsing session expiration string: %s", err))
	}

	if configuration.Inactivity == "" {
		configuration.Inactivity = schema.DefaultSessionConfiguration.Inactivity // 5 min
	} else if _, err := utils.ParseDurationString(configuration.Inactivity); err != nil {
		validator.Push(fmt.Errorf("Error occurred parsing session inactivity string: %s", err))
	}

	if configuration.RememberMeDuration == "" {
		configuration.RememberMeDuration = schema.DefaultSessionConfiguration.RememberMeDuration // 1 month
	} else if _, err := utils.ParseDurationString(configuration.RememberMeDuration); err != nil {
		validator.Push(fmt.Errorf("Error occurred parsing session remember_me_duration string: %s", err))
	}

	if configuration.Domain == "" {
		validator.Push(errors.New("Set domain of the session object"))
	}

	if strings.Contains(configuration.Domain, "*") {
		validator.Push(errors.New("The domain of the session must be the root domain you're protecting instead of a wildcard domain"))
	}
}

func validateRedis(configuration *schema.SessionConfiguration, validator *schema.StructValidator) {
	if configuration.Secret == "" {
		validator.Push(fmt.Errorf(errFmtSessionSecretRedisProvider, "redis"))
	}

	if !strings.HasPrefix(configuration.Redis.Host, "/") && configuration.Redis.Port == 0 {
		validator.Push(errors.New("A redis port different than 0 must be provided"))
	} else if configuration.Redis.Port <= -1 || configuration.Redis.Port > 65535 {
		validator.Push(fmt.Errorf(errFmtSessionRedisPortRange, "redis"))
	}
}

func validateRedisSentinel(configuration *schema.SessionConfiguration, validator *schema.StructValidator) {
	if configuration.Secret == "" {
		validator.Push(fmt.Errorf(errFmtSessionSecretRedisProvider, "redis sentinel"))
	}
}
