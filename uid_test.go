package uid_test

import (
	"testing"

	"github.com/billglover/uid"
)

func TestUID(t *testing.T) {
	g := uid.NewGenerator()
	id, err := g.NextID()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(id) != 12 {
		t.Fatalf("uid should be 12 bytes, got %d", len(id))
	}
}

func TestUIDString(t *testing.T) {
	g := uid.NewGenerator()
	id, err := g.NextStringID()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == "" {
		t.Fatalf("uid should not be empty")
	}
}

func TestUnique(t *testing.T) {
	g := uid.NewGenerator()

	m := make(map[string]bool, 100000)

	for i := 0; i < 100000; i++ {
		id, err := g.NextStringID()
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
	g := uid.NewGenerator()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.NextID()
		}
	})
}

func BenchmarkUIDString(b *testing.B) {
	g := uid.NewGenerator()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.NextID()
		}
	})
}
