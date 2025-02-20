// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/structarg"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/openstack"
	_ "yunion.io/x/cloudmux/pkg/multicloud/openstack/shell"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

type BaseOptions struct {
	Debug         bool   `help:"debug mode"`
	AuthURL       string `help:"Auth URL" default:"$OPENSTACK_AUTH_URL" metavar:"OPENSTACK_AUTH_URL"`
	Username      string `help:"Username" default:"$OPENSTACK_USERNAME" metavar:"OPENSTACK_USERNAME"`
	Password      string `help:"Password" default:"$OPENSTACK_PASSWORD" metavar:"OPENSTACK_PASSWORD"`
	Project       string `help:"Project" default:"$OPENSTACK_PROJECT" metavar:"OPENSTACK_PROJECT"`
	EndpointType  string `help:"Project" default:"$OPENSTACK_ENDPOINT_TYPE|internal" metavar:"OPENSTACK_ENDPOINT_TYPE"`
	DomainName    string `help:"Domain of user" default:"$OPENSTACK_DOMAIN_NAME|Default" metavar:"OPENSTACK_DOMAIN_NAME"`
	ProjectDomain string `help:"Domain of project" default:"$OPENSTACK_PROJECT_DOMAIN|Default" metavar:"OPENSTACK_PROJECT_DOMAIN"`
	RegionID      string `help:"RegionId" default:"$OPENSTACK_REGION_ID" metavar:"OPENSTACK_REGION_ID"`
	SUBCOMMAND    string `help:"openstackcli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, e := structarg.NewArgumentParserWithHelp(&BaseOptions{},
		"openstackcli",
		"Command-line interface to openstack API.",
		`See "openstackcli COMMAND --help" for help on a specific command.`)

	if e != nil {
		return nil, e
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, fmt.Errorf("No subcommand argument.")
	}
	for _, v := range shellutils.CommandTable {
		_, e := subcmd.AddSubParserWithHelp(v.Options, v.Command, v.Desc, v.Callback)
		if e != nil {
			return nil, e
		}
	}
	return parse, nil
}

func showErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "%s", e)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func newClient(options *BaseOptions) (*openstack.SRegion, error) {
	if len(options.AuthURL) == 0 {
		return nil, fmt.Errorf("Missing AuthURL")
	}

	if len(options.Username) == 0 {
		return nil, fmt.Errorf("Missing Username")
	}

	if len(options.Password) == 0 {
		return nil, fmt.Errorf("Missing Password")
	}

	cfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := cfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	cli, err := openstack.NewOpenStackClient(
		openstack.NewOpenstackClientConfig(
			options.AuthURL,
			options.Username,
			options.Password,
			options.Project,
			options.ProjectDomain,
		).
			EndpointType(options.EndpointType).
			DomainName(options.DomainName).
			Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc: proxyFunc,
				},
			),
	)
	if err != nil {
		return nil, err
	}
	region := cli.GetRegion(options.RegionID)
	if region == nil {
		return nil, fmt.Errorf("No such region %s", options.RegionID)
	}
	return region, nil
}

func main() {
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}
	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*BaseOptions)

	if parser.IsHelpSet() {
		fmt.Print(parser.HelpString())
		return
	}
	subcmd := parser.GetSubcommand()
	subparser := subcmd.GetSubParser()
	if e != nil || subparser == nil {
		if subparser != nil {
			fmt.Print(subparser.Usage())
		} else {
			fmt.Print(parser.Usage())
		}
		showErrorAndExit(e)
	}
	suboptions := subparser.Options()
	if subparser.IsHelpSet() {
		fmt.Print(subparser.HelpString())
		return
	}
	var region *openstack.SRegion
	if len(options.RegionID) == 0 {
		options.RegionID = openstack.OPENSTACK_DEFAULT_REGION
	}
	region, e = newClient(options)
	if e != nil {
		showErrorAndExit(e)
	}
	e = subcmd.Invoke(region, suboptions)
	if e != nil {
		showErrorAndExit(e)
	}
}
