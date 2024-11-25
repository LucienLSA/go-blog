package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// 获取当前目录
func Getwd() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fmt.Println("Current working directory:", directory)
	return directory, nil
}

// 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// 获取文件扩展名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 检查文件存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// 检查允许访问
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 检查是否存在 创建目录
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		err := MkDir(src)
		if err != nil {
			return err
		}
	}
	return nil
}

// 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
