Feature: always testing certain subprojects

  As a developer
  I want Morula to always run my end-to-end tests even if they don't contain any changes
  So that I can catch bugs in other subprojects that break the product.

  - subprojects listed under the configuration key "always" are always tested,
    even if they don't contain changes


  Scenario Outline: valid "always" entry in configuration file
    Given a project with the subprojects "one", "e2e", and the configuration file:
      """
      always: e2e
      """
    And I am on the "feature" branch
    And subproject "one" has changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | e2e |
      | one |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "always" entry via command-line parameter
    Given a project with the subprojects "one", "e2e", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subproject "one" has changes
    When running "morula <COMMAND> --always=e2e bin/spec"
    Then it runs that command in the directories:
      | e2e |
      | one |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: non-existing "always" entry
    Given a project with the subprojects "one", "e2e", and the configuration file:
      """
      always: zonk
      """
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to always run subproject zonk,
      but such a subproject does not exist.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: "always" entry points to a file
    Given a project with the configuration file:
      """
      always: myfile
      """
    And the project contains a file "myfile"
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to always run subproject myfile,
      but this path is not a directory.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |
