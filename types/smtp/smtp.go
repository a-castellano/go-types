package smtp

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	from        string // Sender email username (without domain)
	domain      string // Sender email domain
	host        string // SMTP server hostname
	port        int    // SMTP server port
	username    string // SMTP authentication username
	password    string // SMTP authentication password
	validateTLS bool   // Whether to validate TLS certificates
}

func NewConfig() (*Config, error) {
	config := Config{}

	// Check if all required environment variables are defined
	requiredEnvVariables := []string{"SMTP_FROM", "SMTP_DOMAIN", "SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD"}

	for _, requiredEnvVariable := range requiredEnvVariables {
		if _, envVariableFound := os.LookupEnv(requiredEnvVariable); !envVariableFound {
			errorString := fmt.Sprintf("env variable \"%s\" must be set, cannot load smtp config instead", requiredEnvVariable)
			return nil, errors.New(errorString)
		}
	}

	// Parse SMTP port from string to integer
	port, portAtoiError := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if portAtoiError != nil {
		return nil, errors.New("failed to parse \"SMTP_PORT\" value")
	}

	config.port = port

	// Load SMTP configuration from environment variables
	config.from = os.Getenv("SMTP_FROM")
	config.domain = os.Getenv("SMTP_DOMAIN")
	config.host = os.Getenv("SMTP_HOST")
	config.username = os.Getenv("SMTP_USERNAME")
	config.password = os.Getenv("SMTP_PASSWORD")

	// Check SMTP TLS validation setting (defaults to true if not specified)
	config.validateTLS = cmp.Or(os.Getenv("SMTP_VALIDATE_TLS"), "true") == "true"

	return &config, nil
}

func (config Config) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("from", config.from),
		slog.String("domain", config.domain),
		slog.String("host", config.host),
		slog.String("user", config.username),
		slog.String("password", "*****"),
	}

	return slog.GroupValue(attrs...)
}
