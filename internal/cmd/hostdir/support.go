package hostdir

import (
	"bufio"
	"os"
	"strings"

	"github.com/lucasepe/x/config"
)

func loadDict(filename string) (map[string]string, error) {
	dict := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return dict, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		dict[key] = val
	}

	return dict, scanner.Err()
}

func loadConfig(filename string) (config.Config, error) {
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	return config.Parse(fin)
}
