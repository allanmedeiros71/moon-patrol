Feature: Fixed-step headless simulation
  The simulation advances one tick per step at 60 ticks per second, never reads the clock,
  and is deterministic for a given seed and input sequence.
  Phase 1 placeholder world: a marker on the ground at world X 480 that moves 4 world units per
  tick while Left or Right is held (clamped to 0..1920), and a noise value drawn from the
  seeded random generator every tick.

  @R-02 @T-01-03
  Scenario: Stepping counts ticks and simulated time
    Given a new simulation with seed 42
    When it is stepped 600 ticks with no input
    Then the tick count is 600
    And the simulated time is 10.000 seconds

  @R-02 @T-01-03
  Scenario: Movement depends on ticks, not on frames or wall-clock time
    Given a new simulation with seed 42
    When it is stepped 60 ticks holding "Right"
    Then the marker X is 720.0

  @R-02 @T-01-03
  Scenario: Holding Left and Right together leaves the marker in place
    Given a new simulation with seed 42
    When it is stepped 30 ticks holding "Left Right"
    Then the marker X is 480.0

  @R-02 @T-01-03
  Scenario: The marker stays inside the world
    Given a new simulation with seed 42
    When it is stepped 200 ticks holding "Left"
    Then the marker X is 0.0
    When it is stepped 600 ticks holding "Right"
    Then the marker X is 1920.0

  @R-02 @T-01-03
  Scenario: Same seed and same input give the same world
    Given two new simulations with seed 42
    When both are stepped with the script "30 Right, 10 none, 20 Left Fire, 15 Jump"
    Then both world snapshots are identical

  @R-02 @T-01-03
  Scenario: A different seed gives a different world
    Given a new simulation with seed 42 and a new simulation with seed 43
    When both are stepped with the script "30 Right, 10 none"
    Then the world snapshots differ in the noise value only
