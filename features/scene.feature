Feature: Placeholder scene geometry
  The phase 1 scene is a sky, a ground band whose top edge is at world Y 900, and the marker:
  a 96x48 world-unit rectangle centred on the marker X, resting on the ground.

  @R-01 @T-01-07
  Scenario: Marker rectangle on an ultrawide screen
    Given a screen of 2560x1080 pixels
    And the marker is at world X 480
    When the scene shapes are computed
    Then the marker shape is at screen 752,852 with size 96x48
    And the ground shape is at screen 320,900 with size 1920x180

  @R-01 @T-01-07
  Scenario: Shapes scale with a 4K screen
    Given a screen of 3840x2160 pixels
    And the marker is at world X 480
    When the scene shapes are computed
    Then the marker shape is at screen 864,1704 with size 192x96
    And the ground shape is at screen 0,1800 with size 3840x360
