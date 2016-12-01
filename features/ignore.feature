Feature: ignoring hidden folders

  As a developer
  I want that Morula ignores directories starting with a dot
  So that I can use tools that create hidden folders in my code's directory.

  - all directories starting with a dot are ignored


  Scenario Outline: the workspace contains directories with a dot
    Given a project with the subprojects:
      | NAME | TEMPLATE |
      | one  | passing  |
      | .two | passing  |
    And I am on the "feature" branch
    And subprojects "one" and ".two" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | one |

    Examples:
      | COMMAND |
      | all     |
      | changed |
