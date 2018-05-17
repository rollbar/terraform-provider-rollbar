# rollbar-go

A client written in [go](https://golang.org/) for provisioning rollbar https://rollbar.com.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Usage

```go
import "https://github.com/babbel/rollbar-go/rollbar"
```
Construct a new rollbar client, then use the  services on the client to add or remove users from the Rollbar API. For example:

```go
client, err := rollbar.NewClient("your_api_key")

// list all invites for a team

invites, err := client.ListInvites("team_id")
```

## Prerequisites

You will need to have [go](https://golang.org/) up on running on your system.

## Running the tests

All of the functions are tested and the tests can be run with golang test capability:

```
go test -v
```

## Built With

* [golang](https://golang.org/) - The programming language used.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).

## Authors

* **Nikola Velkovski** - *Initial work* - [parabolic](https://github.com/parabolic)

## License

This project is licensed under the Mozilla License - see the LICENSE.md file for details
