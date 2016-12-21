Feature: testing some subprojects after the others

  As a developer
  I want to be able to run end-to-end tests after the other subprojects
  So that I don't waste time running expensive browser tests that are broken by their dependencies.

  - subprojects listed under the configuration key "after-all" are tested last


  Scenario Outline: valid "after-all" entry in configuration file
    Given a project with the subprojects "one", "two", "e2e", and the configuration file:
      """
      after-all: e2e
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "e2e" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | one |
      | two |
      | e2e |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "after-all" entries in configuration file
    Given a project with the subprojects "one", "two", "e2e-1", "e2e-2", and the configuration file:
      """
      after-all:
        - e2e-1
        - e2e-2
      """
    And I am on the "feature" branch
    And subprojects "one", "two", "e2e-1", and "e2e-2" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | one   |
      | two   |
      | e2e-1 |
      | e2e-2 |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "after-all" entry via command-line parameter
    Given a project with the subprojects "one", "two", "e2e", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "e2e" have changes
    When running "morula <COMMAND> --after-all=e2e bin/spec"
    Then it runs that command in the directories:
      | one |
      | two |
      | e2e |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: without an "after-all" entry
    Given a project with the subprojects "one", "two", "e2e", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "e2e" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | e2e |
      | one |
      | two |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: non-existing "after-all" entry
    Given a project with the subprojects "one", "e2e", and the configuration file:
      """
      after-all: zonk
      """
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to run subproject zonk after all others,
      but such a subproject does not exist.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: "after-all" entry points to a file
    Given a project with the configuration file:
      """
      after-all: myfile
      """
    And the project contains a file "myfile"
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to run subproject myfile after all others,
      but this path is not a directory.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |
