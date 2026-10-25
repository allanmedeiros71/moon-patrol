Feature: Input mapping
  Raw keyboard and gamepad state is mapped to game actions without a window.
  Game actions: Left, Right, Jump, Fire. System actions: ToggleFullscreen, ToggleDebug, Quit.
  "Held" is the set of actions active this tick; "pressed" is the set that became active this tick.

  @R-03 @T-01-02
  Scenario Outline: Default keyboard bindings
    Given the default bindings
    When the key "<key>" is held
    Then the held actions are exactly "<actions>"

    Examples:
      | key        | actions |
      | ArrowLeft  | Left    |
      | ArrowRight | Right   |
      | ArrowUp    | Jump    |
      | Z          | Jump    |
      | X          | Fire    |
      | Space      | Fire    |

  @R-03 @T-01-02
  Scenario: An unbound key produces no action
    Given the default bindings
    When the key "A" is held
    Then the held actions are exactly ""

  @R-03 @T-01-02
  Scenario: Several keys held together produce every mapped action
    Given the default bindings
    When the keys "ArrowRight, Z, Space" are held
    Then the held actions are exactly "Right Jump Fire"

  @R-03 @T-01-02
  Scenario Outline: Default standard gamepad bindings
    Given the default bindings
    And a standard gamepad "Xbox Wireless Controller" is connected
    When the gamepad button "<button>" is held
    Then the held actions are exactly "<actions>"

    Examples:
      | button      | actions |
      | LeftLeft    | Left    |
      | LeftRight   | Right   |
      | LeftTop     | Jump    |
      | RightBottom | Jump    |
      | RightRight  | Fire    |
      | RightLeft   | Fire    |

  @R-03 @T-01-02
  Scenario Outline: Left stick beyond the dead zone steers
    Given the default bindings
    And a standard gamepad "Xbox Wireless Controller" is connected
    When the gamepad left stick horizontal axis is <value>
    Then the held actions are exactly "<actions>"

    Examples:
      | value | actions |
      | -0.8  | Left    |
      | 0.8   | Right   |
      | 0.3   |         |
      | -0.5  | Left    |
      | 0.49  |         |

  @R-03 @T-01-02
  Scenario: Keyboard and gamepad combine
    Given the default bindings
    And a standard gamepad "Xbox Wireless Controller" is connected
    When the key "Z" is held
    And the gamepad button "RightRight" is held
    Then the held actions are exactly "Jump Fire"

  @R-03 @T-01-02
  Scenario: A gamepad without the standard layout is ignored
    Given the default bindings
    And a non-standard gamepad "Generic USB Joystick" is connected
    When its raw button 0 is held
    Then the held actions are exactly ""
    And the gamepad "Generic USB Joystick" is reported as unsupported

  @R-03 @T-01-02
  Scenario: An action counts as pressed only on the tick it starts
    Given the default bindings
    When the key "X" is held for 3 ticks
    Then the pressed actions on tick 1 are "Fire"
    And the pressed actions on tick 2 are ""
    And the pressed actions on tick 3 are ""
    And the held actions on tick 3 are "Fire"

  @R-03 @T-01-02
  Scenario: Releasing and pressing again counts as a new press
    Given the default bindings
    When the key "X" is held for 1 tick, released for 1 tick and held for 1 tick
    Then the pressed actions on tick 3 are "Fire"

  @R-01 @R-03 @T-01-02
  Scenario Outline: System keys
    Given the default bindings
    When the keys "<keys>" are held
    Then the pressed actions are exactly "<actions>"

    Examples:
      | keys            | actions          |
      | F11             | ToggleFullscreen |
      | AltLeft, Enter  | ToggleFullscreen |
      | AltRight, Enter | ToggleFullscreen |
      | Enter           |                  |
      | F3              | ToggleDebug      |
      | Escape          | Quit             |
