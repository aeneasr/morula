Feature: testing all subprojects

  As a developer
  I want to be able to run all tests irrespective of changes
  So that I can be sure the test status for my main development branch indicates broken subprojects that I didn't work on in my current change.

  - running "morula test --all" tests all subprojects, even if they don't contain changes


  Scenario: testing all subprojects
    Given a project with the subprojects "one" and "two"
    When running "morula test --all"
    Then it runs the tests:
      | one |
      | two |
