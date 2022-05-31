package options

type ConfigFlagStruct struct {
	GeneralFlags   *GeneralConfigFlagStruct
	OperationFlags *OperationConfigFlagStruct
	CPUFlags       *CPUConfigFlagStruct
}

func (f *ConfigFlagStruct) WithGeneralOption(gos *GeneralConfigFlagStruct) *ConfigFlagStruct {
	f.GeneralFlags = gos
	return f
}

func (f *ConfigFlagStruct) WithOperationOption(ops *OperationConfigFlagStruct) *ConfigFlagStruct {
	f.OperationFlags = ops
	return f
}

func (f *ConfigFlagStruct) WithCpuOption(cf *CPUConfigFlagStruct) *ConfigFlagStruct {
	f.CPUFlags = cf
	return f
}
func NewConfigFlags() *ConfigFlagStruct {
	defaultGeneralFlag := NewGeneralConfigFlagStruct()
	defaultOperationFlag := NewOperationConfigFlagStruct()
	defaultCpuFlag := NewCpuConfigFlagStruct()
	return &ConfigFlagStruct{
		GeneralFlags:   defaultGeneralFlag,
		OperationFlags: defaultOperationFlag,
		CPUFlags:       defaultCpuFlag,
	}
}
