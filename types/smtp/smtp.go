package smtp

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"strconv"
)

// Config is a type that defines required data for connecting to a SMTP server
type Config struct {
	from        string
	host        string // SMTP server domain
	port        int    // SMTP server port
	username    string // SMTP authentication username
	password    string // SMTP authentication password
	validateTLS bool   // Whether to validate TLS certificates
}

// NewConfig is the function that validates and returns Config instance
func NewConfig() (*Config, error) {
	config := Config{}

	// Check if all required environment variables are defined
	requiredEnvVariables := []string{"SMTP_FROM", "SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD"}

	for _, requiredEnvVariable := range requiredEnvVariables {
		if _, envVariableFound := os.LookupEnv(requiredEnvVariable); !envVariableFound {
			errorString := fmt.Sprintf("env variable \"%s\" must be set, cannot load smtp config", requiredEnvVariable)
			return nil, errors.New(errorString)
		}
	}

	// Parse SMTP port from string to integer
	port, portAtoiError := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if portAtoiError != nil {
		return nil, errors.New("failed to parse \"SMTP_PORT\" value")
	}

	if port <= 0 || port >= 65536 {
		return nil, errors.New("SMTP port value must be between 1 and 65535")
	} else {
		config.port = port
	}

	if _, err := mail.ParseAddress(os.Getenv("SMTP_FROM")); err != nil {
		return nil, errors.New("\"SMTP_FROM\" is not a valid email address")
	}

	// Load SMTP configuration from environment variables
	config.host = os.Getenv("SMTP_HOST")
	config.from = os.Getenv("SMTP_FROM")
	config.username = os.Getenv("SMTP_USERNAME")
	config.password = os.Getenv("SMTP_PASSWORD")

	// Check SMTP TLS validation setting (defaults to true if not specified)
	config.validateTLS = cmp.Or(os.Getenv("SMTP_VALIDATE_TLS"), "true") == "true"

	return &config, nil
}

func (config *Config) From() string {
	return config.from
}

func (config *Config) Host() string {
	return config.host
}

func (config *Config) Port() int {
	return config.port
}

func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.host, config.port)
}

func (config *Config) Username() string {
	return config.username
}

func (config *Config) Password() string {
	return config.password
}

func (config *Config) ValidateTLS() bool {
	return config.validateTLS
}

// LogValue allows to log Config masking password value
func (config Config) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("from", config.from),
		slog.String("address", config.Address()),
		slog.String("user", config.username),
		slog.String("password", "*****"),
		slog.Bool("validate_tls", config.validateTLS),
	}

	return slog.GroupValue(attrs...)
}
