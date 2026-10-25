Feature: Debug overlay
  A text overlay shows timing, sizes, held actions, the marker and connected gamepads.
  It is visible at start in phase 1 and F3 toggles it.

  @R-03 @T-01-05
  Scenario: Overlay is visible at start and F3 toggles it
    Given the display controller is created with default options
    Then the debug overlay is visible
    When "ToggleDebug" is pressed
    Then the debug overlay is hidden
    When "ToggleDebug" is pressed
    Then the debug overlay is visible

  @R-01 @R-03 @T-01-05
  Scenario: Overlay lines
    Given the overlay info:
      | tps      | 60.0                     |
      | fps      | 59.94                    |
      | tick     | 1234                     |
      | screen   | 2560x1080                |
      | held     | Right Fire               |
      | marker   | 480                      |
      | gamepads | Xbox Wireless Controller (standard) |
    When the overlay lines are built
    Then the overlay lines are:
      | TPS 60.0  FPS 59.9                                     |
      | Tick 1234                                              |
      | Window 2560x1080  Scene 1920x1080 at 320,0             |
      | Actions: Right Fire                                    |
      | Marker X 480.0                                         |
      | Gamepads: 1                                            |
      | Pad 0: Xbox Wireless Controller (standard)             |
      | F3: hide overlay  F11/Alt+Enter: fullscreen  Esc: quit |

  @R-03 @T-01-05
  Scenario: Overlay with no actions and an unsupported gamepad
    Given the overlay info:
      | tps      | 60.0                            |
      | fps      | 60.0                            |
      | tick     | 0                               |
      | screen   | 1920x1080                       |
      | held     |                                 |
      | marker   | 480                             |
      | gamepads | Generic USB Joystick (unsupported) |
    When the overlay lines are built
    Then overlay line 4 is "Actions: -"
    And overlay line 7 is "Pad 0: Generic USB Joystick (unsupported)"

  @R-03 @T-01-05
  Scenario: Overlay with no gamepads
    Given the overlay info:
      | tps      | 60.0      |
      | fps      | 60.0      |
      | tick     | 0         |
      | screen   | 1920x1080 |
      | held     |           |
      | marker   | 480       |
      | gamepads |           |
    When the overlay lines are built
    Then overlay line 6 is "Gamepads: 0"
    And overlay line 7 is "F3: hide overlay  F11/Alt+Enter: fullscreen  Esc: quit"
