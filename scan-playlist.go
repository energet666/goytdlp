package goytdlp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

type Video struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Id    string `json:"id"`
}

func (ytDlp *YtDlp) ScanPlaylist(url string) ([]Video, error) {
	var videos []Video
	cmd := exec.Command("yt-dlp", "--dump-json", "--flat-playlist", "--config-location", ytDlp.configPath, url)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения stdout: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ошибка запуска: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var video Video
			if err := json.Unmarshal(scanner.Bytes(), &video); err != nil {
				fmt.Println("ошибка разбора JSON:", err)
				continue
			}
			videos = append(videos, video)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("ошибка чтения вывода:", err)
		}
		wg.Done()
	}()

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ошибка выполнения: %v: %s", err, bytes.TrimSpace(stderr.Bytes()))
	}
	return videos, nil
}
