Feature: testing some subprojects before the others

  As a developer
  I want to be able to test shared libraries before subprojects that use them
  So that I don't waste time testing code that are is by its dependencies.

  - subprojects listed under the configuration key `before-all` are tested first


  Scenario Outline: valid "before-all" entry in configuration file
    Given a project with the subprojects "one", "two", "shared", and the configuration file:
      """
      before-all: shared
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "shared" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | shared |
      | one    |
      | two    |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "before-all" entries in configuration file
    Given a project with the subprojects "one", "two", "shared-1", "shared-2", and the configuration file:
      """
      before-all:
        - shared-1
        - shared-2
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "shared" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | shared-1 |
      | shared-2 |
      | one      |
      | two      |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: valid "before-all" entry via command-line parameter
    Given a project with the subprojects "one", "two", "shared", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "shared" have changes
    When running "morula <COMMAND> --before-all=shared bin/spec"
    Then it runs that command in the directories:
      | shared |
      | one    |
      | two    |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: without an "before-all" entry
    Given a project with the subprojects "one", "two", "shared", and the configuration file:
      """
      """
    And I am on the "feature" branch
    And subprojects "one", "two", and "shared" have changes
    When running "morula <COMMAND> bin/spec"
    Then it runs that command in the directories:
      | one |
      | shared |
      | two |

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: non-existing "before-all" entry
    Given a project with the subprojects "one", "shared", and the configuration file:
      """
      before-all: zonk
      """
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to run subproject zonk before all others,
      but such a subproject does not exist.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |


  Scenario Outline: "before-all" entry points to a file
    Given a project with the configuration file:
      """
      before-all: myfile
      """
    And the project contains a file "myfile"
    When trying to run "morula <COMMAND> bin/spec"
    Then it fails with an error code and the message:
      """
      The config file specifies to run subproject myfile before all others,
      but this path is not a directory.
      """

    Examples:
      | COMMAND |
      | all     |
      | changed |
