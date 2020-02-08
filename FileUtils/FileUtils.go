package FiltUtils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

func PathExisted(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}

func FileSha256Hex(filePath string) string {
	var returnSHA1String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String
	}

	hashInBytes := hash.Sum(nil)
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String
}

func encryptedWrite(filePath string, inReader io.Reader, key []byte) error {
	exited := PathExisted(filePath)
	if exited {
		fmt.Println("File is existed: ", filePath)
		return nil
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	writer := &cipher.StreamWriter{
		S: stream,
		W: out,
	}

	if _, err := io.Copy(writer, inReader); err != nil {
		panic(err)
	}

	return nil
}

func DownloadToEncryptedFile(url string, filePath string, key []byte) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return encryptedWrite(filePath, resp.Body, key)
}

func DownloadFile(url string, filePath string) error {
	fmt.Println("Save to ", filePath)
	exited := PathExisted(filePath)
	if exited {
		fmt.Println("File is existed: ", filePath)
		return nil
	}
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DecryptFile(filePath string, outPath string, key []byte) error {
	in, err := os.Open(filePath)
	if err != nil {
		return err
	}

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])
	reader := &cipher.StreamReader{
		S: stream,
		R: in,
	}

	if _, err := io.Copy(out, reader); err != nil {
		panic(err)
	}

	return nil
}
