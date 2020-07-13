# rollbar-go

A client written in [go](https://golang.org/) for provisioning rollbar https://rollbar.com. Currently it only supports adding and removing users from teams.

### Usage

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
teamID := 2131551312
invites, err := client.ListInvites(teamID)
```

### Prerequisites

You will need to have [go](https://golang.org/), [docker](https://www.docker.com/community-edition#/download) and [docker-compose](https://docs.docker.com/compose/install/) up and running on your system.

### Running the tests

All of the methods are tested and the tests can be run with docker-compose or golang's test command:

```bash
docker-compose up --build
```

```bash
go test -v
```

Liniting and go vet can be run locally using the makefile.

### Contributing

Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.

### Authors

* **Nikola Velkovski** - *Initial work* - [parabolic](https://github.com/parabolic)

### License

This project is licensed under the Mozilla License - see the LICENSE.md file for details.
