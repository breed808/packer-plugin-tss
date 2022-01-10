module packer-plugin-tss

replace github.com/thycotic/tss-sdk-go => github.com/breed808/tss-sdk-go v1.0.1-0.20210802091429-0c674ca20814

go 1.16

require (
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/hashicorp/packer-plugin-sdk v0.2.7
	github.com/thycotic/tss-sdk-go v1.0.0
	github.com/zclconf/go-cty v1.10.0
)
