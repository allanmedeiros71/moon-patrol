package render

import (
	"math"
	"testing"
)

func TestSceneRect(t *testing.T) {
	tests := []struct {
		name      string
		sw, sh    int
		want      Rect
		wantEmpty bool
	}{
		{"exact", 1920, 1080, Rect{0, 0, 1920, 1080}, false},
		{"wider", 2560, 1080, Rect{320, 0, 1920, 1080}, false},
		{"taller", 1000, 1000, Rect{0, 218, 1000, 563}, false},
		{"rounding", 1366, 768, Rect{0, 0, 1365, 768}, false},
		{"1x1", 1, 1, Rect{0, 0, 1, 1}, false},
		{"negative", -5, 100, Rect{}, true},
		{"zero height", 100, 0, Rect{}, true},
	}
	for _, tc := range tests {
		l := NewLetterbox(tc.sw, tc.sh, 1920, 1080)
		if l.Scene != tc.want || l.Scene.Empty() != tc.wantEmpty {
			t.Errorf("%s: scene = %+v, want %+v", tc.name, l.Scene, tc.want)
		}
		if l.Screen.W != tc.sw || l.Screen.H != tc.sh {
			t.Errorf("%s: screen = %+v", tc.name, l.Screen)
		}
	}
}

func TestMapping(t *testing.T) {
	l := NewLetterbox(2560, 1080, 1920, 1080)
	if x, y := l.ToScreen(1920, 1080); x != 2240 || y != 1080 {
		t.Errorf("ToScreen = %v,%v", x, y)
	}
	if x, y := l.ToWorld(320, 0); x != 0 || y != 0 {
		t.Errorf("ToWorld = %v,%v", x, y)
	}
	if s := l.Scale(); s != 1 {
		t.Errorf("Scale = %v", s)
	}
	if s := NewLetterbox(3840, 2160, 1920, 1080).Scale(); s != 2 {
		t.Errorf("Scale = %v", s)
	}
}

func TestRoundTrip(t *testing.T) {
	pts := [][2]float64{{0, 0}, {1920, 1080}, {960, 540}, {123.5, 987.25}, {1, 1079}}
	for _, sz := range [][2]int{{1366, 768}, {1000, 1000}, {2560, 1080}} {
		l := NewLetterbox(sz[0], sz[1], 1920, 1080)
		for _, p := range pts {
			sx, sy := l.ToScreen(p[0], p[1])
			wx, wy := l.ToWorld(sx, sy)
			if math.Abs(wx-p[0]) > 1e-9 || math.Abs(wy-p[1]) > 1e-9 {
				t.Errorf("%v: round trip %v -> %v,%v", sz, p, wx, wy)
			}
		}
	}
}
