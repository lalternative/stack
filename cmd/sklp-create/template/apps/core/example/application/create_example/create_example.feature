Feature: Create Example
  As a client of the core API
  I want to create a new Example
  So that the CQRS write side records it as an event stream

  Scenario: Successfully create an Example
    When I create a new Example named "alpha"
    Then the creation should succeed
    And the Example should have name "alpha"

  Scenario: Fail to create an Example without a name
    When I try to create a new Example with no name
    Then the creation should fail with a validation error

  Scenario: Event sourcing verification for Example creation
    When I create a new Example named "beta"
    Then an "example.created" event should be recorded
    And the event should carry the Example ID
    And the event should carry the name "beta"
