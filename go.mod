module github.com/rollbar/terraform-provider-rollbar

go 1.16

// Until https://github.com/rs/zerolog/pull/266 or https://github.com/rs/zerolog/pull/267
// is included in the next release
replace github.com/rs/zerolog => github.com/jmcvetta/zerolog v1.20.1-0.20201102133610-4cc56b8f3f5a

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/dnaeon/go-vcr v1.1.0
	github.com/go-resty/resty/v2 v2.5.0
	github.com/go-test/deep v1.0.7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.29.0
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/jarcoal/httpmock v1.1.0
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/oklog/run v1.1.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.8.3
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
