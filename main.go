package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main(){
	var key = flag.String("key", "", "Encryption Key")
	var iv = flag.String("iv", "", "Initialization vector")
	var path = flag.String("path", ".", "Input folder")
	var output = flag.String("out", "out", "Output folder")
	flag.Parse()
	paths := getFilesPath(*path)


	for _, p := range(paths){
		folders := strings.Split(p, "/")

		fi, err := os.Stat(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			// do directory stuff
			continue
		}
		fmt.Println(folders)
		errF := folders[:len(folders)-1]
		if errF != nil{
			continue
		}
		folder := strings.Join(folders[:len(folders)-1], "/")
		folder = *output+"/"+folder
		createFolder(folder)

		output_path := *output + p
		f, err := os.Create(output_path)
		if err != nil{
			log.Println("Error creating file: ", output_path)
			log.Println(err)
		}
		defer f.Close()
		result := decrypt(p, *key, *iv)
		f.WriteString(result)
	}
}

func createFolder(path string){
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			//fmt.Println(errDir)
		}

	}
}

func unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func decrypt(path, key, iv string) string {
	data, err := ioutil.ReadFile(path)
	b64data := string(data)
	ciphertext, _ := base64.StdEncoding.DecodeString(b64data)
	var block cipher.Block

	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return ""
	}

	//ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, []byte(iv))
	cbc.CryptBlocks(ciphertext, ciphertext)

	return string(unpad(ciphertext))
}

func getFilesPath(path string) []string {
	var files []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files
}
