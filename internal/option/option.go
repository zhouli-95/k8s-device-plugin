package option

const DefaultShareNumber = 4

type Options struct {
	PluginMode     string
	ConfigPath     string
	DeviceStrategy string
	NativeOption   NativeOption
	ShareOption    ShareOption
}

type NativeOption struct {
}

func NewNativeOption() NativeOption {
	return NativeOption{}
}

type ShareOption struct {
	ShareNumber int
}

func NewShareOption() ShareOption {
	return ShareOption{
		ShareNumber: DefaultShareNumber,
	}
}

func New() *Options {
	opt := Options{}
	opt.NativeOption = NewNativeOption()
	opt.ShareOption = NewShareOption()
	return &opt
}
