//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package tss

import (
	"fmt"
	"packer-plugin-tss/common"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	common.AuthConfig `mapstructure:",squash"`
	// Secret ID to retrieve from Thycotic Secret Server.
	SecretID int `mapstructure:"secret_id" required:"true"`
	// Names of the Secret Fields to extract from the Secret ID.
	// If empty, all fields will be retrieved.
	// You can also use `exclude_fields` to omit some fields from the output.
	SecretFields []string `mapstructure:"secret_fields"`
	// Names of the Secret Fields to exclude when extracting fields.
	// The fields provided in this variable will not be returned to Packer,
	// even when they are explicitly defined in `secret_fields`.
	ExcludeFields []string `mapstructure:"exclude_fields"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	// Secret ID in TSS.
	ID int `mapstructure:"id"`
	// Key/value combination of the retrieved Secret Fields, where the key is the name of the field.
	// NOTE: The fields defined in `exclude_fields` will not be retrieved.
	Fields map[string]string `mapstructure:"fields"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	if d.config.SecretID == 0 {
		return fmt.Errorf("you must specify the ID of the secret to get its values")
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	client, err := d.config.CreateClient()
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	// TSS SDK only supports retrieving secrets by ID
	secret, err := client.Secret(d.config.SecretID)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	secretFields := make(map[string]string, len(secret.Fields))

	for _, field := range secret.Fields {
		if (len(d.config.SecretFields) == 0 || containsString(d.config.SecretFields, field.Slug)) &&
			(len(d.config.ExcludeFields) == 0 || !containsString(d.config.ExcludeFields, field.Slug)) {
			secretFields[field.Slug] = field.ItemValue
		}
	}

	output := DatasourceOutput{
		ID:     secret.ID,
		Fields: secretFields,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}

func containsString(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}
