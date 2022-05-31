package operation

import (
	"strconv"
	"time"

	"github.com/huangzixun123/rtebench/pkg/cri"
	"github.com/huangzixun123/rtebench/pkg/options"
	"github.com/huangzixun123/rtebench/pkg/print"
	"github.com/huangzixun123/rtebench/pkg/types"
	"github.com/huangzixun123/rtebench/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// CreateOptions is the commandline options for 'create' sub command
type OperationOptions struct {
	//PrintFlags  *genericclioptions.PrintFlags
	OK         bool
	Report     types.Report
	PrintObj   func(types.Report)
	ConfigFlag *options.ConfigFlagStruct
	options.IOStreams
}

var (
	operationLong = templates.LongDesc(i18n.T(`
		Run the operation benchmark for the container runtime

		Support OP : Create Run CreateAndRun Destory.`))

	operationExample = templates.Examples(i18n.T(`
		# Run the operation with the default value
		rtectl operation

		# Run the operation with the cri is cri-o and the test round is 10
		rtectl operation --cri=cri-o --round=10
		`))
)

// NewCreateOptions returns an initialized CreateOptions instance
func NewOperationOptions(configflag *options.ConfigFlagStruct, ioStreams options.IOStreams) *OperationOptions {
	configflag.OperationFlags = configflag.OperationFlags.WithImage("busybox:latest").WithCommand([]string{"sleep", "60"})
	return &OperationOptions{
		OK:         false,
		PrintObj:   print.NewOperationPrint,
		ConfigFlag: configflag,
		IOStreams:  ioStreams,
		Report:     make(map[string][]float64),
	}
}

func NewCmd(configflag *options.ConfigFlagStruct, ioStreams options.IOStreams) *cobra.Command {

	o := NewOperationOptions(configflag, ioStreams)

	cmd := &cobra.Command{
		Use:     "retctl operation [flags]",
		Short:   i18n.T("Run the operation benchmark for container runtime"),
		Long:    operationLong,
		Example: operationExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.RunBenchMark(); err != nil {
				return err
			}
			o.OK = true
			if o.OK {
				o.PrintObj(o.Report)
			}
			return nil
		},
	}
	o.ConfigFlag.OperationFlags.AddFlags(cmd)
	return cmd
}

func (o *OperationOptions) RunBenchMark() error {
	client, err := cri.NewClient(util.GetCRIEndpoint(o.ConfigFlag.GeneralFlags.Cri))
	if err != nil {
		return err
	}
	k := 0
	round, _ := strconv.Atoi(o.ConfigFlag.GeneralFlags.Round)
	for k < round {
		err := o.run(client)
		if err != nil {
			return err
		}
		k = k + 1
	}
	return nil
}

func (o *OperationOptions) run(client *cri.Client) error {
	var (
		sandboxID                    = "operation." + util.NewUUID()
		containerID                  = "operation." + util.NewUUID()
		image                        = o.ConfigFlag.OperationFlags.Image
		args                         = o.ConfigFlag.OperationFlags.Command
		beginStartup, endStartup     time.Time // measuring total time
		beginSandbox, endSandbox     time.Time // measuring create sandbox and container
		beginContainer, endContainer time.Time // measuring start containerr
		beginShutdown, endShutdown   time.Time // measuring stop container & sandbox
	)
	// Pull image
	if err := client.PullImage(image, nil); err != nil {
		return err
	}

	// Perform benchmark
	sandbox := client.InitLinuxSandbox(sandboxID)
	beginStartup = time.Now()
	beginSandbox = time.Now()
	pod, err := client.StartSandbox(sandbox, o.ConfigFlag.GeneralFlags.Oci)
	if err != nil {
		return err
	}
	container, err := client.CreateContainer(sandbox, pod, containerID, image, args)
	if err != nil {
		return err
	}
	endSandbox = time.Now()
	beginContainer = time.Now()
	if err := client.StartContainer(container); err != nil {
		return err
	}
	endContainer = time.Now()
	endStartup = time.Now()
	beginShutdown = time.Now()

	// Cleanup container and sandbox
	if err := client.StopAndRemoveContainer(container); err != nil {
		return err
	}
	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return err
	}
	endShutdown = time.Now()
	o.Report["CreateAndRun"] = append(o.Report["CreateAndRun"], endStartup.Sub(beginStartup).Seconds())
	o.Report["Create"] = append(o.Report["Create"], endSandbox.Sub(beginSandbox).Seconds())
	o.Report["Run"] = append(o.Report["Run"], endContainer.Sub(beginContainer).Seconds())
	o.Report["Destroy"] = append(o.Report["Destroy"], endShutdown.Sub(beginShutdown).Seconds())
	return nil
}
