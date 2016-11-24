# Morula Developer Documentation


## Set up dev machine

We want to check in code via Git,
so we have to clone the repo using SSH here:

```
$ cd $GOPATH/src/github.com/Originate
$ git clone git@github.com:Originate/morula.git
```

Install the dependencies:

```
$ go get github.com/Masterminds/glide
$ glide install
```


## Development
- compile and run the application: `go run main.go`
- run the tests: `bin/spec`
- compile a binary for the local machine: `go install`


## Updating

To update dependencies, run:

```
$ glide up
```
