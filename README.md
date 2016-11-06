# Morula

> A tool to manage monorepos,
> i.e. repositories containing the source code of several code bases
> that are used together to create one large product.


## Motivation

Large monolithic code bases should be broken up
into more manageable and reusable pieces.
Some of those pieces will be completely independent
from the product they originated from.
They usually become a completely independent project living in its own repository.
Other pieces are just split off the main code base
to keep things modular and manageable.
They are more like subprojects of the main project,
and should remain a part of the product.

There is currently not enough tool support
to work with subprojects that live in their own repositories.
Development often requires changes to several subprojects at the same time.
GitHub doesn't support pull requests across several repos,
necessitating several pull requests.
End-to-end tests need to run each time any of the subprojects are changed
and require combining several repositories.
This is problematic in several ways on CI servers.

If documentation is extracted into its own subprojects,
it is hard to keep it in sync with ongoing development,
since documentation updates cannot be part of the pull request for the changes.

More motivation in the
[monorepo design document of BabelJS](https://github.com/babel/babel/blob/master/doc/design/monorepo.md)


## Approach: monorepositories

- all subprojects move into top-level directories of one repository
- each subproject can still be worked on by itself within this directory
- changes in several subprojects can be reviewed in one pull request
- one can run all tests for all projects plus end-to-end tests,
  or any combination of them on the CI server

React, Meteor, Ember, and Babel take this approach


## Solution: Morula

- each subproject lives in a top-level folder
- each subproject follows [o-tools](https://github.com/Originate/o-tools-node) conventions
  - `bin/setup` makes the subproject runnable (installs dependencies etc)
  - `bin/spec` runs all tests
- running tests for a branch only tests the changed subprojects
  and ignores ones that have no changes:
  - determine which files are changed compared to the master branch
  - determine which top-level directories contain those changed files
  - run only the tests for the corresponding sub-projects
  - always run the full end-to-end tests at the end
  - possibly run tests for different subprojects in parallel
- the repository contains a config file that:
  - defines the order in which the subprojects are tested
  - defines subprojects that should always/never be tested

Example:
- always test the `shared` repo first
- always test the `e2e` subproject last, if all other subprojects work
- test all other repos in between

```yml
before-all:
  - shared
after-all:
  - e2e
```


## Related Work

__[Lerna](https://github.com/lerna/lerna)__
- only works if all subprojects are NPM packages
- enforces a pretty dumb directory structure

