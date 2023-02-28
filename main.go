package main

// Import required libraries to help do stuff
import (
	"os"

	//mordor "github.com/basingh/vault-custom-support-plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

// func main is the entry function to plugin
func main() {

	// Vault plugin communicates with Vault using gRPC
	// This codes is fetched from https://developer.hashicorp.com/vault/docs/plugins/plugin-development#serving-a-plugin
	// VaultPluginTLSProvider is function to get SSL certificate for Vault and returning Vault TLS configuration object

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	// tlsProviderFunc is used to setup communication between vault and plugin
	// BackendFactoryFunc is used to provide secret engine backend factory to function
	// In this case we pass factory of mordor backend framework

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})

	// logging package hclog been used to surface errors
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
