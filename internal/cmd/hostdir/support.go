package hostdir

import (
	"io/fs"

	"github.com/lucasepe/x/config"
)

/*
func loadConfig(filename string) (config.Config, error) {
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	return config.Parse(fin)
}
*/

func loadConfig(fsys fs.FS, filename string) (config.Config, error) {
	fin, err := fsys.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()

	return config.Parse(fin)
}
