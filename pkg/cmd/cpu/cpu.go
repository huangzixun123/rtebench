package cpu

import (
	"fmt"
	"strconv"

	"github.com/huangzixun123/rtebench/pkg/cri"
	"github.com/huangzixun123/rtebench/pkg/options"
	"github.com/huangzixun123/rtebench/pkg/print"
	"github.com/huangzixun123/rtebench/pkg/util"
	"github.com/spf13/cobra"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// CreateOptions is the commandline options for 'create' sub command
type CPUOptions struct {
	//PrintFlags  *genericclioptions.PrintFlags
	OK         bool
	Report     []string
	PrintObj   func([]string)
	ConfigFlag *options.ConfigFlagStruct
	options.IOStreams
}

// NewCreateOptions returns an initialized CreateOptions instance
func NewCPUOptions(configflag *options.ConfigFlagStruct, ioStreams options.IOStreams) *CPUOptions {
	configflag.CPUFlags = configflag.CPUFlags.WithCpuMaxPrime("10").WithNumThreads("1")
	return &CPUOptions{
		OK:         false,
		PrintObj:   print.NewCpuPrint, // TO DO
		ConfigFlag: configflag,
		IOStreams:  ioStreams,
		Report:     []string{},
	}
}

func NewCmd(configflag *options.ConfigFlagStruct, ioStreams options.IOStreams) *cobra.Command {

	o := NewCPUOptions(configflag, ioStreams)

	cmd := &cobra.Command{
		Use:     "cpu --cpu-max-prime=cpu_max_prime",
		Short:   "",
		Long:    "",
		Example: "",
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
	o.ConfigFlag.CPUFlags.AddFlags(cmd)
	return cmd
}

func (o *CPUOptions) RunBenchMark() error {
	client, err := cri.NewClient(util.GetCRIEndpoint(o.ConfigFlag.GeneralFlags.Cri))
	if err != nil {
		return err
	}
	k := 0
	round, _ := strconv.Atoi(o.ConfigFlag.GeneralFlags.Round)
	for k < round {
		report, err := o.run(client)
		if err != nil {
			return err
		}
		o.Report = append(o.Report, report)
		k = k + 1
	}
	return nil
}

func (o *CPUOptions) run(client *cri.Client) (string, error) {
	var (
		sandboxID   = util.NewUUID()
		containerID = util.NewUUID()
		image       = "lnsp/sysbench:latest"
	)

	if err := client.PullImage(image, nil); err != nil {
		return "", err
	}

	sandbox := client.InitLinuxSandbox(sandboxID)

	pod, err := client.StartSandbox(sandbox, o.ConfigFlag.GeneralFlags.Oci)
	if err != nil {
		return "", err
	}
	args := []string{"sysbench", "--test=cpu", fmt.Sprintf("--cpu-max-prime=%s", o.ConfigFlag.CPUFlags.CpuMaxPrime),
		fmt.Sprintf("--num-threads=%s", o.ConfigFlag.CPUFlags.NumThreads), "run"}

	resources := &runtimeapi.LinuxContainerResources{
		CpuPeriod: 100000,
		CpuQuota:  10000,
	}

	container, err := client.CreateContainerWithResources(sandbox, pod, containerID, image, args, resources)
	if err != nil {
		return "", err
	}

	if err := client.StartContainer(container); err != nil {
		return "", err
	}
	logs, err := client.WaitForLogs(container)
	if err != nil {
		return "", err
	}

	if err := client.StopAndRemoveContainer(container); err != nil {
		return "", err
	}

	if err := client.StopAndRemoveSandbox(pod); err != nil {
		return "", err
	}

	return string(logs), nil
}
