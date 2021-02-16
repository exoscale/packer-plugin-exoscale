package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"

	"github.com/exoscale/packer-plugin-exoscale/builder/exoscale"
	exoscaleimport "github.com/exoscale/packer-plugin-exoscale/post-processor/exoscale-import"
)

var (
	// Version is the main version number that is being run at the moment.
	Version = "0.0.0"

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"

	// PluginVersion is used by the plugin set to allow Packer to recognize
	// what version this plugin is.
	PluginVersion = version.InitializePluginVersion(Version, VersionPrerelease)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(exoscale.Builder))
	pps.RegisterPostProcessor("import", new(exoscaleimport.PostProcessor))
	pps.SetVersion(PluginVersion)

	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
