module github.com/asecurityteam/serverfull-gateway

go 1.12

require (
	github.com/asecurityteam/component-aws v0.1.0
	github.com/asecurityteam/httpstats v0.0.0-20191007213332-05cb203c96fb // indirect
	github.com/asecurityteam/transportd v1.0.0
	github.com/aws/aws-sdk-go v1.25.9
	github.com/getkin/kin-openapi v0.2.1-0.20190729060947-8785b416cb32 // indirect
	github.com/golang/mock v1.4.0
	github.com/justinas/alice v0.0.0-20171023064455-03f45bd4b7da // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191009170851-d66e71096ffb // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)
