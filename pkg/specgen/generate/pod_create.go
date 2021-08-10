package generate

import (
	"context"

	"github.com/containers/podman/v3/libpod"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func MakePod(p *specgen.PodSpecGenerator, rt *libpod.Runtime) (*libpod.Pod, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if !p.NoInfra && p.InfraContainerSpec != nil {
		p.InfraContainerSpec.IsInfra = true
	}

	options, err := createPodOptions(p, rt, p.InfraContainerSpec)
	if err != nil {
		return nil, err
	}
	pod, err := rt.NewPod(context.Background(), *p, options...)
	if err != nil {
		return nil, err
	}
	if !p.NoInfra && p.InfraContainerSpec != nil {
		p.InfraContainerSpec.Pod = pod.ID()
		if p.InfraContainerSpec.Name == "" {
			p.InfraContainerSpec.Name = pod.ID()[:12] + "-infra"
		}
		_, err = CompleteSpec(context.Background(), rt, p.InfraContainerSpec)
		if err != nil {
			return nil, err
		}
		infraCtr, err := MakeContainer(context.Background(), rt, p.InfraContainerSpec)
		if err != nil {
			return nil, err
		}
		pod, err = rt.AddInfra(context.Background(), pod, infraCtr)
		if err != nil {
			return nil, err
		}
	}
	return pod, nil
}

func createPodOptions(p *specgen.PodSpecGenerator, rt *libpod.Runtime, infraSpec *specgen.SpecGenerator) ([]libpod.PodCreateOption, error) {
	var (
		options []libpod.PodCreateOption
	)
	if !p.NoInfra { //&& infraSpec != nil {
		options = append(options, libpod.WithInfraContainer(p.InfraContainerSpec))
		nsOptions, err := GetNamespaceOptions(p.SharedNamespaces, p.InfraContainerSpec.NetNS.IsHost())
		if err != nil {
			return nil, err
		}
		options = append(options, nsOptions...)
		// Use pod user and infra userns only when --userns is not set to host
		if !p.InfraContainerSpec.UserNS.IsHost() && !p.InfraContainerSpec.UserNS.IsDefault() {
			options = append(options, libpod.WithPodUser())
		}
		// Make our exit command
		storageConfig := rt.StorageConfig()
		runtimeConfig, err := rt.GetConfig()
		if err != nil {
			return nil, err
		}
		exitCommand, err := CreateExitCommandArgs(storageConfig, runtimeConfig, logrus.IsLevelEnabled(logrus.DebugLevel), false, false)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating infra container exit command")
		}
		libpod.WithPodInfraExitCommand(exitCommand, p.InfraContainerSpec)
	}
	if len(p.CgroupParent) > 0 {
		options = append(options, libpod.WithPodCgroupParent(p.CgroupParent))
	}
	if len(p.Labels) > 0 {
		options = append(options, libpod.WithPodLabels(p.Labels))
	}
	if len(p.Name) > 0 {
		options = append(options, libpod.WithPodName(p.Name))
	}
	if p.PodCreateCommand != nil {
		options = append(options, libpod.WithPodCreateCommand(p.PodCreateCommand))
	}

	if len(p.Hostname) > 0 {
		options = append(options, libpod.WithPodHostname(p.Hostname))
	}

	return options, nil
}
