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
