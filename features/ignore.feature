Feature: ignoring hidden folders

  As a developer
  I want that Morula ignores directories starting with a dot
  So that I can use tools that create hidden folders in my code's directory.

  - all directories starting with a dot are ignored


  Background:
    Given a project with the subprojects:
      | NAME | TEMPLATE |
      | one  | passing  |
      | .two | passing  |


  Scenario: running "morula all"
    When running "morula all bin/spec"
    Then it runs that command in the directories:
      | one |


  Scenario: running "morula changed"
    And I am on the "feature" branch
    And subprojects "one" and ".two" have changes
    When running "morula changed bin/spec"
    Then it runs that command in the directories:
      | one |
