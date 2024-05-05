package temp

import (
	"encoding/json"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"os"
	"path/filepath"
)

type info struct {
	Hash string
	Size int64
	Uuid string
}

func (i *info) writeToFile() error {
	f, err := os.Create(filepath.Join(global.Cfg.BaseDir, "/temp/", i.Uuid))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(i)
}

func readFromFile(uuid string) (*info, error) {
	f, err := os.Open(filepath.Join(global.Cfg.BaseDir, "/temp/", uuid))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var i info
	err = json.NewDecoder(f).Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
