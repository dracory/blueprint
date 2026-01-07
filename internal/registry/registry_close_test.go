package registry

import "testing"

func TestClose_NilReceiverDoesNotPanic(t *testing.T) {
	var r *registryImplementation
	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestClose_NilDatabaseDoesNotPanic(t *testing.T) {
	r := &registryImplementation{}
	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if err := r.Close(); err != nil {
		t.Fatalf("expected nil error on repeated close, got: %v", err)
	}
}
