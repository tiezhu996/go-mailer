package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("MAILER_WORKERS", "")
	t.Setenv("MAILER_BATCH_SIZE", "")
	t.Setenv("MAILER_RETRIES", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 || c.Retries != 1 {
		t.Fatalf("defaults %+v", c)
	}
}

func TestLoadEnv(t *testing.T) {
	t.Setenv("MAILER_WORKERS", "5")
	t.Setenv("MAILER_BATCH_SIZE", "3")
	t.Setenv("MAILER_RETRIES", "2")
	c := Load()
	if c.Workers != 5 || c.BatchSize != 3 || c.Retries != 2 {
		t.Fatalf("env %+v", c)
	}
}
