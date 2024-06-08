//go:build integration_tests || unit_tests || rabbitmq_tests

package rabbitmq

import "testing"

func TestRabbitmqConfigWithoutVariables(t *testing.T) {

	config, err := NewRabbitmqConfig()

	if err != nil {
		t.Errorf("NewRabbitmqConfig method without any env varible suited shouldn't fail.")
	} else {
		if config.Host != "localhost" {
			t.Errorf("RabbitmqConfig Host should be \"localhost.\" but was \"%s\".", config.Host)
		}
	}

}
