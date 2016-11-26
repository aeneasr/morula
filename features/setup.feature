Feature: setting up the application

  As a casual Morula user
  I want to have an option to have the configuration file set up for me
  So that I can configure Morula without wasting time looking up the configuration file format.

  - run "morula setup" to create a configuration file with the default values


  Scenario: creating a configuration file
    Given a project
    When running "morula setup"
    Then it creates a file "morula.yml" with the content:
      """
      before-all: ""
      after-all: ""
      always: ""
      never: ""
      color: true
      """
