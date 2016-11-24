<img src="documentation/logo.png" width="600" height="111" alt="Morula logo">

[![Build Status](https://travis-ci.org/Originate/morula.svg?branch=master)](https://travis-ci.org/Originate/morula)
[![Go Report Card](https://goreportcard.com/badge/github.com/Originate/morula)](https://goreportcard.com/report/github.com/Originate/morula)

Monorepos are Git repositories that contain multiple code bases,
typically for subprojects of the project in the repo.
Morula runs tasks for all those subprojects within a monorepo.
Optionally only for the ones that contain changes.
This makes running administrative tasks on code bases in monorepos easy, reliable, and fast.


## Installation

Download the appropriate binary for your platform from the
[release page](https://github.com/Originate/morula/releases/latest)
and save it somewhere in your PATH.


## Repo structure

Your monorepository should contain the subprojects in top-level folders,
plus a Morula configuration file.


## Commands

- `morula all <command>`:
  runs the given command in every subproject

- `morula changed <command>`:
  runs the given command in the directories of the subprojects
  that contain changes compared to the master branch


## Configuration file (coming soon)

The config file defines which subprojects should be tested first/last,
and which ones should always/never be tested independent of changes.

__morula.yml__
```yml
main-branch-name: master

before-all:
  - shared

after-all:
  - e2e

always:
  - e2e

never:
  - website
```

Certain directories like `.git` are always ignored.


## Why Monorepos

Large monolithic code bases should be broken up
into more manageable and reusable pieces.
Some of these pieces will be completely independent
from the product they originated from
and become a completely separate project.
Those are the straightforward cases.
Problematic are the pieces that are more like _subprojects_ of the main project
and should remain in the vicinity of it.
These projects are best organized as one big monorepo and not as completely separate projects,
for a number of reasons:

- There is currently not enough tool support
  to work with subprojects that live in their own repositories.
  Cloning, setting up, and keeping dozens of repos up to date
  with ongoing development is a lot of boilerplate activity.
- GitHub doesn't support pull requests across several repos,
  necessitating one pull request per repository to implement many changes.
  This means changes
  that break integration with other subprojects
  cannot be found early in the development process.
- End-to-end testing needs to happen on each change in any subproject
  and combine several repositories,
  which is difficult to implement on CI servers.
- If documentation is extracted into its own subproject,
  it is hard to keep it in sync with ongoing development,
  since documentation updates cannot be part of the pull requests
  for the code changes.

More motivation can be found in the
[monorepo design document of BabelJS](https://github.com/babel/babel/blob/master/doc/design/monorepo.md).

Because of this,
many complex open-source projects
like React, Meteor, Ember, and Babel
have moved towards an architecture that puts
all subprojects into the same repository, i.e. a _monorepository_.
Doing so allows to implement, review, and test changes
in several subprojects together and thereby address all challenges mentioned above.


## Related Work

__[Lerna](https://github.com/lerna/lerna)__
- only works if all subprojects are NPM packages
- enforces an unnecessarily nested directory structure
