package main

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/hex"
	"crypto/rand"
)

const (
	ID_SIZE = 16
)

type Vault struct {
	DataPath  string `json:"d"`
	FilesPath string `json:"f"`
	MasterKey []byte `json:"m"`
}

type BlockInfo struct {
	Iv []byte `json:"i"`
}

type GroupInfo struct {
	Id     string      `json:"i"`
	Key    []byte      `json:"k"`
	Hash   string      `json:"h"`
	Blocks []BlockInfo `json:"b"`
}

type FileInfo struct {
	Id        string      `json:"i"`
	Name      string      `json:"n"`
	BlockSize uint        `json:"b"`
	Groups    []GroupInfo `json:"g"`
}

func (v *Vault) Put(filePath string) (*FileInfo, error) {
	fileId, err := v.genFileId()
	if err != nil {
		fmt.Println("Cannot generate file id !")
		os.Exit(1)
	}

	fileName := filepath.Base(filePath)

	return &FileInfo{
		Id: fileId,
		Name: fileName,
	}, nil
}


//
// Helper functions
//

func (v *Vault) genFileId() (string, error) {
	return genId(v.FilesPath)
}

// Fancy Id generator
// WARNING : Not thread safe
func genId(dir string) (string, error) {
	idBytes := make([]byte, ID_SIZE)
	for {
		_, err := rand.Read(idBytes)
		if err != nil {
			return "", err
		}
		id :=  hex.EncodeToString(idBytes)
		fp := filepath.Join(dir, id)
		if _, err = os.Stat(fp); err != nil {
			return id, nil
		}
	}
}
