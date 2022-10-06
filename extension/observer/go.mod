module github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer

go 1.18

require (
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.23.0
)

require (
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/collector => github.com/dmitryax/opentelemetry-collector v0.49.1-0.20221006222824-5b0277b73df3

replace go.opentelemetry.io/collector/pdata => github.com/dmitryax/opentelemetry-collector/pdata v0.0.0-20221006222824-5b0277b73df3
