# Morula

> efficient testing for monorepos


Monorepositories, or monorepos, are repositories that contain multiple code bases that belong together.
They allow efficient work on a complex project
that was broken up into several subprojects
to make things more manageable.

Examples of popular monorepos are:

* https://github.com/docker/docker
* https://github.com/nodejs/node
* https://github.com/kubernetes/kubernetes
* https://github.com/angular/angular.js
* https://github.com/torvalds/linux
* https://github.com/twbs/bootstrap
* https://github.com/facebook/react

The tool Morula provides facilities
to run the tests for all subprojects
affected by a particular change in a monorepo.
This makes testing monorepos via a CI server easy, reliable and fast.


## Repo structure

A monorepository using Morula, also known as a Morula monorepo, contains the subprojects located in top-level folders,
plus a Morula configuration file.
Each subproject can be worked on by itself within its directory,
and is versioned and released independent from the other subprojects.
Each top-level folder subproject follow
[o-tools](https://github.com/Originate/o-tools-node) development structure conventions,
meaning they contain a standardized set of scripts:
- `bin/setup`: makes the subproject runnable, for example by installing dependencies
- `bin/spec`: runs all tests for this subproject


## Configuration file

The config file is named `morula.yml` and
defines which subprojects should be tested first/last,
and which ones should always/never be tested independent of changes.

__morula.yml__
```yml
before-all:
  - shared
after-all:
  - e2e
always:
  - e2e
never:
  - website
```


## Commands

- `morula setup`:
  runs the `bin/setup` scripts for each subproject

- `morula test`:
  determines which folders contain changes
  and runs the tests for only the respective subprojects.


## Why Monorepos

Large monolithic code bases should be broken up
into more manageable and reusable pieces.
Some of these pieces will be completely independent
from the product they originated from
and become a completely separate project.
Those are the straightforward cases.
Problematic are the pieces that are more like _subprojects_ of the main project
and should remain in the vicinity of the main project.

Reasons in problem domain:

* There is currently not enough tool support
to work with subprojects that live in their own repositories.
* Cloning, setting up, and keeping dozens of repos up to date
with ongoing development is a lot of boilerplate activity.
* Configuring subprojects to run against locally checked out
vs published dependencies is another pain.
It requires switching several subprojects to branches under current development.
* GitHub doesn't support pull requests across several repos,
necessitating one pull request per repository to implement many changes.
This means changes
that break integration with other subprojects
can not be found early in the development process.
* End-to-end testing needs to happen on each change in any subproject
and combine several repositories,
which is difficult to implement on CI servers.
* If documentation is extracted into its own subproject,
it is hard to keep it in sync with ongoing development,
since documentation updates cannot be part of the pull requests
for the code changes.
More motivation in the
[monorepo design document of BabelJS](https://github.com/babel/babel/blob/master/doc/design/monorepo.md)

Because of all these reasons,
many complex open-source projects
like React, Meteor, Ember, and Babel
have moved towards an architecture that puts
all subprojects into the same repository, i.e. a _monorepository_.
Doing so allows to implement, review, and test changes
in several subprojects together and thereby address all challenges mentioned above.


## Related Work

__[Lerna](https://github.com/lerna/lerna)__

Issues:
- only works if all subprojects are NPM packages
- enforces an unnecessarily nested directory structure

