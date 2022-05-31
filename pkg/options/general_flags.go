package options

import "github.com/spf13/cobra"

type GeneralConfigFlagStruct struct {
	Cri   string
	Oci   string
	Round string
}

func (f *GeneralConfigFlagStruct) WithCri(cri string) *GeneralConfigFlagStruct {
	f.Cri = cri
	return f
}

func (f *GeneralConfigFlagStruct) WithOci(oci string) *GeneralConfigFlagStruct {
	f.Oci = oci
	return f
}

func (f *GeneralConfigFlagStruct) WithRound(round string) *GeneralConfigFlagStruct {
	f.Round = round
	return f
}

func NewGeneralConfigFlagStruct(options ...func(*GeneralConfigFlagStruct)) *GeneralConfigFlagStruct {
	return &GeneralConfigFlagStruct{
		Oci:   "runc",
		Cri:   "containerd",
		Round: "1",
	}
}

func (f *GeneralConfigFlagStruct) AddFlags(cmd *cobra.Command) {
	if f == nil {
		return
	}
	cmd.PersistentFlags().StringVarP(&f.Cri, "cri", "c", f.Cri, "specificy the CRI")
	cmd.PersistentFlags().StringVarP(&f.Oci, "oci", "o", f.Oci, "specificy the OCI")
	cmd.PersistentFlags().StringVarP(&f.Round, "round", "r", f.Round, "specificy the test round")
}
