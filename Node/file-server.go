package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

type Fileserver struct {
	path string
	// files []string

}

func NewFileserver(path string) Fileserver {
	_, err := os.Stat(path)
	if err != nil {
		mydir, _ := os.Getwd()
		log.Println("Path doesn't exist Creating Seed_folder in Current Dir: ", mydir)
		os.Mkdir("Seed_Folder", 0777)
		return Fileserver{path: "./Seed_Folder"}
	}
	return Fileserver{path: path}
}
func (F Fileserver) CheckFile(name string) bool {
	fmt.Println("Checking ", F.path+"/"+name, "in local disk")
	return fs.ValidPath(F.path + "/" + name)
}

func (F Fileserver) GetFile(name string) []byte {
	jk, err := os.ReadFile(fmt.Sprintf("%s/%s", F.path, name))
	if err != nil {
		log.Fatal("unable to read file")
	}
	return jk
}

func (F Fileserver) SaveFile(name string, data []byte) {
	err := os.WriteFile(name, data, 0755)
	if err != nil {
		log.Fatal("Unable to save file")
	}
}
