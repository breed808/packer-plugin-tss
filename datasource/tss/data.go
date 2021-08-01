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
	SecretID          int `mapstructure:"secret_id" required:"true"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	ID int `mapstructure:"id"`

	// Though TSS stores other fields, retrieve only credential details (username & password) for now.
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
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
	output := DatasourceOutput{}

	emptyOutput := hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec())

	client, err := d.config.CreateClient()
	if err != nil {
		return emptyOutput, err
	}

	// TSS SDK only supports retrieving secrets by ID
	secret, err := client.Secret(d.config.SecretID)
	if err != nil {
		return emptyOutput, err
	}

	output.ID = secret.ID

	var success bool
	output.Username, success = secret.Field("username")
	if !success {
		output.Username = ""
	}

	output.Password, success = secret.Field("password")
	if !success {
		output.Password = ""
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
