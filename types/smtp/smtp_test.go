//go:build unit_tests || smtp_unit_tests

package smtp

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

type envVariable struct {
	Value        string
	IsDefined    bool
	VariableName string
}

var envVariables = map[string]envVariable{
	"from":        {VariableName: "SMTP_FROM"},
	"host":        {VariableName: "SMTP_HOST"},
	"port":        {VariableName: "SMTP_PORT"},
	"username":    {VariableName: "SMTP_USERNAME"},
	"password":    {VariableName: "SMTP_PASSWORD"},
	"validateTLS": {VariableName: "SMTP_VALIDATE_TLS"},
}

func setUp() {

	for key, variable := range envVariables {

		if envValue, found := os.LookupEnv(variable.VariableName); found {
			variable.Value = envValue
			variable.IsDefined = true
		} else {
			variable.IsDefined = false
		}

		os.Unsetenv(variable.VariableName)

		envVariables[key] = variable
	}

}

func teardown() {

	for _, variable := range envVariables {
		if variable.IsDefined {
			os.Setenv(variable.VariableName, variable.Value)
		} else {
			os.Unsetenv(variable.VariableName)
		}
	}
}

func TestConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestConfigWithoutEnvVariables should fail.")
	}

	expectedError := "env variable \"SMTP_FROM\" must be set, cannot load smtp config"

	if err.Error() != expectedError {
		t.Fatalf("TestConfigWithoutEnvVariables error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestConfigWithInvalidSMTPPort(t *testing.T) {

	setUp()
	defer teardown()

	// Set all required variables but with invalid SMTP port
	os.Setenv("SMTP_FROM", "test")
	os.Setenv("SMTP_DOMAIN", "test")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "invalid")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestConfigWithInvalidSMTPPORT should fail.")
	}
	expectedError := "failed to parse \"SMTP_PORT\" value"
	if err.Error() != expectedError {
		t.Errorf("TestConfigWithInvalidSMTPPort error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

// TestConfigWithOutOfRangePort checks that ports outside the 1-65535 range
// are rejected even when they parse as integers.
// This test was written by an AI agent (Claude).
func TestConfigWithOutOfRangePort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test")
	os.Setenv("SMTP_DOMAIN", "test")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "99999")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestConfigWithOutOfRangePort should fail.")
	}
	expectedError := "SMTP port value must be between 1 and 65535"
	if err.Error() != expectedError {
		t.Errorf("TestConfigWithOutOfRangePort error should be \"%s\" but it was \"%s\".", expectedError, err.Error())
	}
}

func TestConfigWithValidConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test@example.com")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USERNAME", "user")
	os.Setenv("SMTP_PASSWORD", "pass")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestConfigWithValidConfig should not fail, error was \"%s\"", err.Error())
	}
	expectedFrom := "test@example.com"
	if config.From() != expectedFrom {
		t.Fatalf("TestConfigWithValidConfig from should be \"%s\", was \"%s\"", expectedFrom, config.From())
	}

	expectedHost := "test"
	if config.Host() != expectedHost {
		t.Fatalf("TestConfigWithValidConfig host should be \"%s\", was \"%s\"", expectedHost, config.Host())
	}

	expectedPort := 25
	if config.Port() != expectedPort {
		t.Fatalf("TestConfigWithValidConfig port should be \"%d\", was \"%d\"", expectedPort, config.Port())
	}

	expectedAddress := "test:25"
	if config.Address() != expectedAddress {
		t.Fatalf("TestConfigWithValidConfig address should be \"%s\", was \"%s\"", expectedAddress, config.Address())
	}

	expectedUser := "user"
	if config.Username() != expectedUser {
		t.Fatalf("TestConfigWithValidConfig smtp user should be \"%s\", was \"%s\"", expectedUser, config.Username())
	}

	expectedPassword := "pass"
	if config.Password() != expectedPassword {
		t.Fatalf("TestConfigWithValidConfig smtp password should be \"%s\", was \"%s\"", expectedPassword, config.Password())
	}

	if !config.ValidateTLS() {
		t.Fatalf("TestConfigWithValidConfig validateTLS should be true")
	}

}

func TestConfigWithValidConfigNotTLSValidation(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test@example.com")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")
	os.Setenv("SMTP_VALIDATE_TLS", "anyvaluedifferentfromtrue")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestConfigWithValidConfigNotTLSValidation should not fail, error was \"%s\"", err.Error())
	}
	if config.validateTLS {
		t.Fatalf("TestConfigWithValidConfig validateTLS should be false")
	}

}

func TestConfigInvalidFrom(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "thisisnotavalidaddress")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")
	os.Setenv("SMTP_VALIDATE_TLS", "anyvaluedifferentfromtrue")

	_, err := NewConfig()

	if err == nil {
		t.Fatalf("TestConfigInvalidFrom should fail")
	}

}

func TestLogValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test@example.com")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")
	os.Setenv("SMTP_VALIDATE_TLS", "anyvaluedifferentfromtrue")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestLogValue should not fail when config is set, error was \"%s\"", err.Error())
	}
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	logger.Info("test log", "smtp config", config)

	bufferLen := buf.Len()

	if bufferLen <= 0 {
		t.Fatalf("TestLogValue has failed, buffer is empty")
	}

	var loggedData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &loggedData); err != nil {
		t.Fatalf("TestLogValue has failed, cannot unmarshal json log")
	}

	smtpConfig := loggedData["smtp config"].(map[string]interface{})
	passwordValue := smtpConfig["password"].(string)
	if passwordValue != "*****" {
		t.Fatalf("TestLogValue has failed, password should be \"*****\" but it was \"%s\"", passwordValue)
	}
}
