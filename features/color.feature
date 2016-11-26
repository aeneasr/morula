Feature: configuring the color output

  As a Morula user
  I want to be able to configure whether the tool outputs colors or not
  So that I can optionally remove colors when running it as part of tools.

  - by default colors are printed
  - provide the parameter "color" as "false" to disable colored output


  Background:
    Given a project with the subprojects:
      | NAME | TEMPLATE  |
      | one  | passing_1 |


  Scenario: enabling colors
    When running "morula all --color=true bin/spec"
    Then it prints the output in color


  Scenario: disabling colors
    When running "morula all --color=false bin/spec"
    Then it prints the output without colors
