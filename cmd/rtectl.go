package cmd

import (
	"os"

	"github.com/huangzixun123/rtebench/pkg/cmd/cpu"
	"github.com/huangzixun123/rtebench/pkg/cmd/operation"
	"github.com/huangzixun123/rtebench/pkg/options"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

type RteOptions struct {
	ConfigFlags *options.ConfigFlagStruct

	options.IOStreams
}

var generalConfigFlag = options.NewGeneralConfigFlagStruct()
var defaultConfigFlags = options.NewConfigFlags().WithGeneralOption(generalConfigFlag)

// NewDefaultRetctlCommand creates the `rtectl` command with default arguments
func NewDefaultRetctlCommand() *cobra.Command {
	return NewDefaultRetctlCommandWithArgs(&RteOptions{
		ConfigFlags: defaultConfigFlags,
		IOStreams:   options.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

// NewDefaultKubectlCommandWithArgs creates the `kubectl` command with arguments
func NewDefaultRetctlCommandWithArgs(o *RteOptions) *cobra.Command {
	cmd := NewRtectlCommand(o)
	return cmd
}

// NewRtectlCommand creates the `retctl` command and its nested children.
func NewRtectlCommand(o *RteOptions) *cobra.Command {

	cmds := &cobra.Command{
		Use:   "rtectl",
		Short: "rtectl is a benchmark tool based on Golang",
		Long: templates.LongDesc(`
		rtectl is a benchmark tool based on Golang.
		It is most frequently used for container runtime benchmarks.
		`),
		Run: runHelp,
	}
	o.ConfigFlags.GeneralFlags.AddFlags(cmds)

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands",
			Commands: []*cobra.Command{
				operation.NewCmd(o.ConfigFlags, o.IOStreams),
				cpu.NewCmd(o.ConfigFlags, o.IOStreams),
			},
		},
	}
	groups.Add(cmds)

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
