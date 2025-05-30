package config

import (
	"os"
	"testing"
)

func TestInitialize_AppServerHostAndPort(t *testing.T) {
	err := TestsConfigureAndInitialize()

	if err != nil {
		t.Fatal(err)
	}

	if WebServerHost == "" {
		t.Fatal("WebServerHost SHOULD NOT BE empty")
	}

	if WebServerPort == "" {
		t.Fatal("WebServerPort SHOULD NOT BE empty")
	}

	if AppUrl == "" {
		t.Fatal("AppUrl SHOULD NOT BE empty")
	}

	if DbDriver == "" {
		t.Fatal("DbDriver SHOULD NOT BE empty")
	}

	if DbHost != "" {
		t.Fatal("DbHost SHOULD BE empty")
	}

	if DbName == "" {
		t.Fatal("DbName SHOULD NOT BE empty")
	}

	if DbPass != "" {
		t.Fatal("DbPass SHOULD BE empty")
	}

	if DbPort != "" {
		t.Fatal("DbPort SHOULD BE empty")
	}
}

func TestInitialize_Debug(t *testing.T) {
	os.Setenv("DEBUG", "yes")
	err := TestsConfigureAndInitialize()

	if err != nil {
		t.Fatal(err)
	}

	if Debug == false {
		t.Fatal("Debug SHOULD NOT BE false")
	}
	if WebServerHost == "" {
		t.Fatal("ServerHost SHOULD NOT BE empty")
	}
	if WebServerPort == "" {
		t.Fatal("ServerPort SHOULD NOT BE empty")
	}
}
