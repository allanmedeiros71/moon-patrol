Feature: Game loop wiring
  The Ebitengine game polls input, maps it, steps the simulation exactly once per Update and
  applies display commands. It is driven headless with a fake window and fake input source.

  @R-02 @T-01-06
  Scenario: Each Update advances the simulation exactly one tick
    Given a headless game with seed 1 and a fake window of 1920x1080
    When Update is called 120 times
    Then the simulation tick count is 120

  @R-02 @T-01-06
  Scenario: The game runs logic at 60 ticks per second
    Given a headless game with seed 1 and a fake window of 1920x1080
    Then the configured ticks per second is 60

  @R-01 @T-01-06
  Scenario: Polled F11 switches the window out of fullscreen
    Given a headless game with seed 1 and a fake window of 1920x1080
    And the fake window is fullscreen
    When the fake input holds the key "F11" for 1 Update
    Then the fake window was set to windowed with size 1280x720

  @R-01 @T-01-06
  Scenario: Escape ends the game
    Given a headless game with seed 1 and a fake window of 1920x1080
    When the fake input holds the key "Escape" for 1 Update
    Then Update returns the termination signal

  @R-01 @R-02 @T-01-06
  Scenario: World state does not depend on the window size
    Given a headless game with seed 1 and a fake window of 800x600
    And a headless game with seed 1 and a fake window of 3840x2160
    When both games receive the key "ArrowRight" for 45 Updates
    Then both simulations have marker X 660.0

  @R-01 @T-01-06
  Scenario: Layout renders at native resolution
    Given a headless game with seed 1 and a fake window of 1280x720
    And the fake window device scale factor is 2.0
    When Layout is called with an outside size of 1280x720
    Then the layout size is 2560x1440
