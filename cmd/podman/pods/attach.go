package pods

import (
	"context"
	"fmt"

	"github.com/containers/podman/v3/cmd/podman/common"
	"github.com/containers/podman/v3/cmd/podman/registry"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	podSplitDescription = `After adding an existing container to a pod, the container inherits all cgroup information and acts as a normal member of the pod.`

	splitCommand = &cobra.Command{
		Use:               "add [options] POD",
		Args:              cobra.ExactArgs(1),
		Short:             "Add an existing container to a pod",
		Long:              podSplitDescription,
		RunE:              attach,
		ValidArgsFunction: common.AutocompletePods,
	}
)

var (
	splitOptions entities.PodAttachOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: splitCommand,
		Parent:  podCmd,
	})
	flags := splitCommand.Flags()
	flags.SetInterspersed(false)

	idFlagName := "id"
	flags.StringVar(&splitOptions.ID, idFlagName, "", "ID of the pod to attach to")
	_ = splitCommand.RegisterFlagCompletionFunc(idFlagName, common.AutocompletePods)

	cidFlagName := "cid"
	flags.StringVar(&splitOptions.ID, cidFlagName, "", "ID of the container to attach to the given pod")
	_ = splitCommand.RegisterFlagCompletionFunc(idFlagName, common.AutocompleteContainers)
}

func attach(cmd *cobra.Command, args []string) error {
	var (
		err      error
		response *entities.PodAttachReport
	)

	response, err = registry.ContainerEngine().PodAttach(context.Background(), splitOptions.ID, splitOptions.CID)

	if err != nil {
		return err
	}

	fmt.Println(response.Id)

	return nil
}
