package input

import "testing"

func TestEbitenTablesComplete(t *testing.T) {
	if len(ebitenKeys) != len(keyNames) {
		t.Errorf("ebitenKeys has %d entries, want %d", len(ebitenKeys), len(keyNames))
	}
	for i := range keyNames {
		if _, ok := ebitenKeys[Key(i)]; !ok {
			t.Errorf("key %s not translated", keyNames[i])
		}
	}
	if len(ebitenPadButtons) != len(padButtonNames) {
		t.Errorf("ebitenPadButtons has %d entries, want %d", len(ebitenPadButtons), len(padButtonNames))
	}
	seen := map[int]bool{}
	for b, eb := range ebitenPadButtons {
		if seen[int(eb)] {
			t.Errorf("duplicate ebiten button for %d", b)
		}
		seen[int(eb)] = true
	}
}
