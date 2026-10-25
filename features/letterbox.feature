Feature: Resolution-independent 16:9 letterbox
  The world is 1920x1080 units. The scene is the largest 16:9 rectangle that fits the screen,
  centred; the rest is black bars. Scene width and height are rounded to the nearest pixel,
  offsets are (screen - scene) / 2 rounded down.

  @R-01 @T-01-04
  Scenario Outline: Scene rectangle for a screen size
    Given a screen of <sw>x<sh> pixels
    When the letterbox is computed
    Then the scene rectangle is at <x>,<y> with size <w>x<h>

    Examples:
      | sw   | sh   | x   | y   | w    | h    |
      | 1920 | 1080 | 0   | 0   | 1920 | 1080 |
      | 3840 | 2160 | 0   | 0   | 3840 | 2160 |
      | 2560 | 1080 | 320 | 0   | 1920 | 1080 |
      | 1280 | 1024 | 0   | 152 | 1280 | 720  |
      | 800  | 600  | 0   | 75  | 800  | 450  |
      | 1366 | 768  | 0   | 0   | 1365 | 768  |
      | 1000 | 1000 | 0   | 218 | 1000 | 563  |

  @R-01 @T-01-04
  Scenario Outline: World points map into the scene rectangle
    Given a screen of <sw>x<sh> pixels
    When the world point <wx>,<wy> is mapped to the screen
    Then the screen point is <px>,<py>

    Examples:
      | sw   | sh   | wx   | wy   | px   | py   |
      | 2560 | 1080 | 0    | 0    | 320  | 0    |
      | 2560 | 1080 | 1920 | 1080 | 2240 | 1080 |
      | 1280 | 1024 | 960  | 540  | 640  | 512  |
      | 3840 | 2160 | 480  | 900  | 960  | 1800 |

  @R-01 @T-01-04
  Scenario: Screen points map back to the same world point
    Given a screen of 1280x1024 pixels
    When the screen point 640,512 is mapped to the world
    Then the world point is 960.0,540.0

  @R-01 @T-01-04
  Scenario: A zero-sized screen yields an empty scene
    Given a screen of 0x0 pixels
    When the letterbox is computed
    Then the scene rectangle is empty
