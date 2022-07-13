package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"encoding/hex"
	"crypto/rand"
	"crypto/sha256"
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

	groups, err := v.writeGroups(filePath)
	if err != nil {
		return &FileInfo{}, err
	}


	fileId, err := v.genFileId()
	if err != nil {
		fmt.Println("Cannot generate file id !")
		os.Exit(1)
	}

	fileName := filepath.Base(filePath)

	return &FileInfo{
		Id: fileId,
		Name: fileName,
		Groups: *groups,
	}, nil
}


//
// Helper functions
//

func (v *Vault) newGroup() *GroupInfo {
	id, err := v.genGroupId()
	if err != nil {
		// TODO : Error
	}
	return &GroupInfo{
		Id:  id,
		Key: GenCryptoRand(AES_KEY_SIZE),
	}
}

func (v *Vault) writeGroups(filePath string) (*[]GroupInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return &[]GroupInfo{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return &[]GroupInfo{}, err
	}
	fileSize := stat.Size()

	nbBlocks, blockSize := NbBlocksPerGroup(fileSize)

	fmt.Printf("file size : %v - %T\n", fileSize, fileSize) // TODO : debug
	fmt.Printf("Nb blocks per group : %v\n", nbBlocks)      // TODO : debug

	// -----

	var groups []GroupInfo
	var group *GroupInfo
	var groupFile *os.File
	groupSHA := sha256.New()
	count := 0
	buffer := make([]byte, blockSize)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				// TODO : Clean up the already created blocks
				return &[]GroupInfo{}, err
			}

			groupFile.Close()
			if len(group.Blocks) != 0 {
				group.Hash = hex.EncodeToString(groupSHA.Sum(nil))
				fmt.Printf("Group: %v, sha: %v\n", group.Id, group.Hash)
				groups = append(groups, *group)
			}
			break
		}
		data := buffer[:bytesRead]

		if count == 0 {
			group = v.newGroup()
			groupSHA.Reset()
			groupSHA.Write([]byte(group.Id))

			fp := filepath.Join(v.DataPath, group.Id)
			groupFile, err = os.Create(fp)
			if err != nil {
				// TODO : Clean up the already created blocks
				return &[]GroupInfo{}, err
			}
			defer groupFile.Close()
		}

		// Encrypt block data
		cipher, iv, err := EncryptWithKey(&data, &group.Key)
		if err != nil {
			// TODO : clean up already created blocks !
			return &[]GroupInfo{}, err
		}
		// Update current group SHA
		groupSHA.Write(*cipher)

		// Write and add block to current group
		_, err = groupFile.Write(*cipher)
		if err != nil {
			// TODO : clean up already created blocks !
			return &[]GroupInfo{}, err
		}
		block := BlockInfo{
			Iv: *iv,
		}
		group.Blocks = append(group.Blocks, block)

		if count += 1; count == nbBlocks {
			groupFile.Close()
			group.Hash = hex.EncodeToString(groupSHA.Sum(nil))

			fmt.Printf("Group: %v, sha: %v\n", group.Id, group.Hash)

			groups = append(groups, *group)
			count = 0
		}
	}

	return &groups, nil
}

// -----

func (v *Vault) genGroupId() (string, error) {
	return genId(v.DataPath)
}

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
