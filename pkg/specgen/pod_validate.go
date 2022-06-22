package specgen

import (
	"github.com/containers/podman/v4/pkg/util"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidPodSpecConfig describes an error given when the podspecgenerator is invalid
	ErrInvalidPodSpecConfig = errors.New("invalid pod spec")
	// containerConfig has the default configurations defined in containers.conf
	containerConfig = util.DefaultContainerConfig()
)

func exclusivePodOptions(opt1, opt2 string) error {
	return errors.Wrapf(ErrInvalidPodSpecConfig, "%s and %s are mutually exclusive pod options", opt1, opt2)
}

// Validate verifies the input is valid
func (p *PodSpecGenerator) Validate() error {
	// PodBasicConfig
	if p.NoInfra {
		if len(p.InfraCommand) > 0 {
			return exclusivePodOptions("NoInfra", "InfraCommand")
		}
		if len(p.InfraImage) > 0 {
			return exclusivePodOptions("NoInfra", "InfraImage")
		}
		if len(p.InfraName) > 0 {
			return exclusivePodOptions("NoInfra", "InfraName")
		}
		if len(p.SharedNamespaces) > 0 {
			return exclusivePodOptions("NoInfra", "SharedNamespaces")
		}
	}

	// PodNetworkConfig
	if err := validateNetNS(&p.InfraContainerSpec.NetNS); err != nil {
		return err
	}
	if p.NoInfra {
		if p.InfraContainerSpec.NetNS.NSMode != Default && p.InfraContainerSpec.NetNS.NSMode != "" {
			return errors.New("NoInfra and network modes cannot be used together")
		}
		// Note that networks might be set when --ip or --mac was set
		// so we need to check that no networks are set without the infra
		if len(p.InfraContainerSpec.Networks) > 0 {
			return errors.New("cannot set networks options without infra container")
		}
		if len(p.InfraContainerSpec.DNSOptions) > 0 {
			return exclusivePodOptions("NoInfra", "DNSOption")
		}
		if len(p.InfraContainerSpec.DNSSearch) > 0 {
			return exclusivePodOptions("NoInfo", "DNSSearch")
		}
		if len(p.InfraContainerSpec.DNSServers) > 0 {
			return exclusivePodOptions("NoInfra", "DNSServer")
		}
		if len(p.InfraContainerSpec.HostAdd) > 0 {
			return exclusivePodOptions("NoInfra", "HostAdd")
		}
		if p.InfraContainerSpec.UseImageResolvConf {
			return exclusivePodOptions("NoInfra", "NoManageResolvConf")
		}
	}
	if p.InfraContainerSpec.NetNS.NSMode != "" && p.InfraContainerSpec.NetNS.NSMode != Bridge && p.InfraContainerSpec.NetNS.NSMode != Slirp && p.InfraContainerSpec.NetNS.NSMode != Default {
		if len(p.InfraContainerSpec.PortMappings) > 0 {
			return errors.New("PortMappings can only be used with Bridge or slirp4netns networking")
		}
	}

	if p.InfraContainerSpec.UseImageResolvConf {
		if len(p.InfraContainerSpec.DNSServers) > 0 {
			return exclusivePodOptions("NoManageResolvConf", "DNSServer")
		}
		if len(p.InfraContainerSpec.DNSSearch) > 0 {
			return exclusivePodOptions("NoManageResolvConf", "DNSSearch")
		}
		if len(p.InfraContainerSpec.DNSOptions) > 0 {
			return exclusivePodOptions("NoManageResolvConf", "DNSOption")
		}
	}
	if p.InfraContainerSpec.UseImageHosts && len(p.InfraContainerSpec.HostAdd) > 0 {
		return exclusivePodOptions("NoManageHosts", "HostAdd")
	}

	return nil
}
