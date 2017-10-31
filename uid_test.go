package uid_test

import (
	"testing"

	"github.com/billglover/uid"
)

func TestUID(t *testing.T) {
	id, err := uid.NextID()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id <= 0 {
		t.Fatalf("uid should be a positive number, got %d", id)
	}
}

func TestUIDString(t *testing.T) {
	id, err := uid.NextStringID()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == "" {
		t.Fatalf("uid should not be empty")
	}
}

func TestUnique(t *testing.T) {
	m := make(map[uint64]bool, 100000)

	for i := 0; i < 100000; i++ {
		id, err := uid.NextID()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if m[id] {
			t.Fatalf("uid %d reissued after %d robots.", id, i)
		}
		m[id] = true
	}
}

func BenchmarkUID(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = uid.NextID()
		}
	})
}

func BenchmarkUIDString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = uid.NextID()
		}
	})
}
