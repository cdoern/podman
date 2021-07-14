package common

import (
	"net"
	"reflect"

	"github.com/containers/common/pkg/completion"
	"github.com/containers/podman/v3/cmd/podman/parse"
	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/containers/podman/v3/pkg/specgenutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DefineNetFlags(cmd *cobra.Command) {
	netFlags := cmd.Flags()

	addHostFlagName := "add-host"
	netFlags.StringSlice(
		addHostFlagName, []string{},
		"Add a custom host-to-IP mapping (host:ip) (default [])",
	)
	_ = cmd.RegisterFlagCompletionFunc(addHostFlagName, completion.AutocompleteNone)

	dnsFlagName := "dns"
	netFlags.StringSlice(
		dnsFlagName, containerConfig.DNSServers(),
		"Set custom DNS servers",
	)
	_ = cmd.RegisterFlagCompletionFunc(dnsFlagName, completion.AutocompleteNone)

	dnsOptFlagName := "dns-opt"
	netFlags.StringSlice(
		dnsOptFlagName, containerConfig.DNSOptions(),
		"Set custom DNS options",
	)
	_ = cmd.RegisterFlagCompletionFunc(dnsOptFlagName, completion.AutocompleteNone)

	dnsSearchFlagName := "dns-search"
	netFlags.StringSlice(
		dnsSearchFlagName, containerConfig.DNSSearches(),
		"Set custom DNS search domains",
	)
	_ = cmd.RegisterFlagCompletionFunc(dnsSearchFlagName, completion.AutocompleteNone)

	ipFlagName := "ip"
	netFlags.String(
		ipFlagName, "",
		"Specify a static IPv4 address for the container",
	)
	_ = cmd.RegisterFlagCompletionFunc(ipFlagName, completion.AutocompleteNone)

	macAddressFlagName := "mac-address"
	netFlags.String(
		macAddressFlagName, "",
		"Container MAC address (e.g. 92:d0:c6:0a:29:33)",
	)
	_ = cmd.RegisterFlagCompletionFunc(macAddressFlagName, completion.AutocompleteNone)

	networkFlagName := "network"
	netFlags.String(
		networkFlagName, containerConfig.NetNS(),
		"Connect a container to a network",
	)
	_ = cmd.RegisterFlagCompletionFunc(networkFlagName, AutocompleteNetworkFlag)

	networkAliasFlagName := "network-alias"
	netFlags.StringSlice(
		networkAliasFlagName, []string{},
		"Add network-scoped alias for the container",
	)
	_ = cmd.RegisterFlagCompletionFunc(networkAliasFlagName, completion.AutocompleteNone)

	publishFlagName := "publish"
	netFlags.StringSliceP(
		publishFlagName, "p", []string{},
		"Publish a container's port, or a range of ports, to the host (default [])",
	)
	_ = cmd.RegisterFlagCompletionFunc(publishFlagName, completion.AutocompleteNone)

	netFlags.Bool(
		"no-hosts", containerConfig.Containers.NoHosts,
		"Do not create /etc/hosts within the container, instead use the version from the image",
	)
}

func GenerateTempNetFlags(net *entities.NetFlags, infra bool) (pflag.FlagSet, error) {
	if net == nil {
		net = &entities.NetFlags{}
	}
	flags := pflag.FlagSet{}
	fields := reflect.TypeOf(*net)
	values := reflect.ValueOf(*net) //.Elem()
	num := fields.NumField()
	for i := 0; i < num; i++ {
		name := fields.Field(i).Tag.Get("json")
		val := values.Field(i)

		switch val.Kind() {
		case reflect.Slice:
			if !val.IsNil() && !infra {
				return pflag.FlagSet{}, define.ErrInvalidArg
			}
			flags.StringSlice(name, val.Interface().([]string), "")
		case reflect.String:
			if val.String() != "" && !infra {
				return pflag.FlagSet{}, define.ErrInvalidArg
			}
			flags.String(name, val.String(), "")
		case reflect.Bool:
			if val.Bool() && !infra {
				return pflag.FlagSet{}, define.ErrInvalidArg
			}
			flags.Bool(name, val.Bool(), "")
		}
	}
	return flags, nil
}

// NetFlagsToNetOptions parses the network flags for the given cmd.
// The netnsFromConfig bool is used to indicate if the --network flag
// should always be parsed regardless if it was set on the cli.
func NetFlagsToNetOptions(opts *entities.NetOptions, flags pflag.FlagSet, netnsFromConfig bool, createOptions *entities.PodCreateOptions) (*entities.NetOptions, error) {
	var (
		err error
	)
	if createOptions != nil {
		createOptions.NetFlags = &entities.NetFlags{}
	} else {
		createOptions = &entities.PodCreateOptions{
			NetFlags: &entities.NetFlags{},
		}
	}
	if opts == nil {
		opts = &entities.NetOptions{}
	}
	if flags.Changed("add-host") {
		opts.AddHosts, err = flags.GetStringSlice("add-host")
		createOptions.NetFlags.AddHosts = opts.AddHosts
		if err != nil {
			return nil, err
		}
		// Verify the additional hosts are in correct format
		for _, host := range opts.AddHosts {
			if _, err := parse.ValidateExtraHost(host); err != nil {
				return nil, err
			}
		}
	}

	if flags.Changed("dns") {
		servers, err := flags.GetStringSlice("dns")
		createOptions.NetFlags.DNS = servers
		if err != nil {
			return nil, err
		}
		for _, d := range servers {
			if d == "none" {
				opts.UseImageResolvConf = true
				if len(servers) > 1 {
					return nil, errors.Errorf("%s is not allowed to be specified with other DNS ip addresses", d)
				}
				break
			}
			dns := net.ParseIP(d)
			if dns == nil {
				return nil, errors.Errorf("%s is not an ip address", d)
			}
			opts.DNSServers = append(opts.DNSServers, dns)
		}
	}

	if flags.Changed("dns-opt") {
		options, err := flags.GetStringSlice("dns-opt")
		if err != nil {
			return nil, err
		}
		opts.DNSOptions = options
		createOptions.NetFlags.DNSOpt = opts.DNSOptions
	}

	if flags.Changed("dns-search") {
		dnsSearches, err := flags.GetStringSlice("dns-search")
		if err != nil {
			return nil, err
		}
		// Validate domains are good
		for _, dom := range dnsSearches {
			if dom == "." {
				if len(dnsSearches) > 1 {
					return nil, errors.Errorf("cannot pass additional search domains when also specifying '.'")
				}
				continue
			}
			if _, err := parse.ValidateDomain(dom); err != nil {
				return nil, err
			}
		}
		opts.DNSSearch = dnsSearches
		createOptions.NetFlags.DNDSearch = opts.DNSSearch
	}

	if flags.Changed("mac-address") {
		m, err := flags.GetString("mac-address")
		createOptions.NetFlags.MacAddr = m
		if err != nil {
			return nil, err
		}
		if len(m) > 0 {
			mac, err := net.ParseMAC(m)
			if err != nil {
				return nil, err
			}
			opts.StaticMAC = &mac
		}
	}

	if flags.Changed("publish") {
		inputPorts, err := flags.GetStringSlice("publish")
		createOptions.NetFlags.Publish = inputPorts
		if err != nil {
			return nil, err
		}
		if len(inputPorts) > 0 {
			opts.PublishPorts, err = specgenutil.CreatePortBindings(inputPorts)
			if err != nil {
				return nil, err
			}
		}
	}

	if flags.Changed("ip") {
		ip, err := flags.GetString("ip")
		createOptions.NetFlags.IP = ip
		if err != nil {
			return nil, err
		}
		if ip != "" {
			staticIP := net.ParseIP(ip)
			if staticIP == nil {
				return nil, errors.Errorf("%s is not an ip address", ip)
			}
			if staticIP.To4() == nil {
				return nil, errors.Wrapf(define.ErrInvalidArg, "%s is not an IPv4 address", ip)
			}
			opts.StaticIP = &staticIP
		}
	}

	opts.NoHosts, err = flags.GetBool("no-hosts")
	if err != nil {
		logrus.Warnf("no-hosts undefined")
		createOptions.NetFlags.NoHosts = false
		opts.NoHosts = false
		err = nil
	} else {
		createOptions.NetFlags.NoHosts = opts.NoHosts
	}

	// parse the --network value only when the flag is set or we need to use
	// the netns config value, e.g. when --pod is not used
	if netnsFromConfig || flags.Changed("network") {
		network, err := flags.GetString("network")
		createOptions.NetFlags.Network = network
		if err != nil {
			return nil, err
		}

		ns, cniNets, options, err := specgen.ParseNetworkString(network)
		if err != nil {
			return nil, err
		}

		if len(options) > 0 {
			opts.NetworkOptions = options
		}
		opts.Network = ns
		opts.CNINetworks = cniNets
	}

	if flags.Changed("network-alias") {
		aliases, err := flags.GetStringSlice("network-alias")
		createOptions.NetFlags.NetworkAlias = aliases
		if err != nil {
			return nil, err
		}
		if len(aliases) > 0 {
			opts.Aliases = aliases
		}
	}
	return opts, err
}
