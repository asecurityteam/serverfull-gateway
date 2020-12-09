module github.com/asecurityteam/serverfull-gateway

go 1.12

require (
	github.com/asecurityteam/component-aws v0.1.0
	github.com/asecurityteam/transportd v1.2.4
	github.com/aws/aws-sdk-go v1.36.4
	github.com/getkin/kin-openapi v0.2.1-0.20190729060947-8785b416cb32 // indirect
	github.com/golang/mock v1.4.4
	github.com/stretchr/testify v1.6.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)
