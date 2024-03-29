module github.com/asecurityteam/serverfull-gateway

go 1.21

require (
	github.com/asecurityteam/component-aws v0.2.0
	github.com/asecurityteam/transportd v1.11.0
	github.com/aws/aws-sdk-go v1.46.3
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.4
)

require (
	bitbucket.org/atlassian/go-asap v0.0.0-20190921160616-bb88d6193af9 // indirect
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc // indirect
	github.com/asecurityteam/component-connstate v0.2.0 // indirect
	github.com/asecurityteam/component-expvar v0.2.0 // indirect
	github.com/asecurityteam/component-log v0.2.1 // indirect
	github.com/asecurityteam/component-signals v0.2.0 // indirect
	github.com/asecurityteam/component-stat v0.2.0 // indirect
	github.com/asecurityteam/httpstats/v2 v2.4.0 // indirect
	github.com/asecurityteam/logevent v1.6.1 // indirect
	github.com/asecurityteam/runhttp v0.6.1 // indirect
	github.com/asecurityteam/settings v1.0.0 // indirect
	github.com/asecurityteam/transport v1.6.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/getkin/kin-openapi v0.69.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chi/chi v4.0.3+incompatible // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xhandler v0.0.0-20170707052532-1eb70cf1520d // indirect
	github.com/rs/xstats v0.0.0-20170813190920-c67367528e16 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.19.0
	github.com/golang/lint => golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.7.0
)
