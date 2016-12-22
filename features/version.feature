Feature: displaying the version

  As a Morula user
  I want to have an easy and standardized way to determine the version of Morula I have installed on my machine
  So that I know whether I am up to date, and can work with customer support on investigating issues.

  - run "morula version" to see the version of Morula you are running


  Scenario: displaying the version
    When running "morula version"
    Then it prints:
      """
      Morula v\d+\.\d+(\.\d+)?
      """
