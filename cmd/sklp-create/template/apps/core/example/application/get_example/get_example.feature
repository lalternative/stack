Feature: Get Example
  As a client of the core API
  I want to fetch an Example from the read model
  So that the CQRS read side answers queries from the projection

  Scenario: Fetch an Example that was created
    Given an Example named "alpha" was created
    When I fetch it by its ID
    Then the query should succeed
    And the returned Example should have name "alpha"

  Scenario: Read model reflects a rename
    Given an Example named "alpha" was created
    And it was renamed to "omega"
    When I fetch it by its ID
    Then the returned Example should have name "omega"

  Scenario: Fetch an unknown Example
    When I fetch an unknown Example
    Then the query should fail with a not-found error
