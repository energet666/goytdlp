package goytdlp

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sync"
)

func (ytDlp *YtDlp) Download(url string) error {
	cmd := exec.Command("yt-dlp", "--config-location", ytDlp.configPath, url)
	var wg sync.WaitGroup

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("ошибка получения stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("ошибка получения stderr: %v", err)
	}

	// Запускаем процесс
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка запуска: %v", err)
	}

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitByCarriageReturn)
		for scanner.Scan() {
			// fmt.Printf("%s\n", scanner.Text())
			re := regexp.MustCompile(`\[download\]\s+(\d{1,3}\.\d{1,2}%)`)

			// Ищем все совпадения
			matches := re.FindStringSubmatch(scanner.Text())
			if len(matches) > 0 {
				fmt.Print("\r[yt-dlp]Процент:", matches[1])
			} else {
				fmt.Printf("[yt-dlp][stdout] %s\n", scanner.Text())
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("ошибка чтения [yt-dlp] stdout: %v\n", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("[yt-dlp][stderr] %s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("ошибка чтения [yt-dlp] stderr: %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ошибка выполнения [yt-dlp]: %v", err)
	}

	return nil
}

func splitByCarriageReturn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// Handle \r\n
		if i+1 < len(data) && data[i+1] == '\n' {
			return i + 2, data[:i], nil
		}
		return i + 1, data[:i], nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[:i], nil
	}

	if !atEOF {
		return 0, nil, nil
	}

	return len(data), data, nil
}
