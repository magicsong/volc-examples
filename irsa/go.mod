module github.com/magicsong/volc-examples/irsa

go 1.23

require (
	github.com/volcengine/ve-tos-golang-sdk/v2 v2.7.9
	github.com/volcengine/volcengine-go-sdk v1.0.186
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/volcengine/volc-sdk-golang v1.0.23 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/volcengine/volcengine-go-sdk v1.0.186 => github.com/magicsong/volcengine-go-sdk v0.0.5
