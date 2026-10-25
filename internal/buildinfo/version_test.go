package buildinfo

import "testing"

func TestString(t *testing.T) {
	old := Version
	t.Cleanup(func() { Version = old })
	Version = "v1.2.3"
	if got, want := String(), "moonpatrol v1.2.3"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
