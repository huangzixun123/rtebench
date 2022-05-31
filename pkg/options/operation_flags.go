package options

import "github.com/spf13/cobra"

type OperationConfigFlagStruct struct {
	Image   string
	Command []string
}

type OperationConfigFlag func(*OperationConfigFlagStruct)

func (f *OperationConfigFlagStruct) WithImage(image string) *OperationConfigFlagStruct {
	f.Image = image
	return f
}

func (f *OperationConfigFlagStruct) WithCommand(command []string) *OperationConfigFlagStruct {
	f.Command = command
	return f
}

func NewOperationConfigFlagStruct(options ...func(*OperationConfigFlagStruct)) *OperationConfigFlagStruct {
	return &OperationConfigFlagStruct{
		Image:   "busybox:latest",
		Command: []string{"sleep", "60"},
	}
}

func (f *OperationConfigFlagStruct) AddFlags(cmd *cobra.Command) {
	if f == nil {
		return
	}
	cmd.Flags().StringVarP(&f.Image, "image", "i", f.Image, "specificy the test images")
	cmd.Flags().StringArrayVar(&f.Command, "command", f.Command, "specificy the test images run command")
}
