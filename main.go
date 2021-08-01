package main

import (
	"fmt"
	"os"
	"packer-plugin-tss/datasource/tss"
	"packer-plugin-tss/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(tss.Datasource))
	pps.SetVersion(version.PluginVersion)

	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
