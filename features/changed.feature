Feature: running a command only in updated subprojects

  As a developer
  I want to be able to run a command only in subprojects that I have made modifications to
  So that I can test my changes as quickly as possible.

  - "morula run --updated <command>" runs the given command in all subprojects that contain changes


  Scenario: all tested subprojects are passing
    Given a project with the subprojects:
      | NAME  | TEMPLATE  |
      | one   | passing_1 |
      | two   | failing   |
      | three | passing_2 |
    And I am on the "feature" branch
    And subprojects "one" and "three" have changes
    When running "morula changed bin/spec"
    Then it runs that command in the directories:
      | one   |
      | three |


  Scenario: calling a command with command-line arguments
    Given a project with the subprojects:
      | NAME  | TEMPLATE  |
      | one   | passing_1 |
      | two   | failing   |
      | three | passing_2 |
    And I am on the "feature" branch
    And subprojects "one" and "three" have changes
    When running "morula changed 'ls -la'"
    Then it runs that command in the directories:
      | one   |
      | three |


  Scenario: some tested subprojects are failing
    Given a project with the subprojects:
      | NAME  | TEMPLATE  |
      | one   | passing_1 |
      | two   | failing   |
      | three | passing_2 |
    And I am on the "feature" branch
    And subprojects "one" and "two" have changes
    When trying to run "morula changed bin/spec"
    Then it fails with an error code and the message:
      """
      subproject two is broken
      """

  Scenario: forgetting to provide the command
    Given a project with the subprojects:
      | NAME  | TEMPLATE|
      | works | passing |
    When trying to run "morula changed"
    Then it fails with an error code and the message:
      """
      Please provide the command to run
      """


  Scenario: providing a command that doesn't exist
    Given a project with the subprojects:
      | NAME  | TEMPLATE|
      | works | passing |
    And I am on the "feature" branch
    And subproject "works" has changes
    When trying to run "morula changed zonk"
    Then it fails with an error code and the message:
      """
      command zonk doesn't exist
      """
