# Morula Developer Documentation

Contributions to Morula are welcome!
Check out the list of [things to do](issues) to get started.


## Set up dev machine

- install [Go](https://golang.org)

- add `./bin` to your PATH

- clone the code base:

  ```
  $ cd $GOPATH/src/github.com/Originate
  $ git clone git@github.com:Originate/morula.git
  ```

- install the dependencies:

  ```
  $ go get github.com/Masterminds/glide github.com/golang/lint/golint
  $ glide install
  ```


## Development
- run the application: `go run main.go`
- run all tests: `bin/spec`
- run only the linter: `bin/lint`
- compile a binary for the local machine: `go install`


## Updating

To update dependencies:

```
$ glide up
```


## Releasing

To publish a new version:

```
$ publish <version>
```

If you want to do it manually:

```
$ git tag -a <version> -m <version>
$ git push --tags
```

When Travis-CI works on a Git tag,
it calls `bin/build_on_travis` and then
deploys the binaries into a GitHub release
via the `deploy` section in [travis.yml](.travis.yml).
