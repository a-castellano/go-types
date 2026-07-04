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
	"domain":      {VariableName: "SMTP_DOMAIN"},
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

	expectedError := "env variable \"SMTP_FROM\" must be set, cannot load smtp config instead"

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

func TestConfigWithValidConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test")
	os.Setenv("SMTP_DOMAIN", "test")
	os.Setenv("SMTP_HOST", "test")
	os.Setenv("SMTP_PORT", "25")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")

	config, err := NewConfig()

	if err != nil {
		t.Fatalf("TestConfigWithValidConfig should not fail, error was \"%s\"", err.Error())
	}
	if config.port != 25 {
		t.Fatalf("TestConfigWithValidConfig port should be 25, was \"%d\"", config.port)
	}
	if !config.validateTLS {
		t.Fatalf("TestConfigWithValidConfig validateTLS should be true")
	}

}

func TestConfigWithValidConfigNotTLSValidation(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test")
	os.Setenv("SMTP_DOMAIN", "test")
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

func TestLogValue(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("SMTP_FROM", "test")
	os.Setenv("SMTP_DOMAIN", "test")
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
