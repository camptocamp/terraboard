package config

import (
	log "github.com/Sirupsen/logrus"
)

import "testing"

func TestSetLogging_debug(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.DebugLevel {
		t.Fatalf("Expected debug level, got %v", log.GetLevel())
	}
}

func TestSetLogging_info(t *testing.T) {
	c := Config{}
	c.Log.Level = "info"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.InfoLevel {
		t.Fatalf("Expected info level, got %v", log.GetLevel())
	}
}

func TestSetLogging_warn(t *testing.T) {
	c := Config{}
	c.Log.Level = "warn"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.WarnLevel {
		t.Fatalf("Expected warn level, got %v", log.GetLevel())
	}
}

func TestSetLogging_error(t *testing.T) {
	c := Config{}
	c.Log.Level = "error"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.ErrorLevel {
		t.Fatalf("Expected error level, got %v", log.GetLevel())
	}
}

func TestSetLogging_fatal(t *testing.T) {
	c := Config{}
	c.Log.Level = "fatal"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.FatalLevel {
		t.Fatalf("Expected fatal level, got %v", log.GetLevel())
	}
}

func TestSetLogging_panic(t *testing.T) {
	c := Config{}
	c.Log.Level = "panic"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if log.GetLevel() != log.PanicLevel {
		t.Fatalf("Expected panic level, got %v", log.GetLevel())
	}
}

func TestSetLogging_wronglevel(t *testing.T) {
	c := Config{}
	c.Log.Level = "wrong"
	c.Log.Format = "plain"
	err := c.SetupLogging()

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedError := "Wrong log level 'wrong'"

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}

func TestSetLogging_json(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "json"
	err := c.SetupLogging()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSetLogging_wrongformat(t *testing.T) {
	c := Config{}
	c.Log.Level = "debug"
	c.Log.Format = "yaml"
	err := c.SetupLogging()

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedError := "Wrong log format 'yaml'"

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}
