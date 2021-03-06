package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	tokSalt = "*#811"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func GenToken(nickname string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := MD5([]byte(nickname + ts + tokSalt))
	return tokenPrefix + ts[:8]
}

func GenUploadID(nickname string) string {
	return nickname + fmt.Sprintf("%x", time.Now().UnixNano())
}

func Hex2Dec(val string) int64 {
	n, err := strconv.ParseInt(val, 16, 0)
	if err != nil {
		fmt.Println(err)
	}
	return n
}

func MergeChunks(desDirPath string, srcDirPath string, fileName string, n int) error {
	desFilePath := desDirPath + "/" + fileName
	mergedFile, err := os.OpenFile(desFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer mergedFile.Close()

	dir, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	fileInfos, err := dir.Readdir(n)
	if err != nil {
		return err
	}
	defer dir.Close()

	for _, fileInfo := range fileInfos {
		srcFilePath := srcDirPath + "/" + fileInfo.Name()
		chunkFile, err := os.Open(srcFilePath)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		_, err = io.Copy(mergedFile, chunkFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearChunks(srcDirPath string, n int) error {
	dir, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	fileInfos, err := dir.Readdir(n)
	if err != nil {
		return err
	}
	defer dir.Close()

	for _, fileInfo := range fileInfos {
		srcFilePath := srcDirPath + "/" + fileInfo.Name()
		if err := os.Remove(srcFilePath); err != nil {
			return err
		}
	}

	return nil
}
