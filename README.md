[![Build Status](https://travis-ci.org/babbel/rollbar-go.svg?branch=master)](https://travis-ci.org/babbel/rollbar-go)

# rollbar-go

A client written in [go](https://golang.org/) for provisioning rollbar https://rollbar.com. Currently it only supports adding and removing users from teams.

## Usage

```go
go get -u "github.com/babbel/rollbar-go/rollbar"
```

```go
import "github.com/babbel/rollbar-go/rollbar"
```
Construct a new rollbar client, then use the  services on the client to add or remove users from the Rollbar API. For example:

```go
client, err := rollbar.NewClient("your_api_key")


// List all invites for a team.
invites, err := client.ListInvites("team_id")
```

## Prerequisites

You will need to have [go](https://golang.org/), [docker](https://www.docker.com/community-edition#/download) and [docker-compose](https://docs.docker.com/compose/install/) up on running on your system.

## Running the tests

All of the methods are tested and the tests can be run with docker-compose or golang's test command:

```bash
docker-compose up --build
```
```bash
go test -v
```

## Built With

* [golang](https://golang.org/) - The programming language used.
* [docker](https://www.docker.com/community-edition) - Docker CE.
* [docker-compose](https://docs.docker.com/compose/) - Used for building the application.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Nikola Velkovski** - *Initial work* - [parabolic](https://github.com/parabolic)

## License

This project is licensed under the Mozilla License - see the LICENSE.md file for details.
