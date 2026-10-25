Feature: Fullscreen and window mode
  The game starts fullscreen unless started with -windowed. F11 and Alt+Enter toggle between
  fullscreen and a resizable 1280x720 window. A key held down toggles only once.

  @R-01 @T-01-05
  Scenario: Starts fullscreen by default
    Given the display controller is created with default options
    Then the display mode is fullscreen

  @R-01 @T-01-05
  Scenario: Starts windowed with the -windowed option
    Given the display controller is created with the windowed option
    Then the display mode is windowed with size 1280x720

  @R-01 @T-01-05
  Scenario: F11 toggles fullscreen off and on
    Given the display controller is created with default options
    When "ToggleFullscreen" is pressed
    Then the display mode is windowed with size 1280x720
    When "ToggleFullscreen" is pressed
    Then the display mode is fullscreen

  @R-01 @T-01-05
  Scenario: Holding the toggle key toggles only once
    Given the display controller is created with default options
    When the keys "F11" are held for 30 ticks
    Then the display mode is windowed with size 1280x720
