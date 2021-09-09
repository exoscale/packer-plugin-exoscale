module github.com/exoscale/packer-plugin-exoscale

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.2.1
	github.com/aws/aws-sdk-go-v2/config v1.1.2
	github.com/aws/aws-sdk-go-v2/credentials v1.1.2
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.0.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.2.1
	github.com/exoscale/egoscale v0.72.2
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/hashicorp/packer-plugin-sdk v0.2.5
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/rs/xid v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/zclconf/go-cty v1.9.1
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)
