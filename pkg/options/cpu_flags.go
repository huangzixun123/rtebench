package options

import "github.com/spf13/cobra"

//TODO
type CPUConfigFlagStruct struct {
	CpuMaxPrime string
	NumThreads  string
}

func (f *CPUConfigFlagStruct) WithCpuMaxPrime(cpumaxprime string) *CPUConfigFlagStruct {
	f.CpuMaxPrime = cpumaxprime
	return f
}

func (f *CPUConfigFlagStruct) WithNumThreads(threads string) *CPUConfigFlagStruct {
	f.NumThreads = threads
	return f
}

func (f *CPUConfigFlagStruct) AddFlags(cmd *cobra.Command) {
	if f == nil {
		return
	}
	cmd.Flags().StringVar(&f.CpuMaxPrime, "cpu-max-prime", f.CpuMaxPrime, "specificy the max prime")
	cmd.Flags().StringVar(&f.NumThreads, "num-threads", f.NumThreads, "specificy the num of thread")
}

func NewCpuConfigFlagStruct() *CPUConfigFlagStruct {
	return &CPUConfigFlagStruct{
		CpuMaxPrime: "10",
		NumThreads:  "1",
	}
}
