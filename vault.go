package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"encoding/hex"
	"encoding/json"
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
	BlockSize int         `json:"b"`
	Groups    []GroupInfo `json:"g"`
}

func (v *Vault) Put(filePath string) (*FileInfo, error) {
	// Write groups of blocks
	groups, blockSize, err := v.writeGroups(filePath)
	if err != nil {
		return &FileInfo{}, err
	}
	// Write file info
	fileName := filepath.Base(filePath)
	fileInfo, err := v.writeFile(fileName, groups, blockSize)
	if err != nil {
		return &FileInfo{}, err
	}
	return fileInfo, nil
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

func (v *Vault) writeGroups(filePath string) (*[]GroupInfo, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return &[]GroupInfo{}, 0, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return &[]GroupInfo{}, 0, err
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
				return &[]GroupInfo{}, 0, err
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
				return &[]GroupInfo{}, 0, err
			}
			defer groupFile.Close()
		}

		// Encrypt block data
		cipher, iv, err := EncryptWithKey(&data, &group.Key)
		if err != nil {
			// TODO : clean up already created blocks !
			return &[]GroupInfo{}, 0, err
		}
		// Update current group SHA
		groupSHA.Write(*cipher)

		// Write and add block to current group
		_, err = groupFile.Write(*cipher)
		if err != nil {
			// TODO : clean up already created blocks !
			return &[]GroupInfo{}, 0, err
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

	return &groups, blockSize, nil
}

func (v *Vault) writeFile(fileName string, groups *[]GroupInfo, blockSize int) (*FileInfo, error) {
	fileId, err := v.genFileId()
	if err != nil {
		return &FileInfo{}, err
	}

	fileInfo := FileInfo{
		Id:        fileId,
		Name:      fileName,
		BlockSize: blockSize,
		Groups:    *groups,
	}
	jsonFileInfo, err := json.Marshal(fileInfo)
	if err != nil {
		return &FileInfo{}, err
	}

	// Encrypt the file info with the master key
	eFileInfo, iv, err := EncryptWithKey(&jsonFileInfo, &v.MasterKey)
	if err != nil {
		return &FileInfo{}, err
	}

	// Write file info to disk
	fp := filepath.Join(v.FilesPath, fileId)
	fh, err := os.Create(fp)
	if err != nil {
		return &FileInfo{}, err
	}

	defer func() {
		if err := fh.Close(); err != nil {
			panic("An error occur")
		}
	}()

	if _, err := fh.Write(*iv); err != nil {
		return &FileInfo{}, err
	}
	if _, err := fh.Write(*eFileInfo); err != nil {
		return &FileInfo{}, err
	}

	return &fileInfo, nil
}

// -----

func (v *Vault) genGroupId() (string, error) {
	return genId(v.DataPath)
}

func (v *Vault) genFileId() (string, error) {
	return genId(v.FilesPath)
}

// Fancy Id generator
// WARNING : Not concurrent safe
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
