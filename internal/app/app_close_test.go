package app

import (
	"testing"

	"project/internal/config"

	_ "modernc.org/sqlite"
)

func TestClose_NilReceiverDoesNotPanic(t *testing.T) {
	var r *appImplementation
	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestClose_NilDatabaseDoesNotPanic(t *testing.T) {
	r := &appImplementation{}
	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error on repeated close, got: %v", err)
	}
}

func TestClose_ClosesNeatDatabase(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(":memory:")

	app, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	impl, ok := app.(*appImplementation)
	if !ok {
		t.Fatal("expected *appImplementation")
	}

	if impl.neatDB == nil {
		t.Fatal("expected neatDB to be set")
	}
	if impl.db == nil {
		t.Fatal("expected sql db to be set")
	}

	if err := app.Close(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if impl.neatDB != nil {
		t.Error("expected neatDB to be nil after close")
	}
	if impl.db != nil {
		t.Error("expected sql db to be nil after close")
	}
}
