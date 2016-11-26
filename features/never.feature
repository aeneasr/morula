Feature: never testing certain subprojects

  As a developer
  I want Morula to ignore folders that are not subprojects
  So that I don't get misleading warnings when trying to run tasks meant for subprojects in them.

  - subprojects listed under the configuration key "never" are ignored,
    even if they contain changes


  Scenario Outline: valid "never" entry in configuration file
    Given a project with the subprojects "one", "assets", and the configuration file:
      """
      never: assets
      """
    And I am on the "feature" branch
    And subproject "one" has changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | one |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "never" entry via command-line parameter
    Given a project with the subprojects "one", "assets", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subproject "one" has changes
    When running "morula <COMMAND> --never=assets bin/spec"
    Then it runs that command in the directories:
      | one |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: non-existing "never" entry
    Given a project with the subprojects "one", "assets", and the configuration file:
      """
      never: zonk
      """
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to never run subproject zonk,
      but such a subproject does not exist.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: "never" entry points to a file
    Given a project with the configuration file:
      """
      never: myfile
      """
    And the project contains a file "myfile"
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to never run subproject myfile,
      but this path is not a directory.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |
