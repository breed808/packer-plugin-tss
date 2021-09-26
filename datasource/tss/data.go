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
	SecretID          int      `mapstructure:"secret_id" required:"true"`
	SecretFields      []string `mapstructure:"secret_fields" required:"true"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	// Secret ID in TSS.
	ID int `mapstructure:"id"`
	// Values of the requested Secret Fields.
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

	secretFields := make(map[string]string, len(d.config.SecretFields))

	for _, field := range d.config.SecretFields {
		var success bool
		secretFields[field], success = secret.Field(field)
		if !success {
			secretFields[field] = ""
		}
	}

	output := DatasourceOutput{
		ID:     secret.ID,
		Fields: secretFields,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
