module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusremotewritereceiver

go 1.22.0

require (
	go.opentelemetry.io/collector/component v0.110.0
	go.opentelemetry.io/collector/consumer v0.110.0
	go.opentelemetry.io/collector/receiver v0.110.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.110.0 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.110.0 // indirect
	go.opentelemetry.io/collector/internal/globalsignal v0.110.0 // indirect
	go.opentelemetry.io/collector/pdata v1.16.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.110.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.110.0 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240822170219-fc7c04adadcd // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace go.opentelemetry.io/collector/exporter/debugexporter => /Users/danoshin/Projects/opentelemetry-collector/exporter/debugexporter

replace go.opentelemetry.io/collector/exporter => /Users/danoshin/Projects/opentelemetry-collector/exporter

replace go.opentelemetry.io/collector/consumer => /Users/danoshin/Projects/opentelemetry-collector/consumer

replace go.opentelemetry.io/collector/exporter/exporterprofiles => /Users/danoshin/Projects/opentelemetry-collector/exporter/exporterprofiles

replace go.opentelemetry.io/collector/receiver => /Users/danoshin/Projects/opentelemetry-collector/receiver

replace go.opentelemetry.io/collector/receiver/receiverprofiles => /Users/danoshin/Projects/opentelemetry-collector/receiver/receiverprofiles

replace go.opentelemetry.io/collector/receiver/otlpreceiver => /Users/danoshin/Projects/opentelemetry-collector/receiver/otlpreceiver

replace go.opentelemetry.io/collector/service => /Users/danoshin/Projects/opentelemetry-collector/service

replace go.opentelemetry.io/collector/connector/connectorprofiles => /Users/danoshin/Projects/opentelemetry-collector/connector/connectorprofiles

replace go.opentelemetry.io/collector/connector => /Users/danoshin/Projects/opentelemetry-collector/connector

replace go.opentelemetry.io/collector/processor => /Users/danoshin/Projects/opentelemetry-collector/processor

replace go.opentelemetry.io/collector/processor/processorprofiles => /Users/danoshin/Projects/opentelemetry-collector/processor/processorprofiles

replace go.opentelemetry.io/collector/processor/memorylimiterprocessor => /Users/danoshin/Projects/opentelemetry-collector/processor/memorylimiterprocessor

replace go.opentelemetry.io/collector/component => /Users/danoshin/Projects/opentelemetry-collector/component

replace go.opentelemetry.io/collector/component/componentstatus => /Users/danoshin/Projects/opentelemetry-collector/component/componentstatus

replace go.opentelemetry.io/collector/component/componentprofiles => /Users/danoshin/Projects/opentelemetry-collector/component/componentprofiles
