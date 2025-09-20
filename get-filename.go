package goytdlp

import (
	"fmt"
	"os/exec"
	"strings"
)

func (ytDlp *YtDlp) GetFilename(url string) (filename string, err error) {
	cmd := exec.Command("yt-dlp", "--print", "filename", "--config-location", ytDlp.configPath, url)
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ошибка получения stdout: %v", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
