package inputebiten

import (
	"testing"

	"github.com/allanmedeiros/moonpatrol/internal/input"
)

func TestEbitenTablesComplete(t *testing.T) {
	nKeys := int(input.KeyEscape) + 1
	if len(ebitenKeys) != nKeys {
		t.Errorf("ebitenKeys has %d entries, want %d", len(ebitenKeys), nKeys)
	}
	for i := range nKeys {
		if _, ok := ebitenKeys[input.Key(i)]; !ok {
			t.Errorf("key %d not translated", i)
		}
	}
	nPads := int(input.CenterCenter) + 1
	if len(ebitenPadButtons) != nPads {
		t.Errorf("ebitenPadButtons has %d entries, want %d", len(ebitenPadButtons), nPads)
	}
	seen := map[int]bool{}
	for b, eb := range ebitenPadButtons {
		if seen[int(eb)] {
			t.Errorf("duplicate ebiten button for %d", b)
		}
		seen[int(eb)] = true
	}
}
