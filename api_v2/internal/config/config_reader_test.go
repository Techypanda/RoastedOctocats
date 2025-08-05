package config

import (
	"testing"
)

func TestConfigReader(t *testing.T) {
	localConfig := New().ReadConfiguration("local")
	if localConfig.GetTestPem() == "" {
		t.Fatal("expected a test pem in local")
	}
	if localConfig.GetPAT() == "" {
		t.Fatal("expected a test pat in local")
	}
}
