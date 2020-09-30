module github.com/heroku/terraform-provider-librato

go 1.12

require (
	github.com/hashicorp/terraform v0.12.0
	github.com/heroku/go-librato v0.0.0-20171209220328-a5ddfd932e56
	golang.org/x/sys v0.0.0-20190619223125-e40ef342dc56 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
