package uid_test

import (
	"os"
	"testing"

	"github.com/billglover/uid"
)

func TestUID(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("unable to get hostname: %v", err)
	}
	pid := os.Getpid()

	id, err := uid.NextID(hostname, pid)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(id) != 12 {
		t.Fatalf("uid should be 12 bytes, got %d", len(id))
	}
}

func TestUIDString(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("unable to get hostname: %v", err)
	}
	pid := os.Getpid()

	id, err := uid.NextStringID(hostname, pid)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == "" {
		t.Fatalf("uid should not be empty")
	}
}

func TestUnique(t *testing.T) {

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("unable to get hostname: %v", err)
	}
	pid := os.Getpid()

	m := make(map[string]bool, 100000)

	for i := 0; i < 100000; i++ {
		id, err := uid.NextStringID(hostname, pid)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if m[id] {
			t.Fatalf("uid %s reissued after %d robots.", id, i)
		}
		m[id] = true
	}
}

func BenchmarkUID(b *testing.B) {
	hostname, err := os.Hostname()
	if err != nil {
		b.Fatalf("unable to get hostname: %v", err)
	}
	pid := os.Getpid()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = uid.NextID(hostname, pid)
		}
	})
}

func BenchmarkUIDString(b *testing.B) {
	hostname, err := os.Hostname()
	if err != nil {
		b.Fatalf("unable to get hostname: %v", err)
	}
	pid := os.Getpid()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = uid.NextStringID(hostname, pid)
		}
	})
}
