module github.com/rollbar/terraform-provider-rollbar

go 1.14

// Until https://github.com/rs/zerolog/pull/266 or https://github.com/rs/zerolog/pull/267
// is included in the next release
replace github.com/rs/zerolog => github.com/jmcvetta/zerolog v1.20.1-0.20201102133610-4cc56b8f3f5a

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.31.9 // indirect
	github.com/brianvoe/gofakeit/v5 v5.11.1
	github.com/dnaeon/go-vcr v1.1.0
	github.com/fatih/color v1.10.0 // indirect
	github.com/go-resty/resty/v2 v2.3.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.8.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.2.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.0
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/jarcoal/httpmock v1.0.6
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.0
	github.com/oklog/run v1.1.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.6.1
	github.com/zclconf/go-cty v1.7.0 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201207223542-d4d67f95c62d // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201209185603-f92720507ed4 // indirect
	google.golang.org/grpc v1.34.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
