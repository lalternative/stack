Feature: Rename Example
  As a client of the core API
  I want to rename an existing Example
  So that the CQRS write side appends a rename event to its stream

  Background:
    Given an Example named "alpha" exists

  Scenario: Successfully rename an existing Example
    When I rename it to "omega"
    Then the rename should succeed
    And the Example should have name "omega"

  Scenario: Fail to rename with an empty name
    When I try to rename it to ""
    Then the rename should fail with a validation error

  Scenario: Fail to rename an unknown Example
    When I try to rename an unknown Example
    Then the rename should fail with a not-found error

  Scenario: Event sourcing verification for Example rename
    When I rename it to "omega"
    Then an "example.renamed" event should be recorded
