package render

import (
	"math"
	"testing"
)

func near(a, b FRect) bool {
	return math.Abs(a.X-b.X) < 1e-6 && math.Abs(a.Y-b.Y) < 1e-6 &&
		math.Abs(a.W-b.W) < 1e-6 && math.Abs(a.H-b.H) < 1e-6
}

func TestSceneShapesOrderAndColors(t *testing.T) {
	s := SceneShapes(NewLetterbox(1920, 1080, WorldW, WorldH), 480)
	wantColors := []interface{}{ColorSky, ColorGround, ColorGroundEdge, ColorMarker}
	if len(s) != 4 {
		t.Fatalf("len = %d", len(s))
	}
	for i, c := range wantColors {
		if s[i].Color != c {
			t.Errorf("shape %d colour = %v, want %v", i, s[i].Color, c)
		}
	}
	wantRects := []FRect{
		{0, 0, 1920, 1080}, {0, 900, 1920, 180}, {0, 900, 1920, 4}, {432, 852, 96, 48},
	}
	for i, r := range wantRects {
		if !near(s[i].Rect, r) {
			t.Errorf("shape %d rect = %+v, want %+v", i, s[i].Rect, r)
		}
	}
}

func TestSceneShapesUltrawideAnd4K(t *testing.T) {
	s := SceneShapes(NewLetterbox(2560, 1080, WorldW, WorldH), 480)
	if !near(s[3].Rect, FRect{752, 852, 96, 48}) || !near(s[1].Rect, FRect{320, 900, 1920, 180}) {
		t.Errorf("ultrawide: %+v %+v", s[3].Rect, s[1].Rect)
	}
	s = SceneShapes(NewLetterbox(3840, 2160, WorldW, WorldH), 480)
	if !near(s[3].Rect, FRect{864, 1704, 192, 96}) || !near(s[1].Rect, FRect{0, 1800, 3840, 360}) {
		t.Errorf("4K: %+v %+v", s[3].Rect, s[1].Rect)
	}
}

func TestSceneShapesBarsTopAndBottom(t *testing.T) {
	s := SceneShapes(NewLetterbox(1280, 1024, WorldW, WorldH), 480)
	sky := s[0].Rect
	if !near(sky, FRect{0, 152, 1280, 720}) {
		t.Errorf("sky = %+v", sky)
	}
}

func TestSceneShapesEmpty(t *testing.T) {
	if s := SceneShapes(NewLetterbox(0, 0, WorldW, WorldH), 480); len(s) != 0 {
		t.Errorf("got %d shapes", len(s))
	}
}
