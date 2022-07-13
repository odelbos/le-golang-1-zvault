package main

type Vault struct {
	DataPath string `json:"d"`
	FilesPath string `json:"f"`
	MasterKey []byte `json:"m"`
}
