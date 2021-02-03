module github.com/rollbar/terraform-provider-rollbar

go 1.14

// Until https://github.com/rs/zerolog/pull/266 or https://github.com/rs/zerolog/pull/267
// is included in the next release
replace github.com/rs/zerolog => github.com/jmcvetta/zerolog v1.20.1-0.20201102133610-4cc56b8f3f5a

require (
	cloud.google.com/go v0.76.0 // indirect
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.37.2 // indirect
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/dnaeon/go-vcr v1.1.0
	github.com/fatih/color v1.10.0 // indirect
	github.com/gliderlabs/ssh v0.3.2 // indirect
	github.com/go-git/go-git/v5 v5.2.0 // indirect
	github.com/go-resty/resty/v2 v2.4.0
	github.com/go-test/deep v1.0.7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-getter v1.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.2
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jarcoal/httpmock v1.0.8
	github.com/jhump/protoreflect v1.8.1 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/mitchellh/copystructure v1.1.1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/oklog/run v1.1.0 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/zclconf/go-cty v1.7.1 // indirect
	go.opencensus.io v0.22.6 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
