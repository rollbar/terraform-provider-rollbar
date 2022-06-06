module github.com/rollbar/terraform-provider-rollbar

go 1.16

// Until https://github.com/rs/zerolog/pull/266 or https://github.com/rs/zerolog/pull/267
// is included in the next release
replace github.com/rs/zerolog => github.com/jmcvetta/zerolog v1.20.1-0.20201102133610-4cc56b8f3f5a

require (
	cloud.google.com/go v0.78.0 // indirect
	cloud.google.com/go/storage v1.14.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.37.25 // indirect
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/dnaeon/go-vcr v1.1.0
	github.com/fatih/color v1.10.0 // indirect
	github.com/go-resty/resty/v2 v2.5.0
	github.com/go-test/deep v1.0.7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/hcl/v2 v2.9.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/jarcoal/httpmock v1.1.0
	github.com/jhump/protoreflect v1.8.2 // indirect
	github.com/klauspost/compress v1.11.12 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nsf/jsondiff v0.0.0-20210303162244-6ea32392771e // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.2
	github.com/ulikunitz/xz v0.5.10 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93 // indirect
	golang.org/x/sys v0.0.0-20210902050250-f475640dd07b // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	google.golang.org/grpc v1.36.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
