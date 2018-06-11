package main

import "testing"

func TestAllocate(t *testing.T) {
	table := []struct {
		n, bulk  uint64
		from, to uint64
		name     string
	}{
		{1, 500, 0, 499, "first possible n value"},
	}
	for _, row := range table {
		from, to := allocate(row.n, row.bulk)
		if from != row.from {
			t.Errorf("%s: expected from to be %d, got %d", row.name, row.from, from)
		}
		if to != row.to {
			t.Errorf("%s: expected to to be %d, got %d", row.name, row.to, to)
		}
	}
}

func TestAllocate0Panic(t *testing.T) {
	defer func() {
		switch x := recover().(type) {
		case string:
			s := "n to be 0 is not allowed"
			if x != s {
				t.Errorf("expected string panic to be %q, got %q", s, x)
			}
		default:
			t.Errorf("expected panic for n to be 0, got %v", x)
		}
	}()
	allocate(0, 500)
}
