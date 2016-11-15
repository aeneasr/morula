Feature: testing all subprojects

  As a developer
  I want to be able to run all tests irrespective of changes
  So that I can be sure the test status for my main development branch indicates broken subprojects that I didn't work on in my current change.

  - running "morula test --all" tests all subprojects, even if they don't contain changes


  Scenario: all subprojects are passing
    Given a project with the subprojects:
      | NAME | TEMPLATE  |
      | one  | passing_1 |
      | two  | passing_2 |
    When running "morula test --all"
    Then it runs the tests:
      | one |
      | two |


  Scenario: some subprojects are failing
    Given a project with the subprojects:
      | NAME  | TEMPLATE|
      | works | passing |
      | fails | failing |
    When trying to run "morula test --all"
    Then it fails with an error code and the message:
      """
      subproject fails is broken
      """
