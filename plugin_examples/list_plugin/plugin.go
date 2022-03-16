package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	bhttp "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/http"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/api"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/commands"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin_examples/list_plugin/i18n"
)

type ListPlugin struct{}

func (p *ListPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "ibmcloud-list",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:        "list",
				Description: "List your apps, containers and services in the target space.",
				Usage:       "ibmcloud list",
			},
		},
	}
}

func (p *ListPlugin) Run(context plugin.PluginContext, args []string) {
	p.init(context)

	cf := context.CF()

	ui := terminal.NewStdUI()
	client := &rest.Client{
		HTTPClient:    NewHTTPClient(context),
		DefaultHeader: DefaultHeader(cf),
	}

	var err error
	switch args[0] {
	case "list":
		err = commands.NewList(ui,
			context,
			api.NewCCClient(cf.APIEndpoint(), client),
			api.NewContainerClient(containerEndpoint(cf.APIEndpoint()), client),
		).Run(args[1:])
	}

	if err != nil {
		ui.Failed("%v\n", err)
		os.Exit(1)
	}
}

func (p *ListPlugin) init(context plugin.PluginContext) {
	i18n.T = i18n.Init(context)

	trace.Logger = trace.NewLogger(context.Trace())

	terminal.UserAskedForColors = context.ColorEnabled()
	terminal.InitColorSupport()
}

func NewHTTPClient(context plugin.PluginContext) *http.Client {
	transport := bhttp.NewTraceLoggingTransport(
		&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: context.IsSSLDisabled(),
			},
		})

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(context.HTTPTimeout()) * time.Second,
	}
}

func DefaultHeader(cf plugin.CFContext) http.Header {
	_, err := cf.RefreshUAAToken()
	if err != nil {
		fmt.Println(err.Error())
	}

	h := http.Header{}
	h.Add("Authorization", cf.UAAToken())
	return h
}

func containerEndpoint(cfAPIEndpoint string) string {
	return regexp.MustCompile(`(^https?://)?[^\.]+(\..+)+`).ReplaceAllString(cfAPIEndpoint, "${1}containers-api${2}")
}

func main() {
	plugin.Start(new(ListPlugin))
}
