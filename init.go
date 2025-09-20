package goytdlp

type YtDlp struct {
	configPath string
}

func NewYtDlp(configPath string) *YtDlp {
	return &YtDlp{configPath: configPath}
}
