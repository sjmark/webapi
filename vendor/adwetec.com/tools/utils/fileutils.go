package utils

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 文件读写
func FileWrite(filename string, m map[string]int64, s string) {

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	defer f.Close()

	if err != nil {
		fmt.Println(fmt.Sprintf("FileWrite ioutil.ReadAll err: %v, s: %v", err, s))
		return
	}

	wb, _ := json.Marshal(m)

	f.WriteString(string(wb))
}

func FileRead(filename string) map[string]int64 {

	f, err := os.Open(filename)

	if err != nil {

		fmt.Println("os.Open err: %v" + err.Error())

		return map[string]int64{}
	}

	b, err := ioutil.ReadAll(f)

	if err != nil {

		fmt.Println("ioutil.ReadAll err: %v" + err.Error())

		return map[string]int64{}
	}

	var data = map[string]int64{}

	json.Unmarshal(b, &data)

	return data
}

// 下载文件
func UploadFileByPath(url, fileName string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, res.Body)

	return err
}

// 解压文件
func Unzip(dirPth string) error {

	// 读文件
	files, err := ioutil.ReadDir(dirPth)

	if err != nil {
		return err
	}
	// 循环下面文件
	for _, file := range files {

		if file.IsDir() { // 忽略目录
			continue
		}
		// zip 打开文件
		File, err := zip.OpenReader(dirPth + file.Name())
		defer File.Close()

		if err != nil {
			return err
		}

		for _, v := range File.File {

			srcFile, err := v.Open()

			defer srcFile.Close()

			if err != nil {
				return err
			}

			newFile, err := os.Create(dirPth + v.Name)

			if err != nil {
				return err
			}

			_, err = io.Copy(newFile, srcFile)
			if err != nil {
				return err
			}
			newFile.Close()
		}
		os.Remove(dirPth + file.Name())
	}
	return nil
}

func FileReadBase64(path string) string {

	f, err := os.Open(path)

	if err != nil {

		fmt.Println("os.Open err: %v" + err.Error())

		return ""
	}

	b, err := ioutil.ReadAll(f)

	if err != nil {

		fmt.Println("ioutil.ReadAll err: %v" + err.Error())

		return ""
	}

	return base64.StdEncoding.EncodeToString(b)
}
func FileReadBytes(path string) []byte {

	f, err := os.Open(path)

	if err != nil {

		fmt.Println("os.Open err: %v" + err.Error())

		return nil
	}

	b, err := ioutil.ReadAll(f)

	if err != nil {

		fmt.Println("ioutil.ReadAll err: %v" + err.Error())

		return nil
	}

	return b
}
func FileWriteBase64(name string, data string) string {

	b, err := base64.StdEncoding.DecodeString(data)

	if err != nil {

		fmt.Println("base64.StdEncoding.DecodeString err: %v" + err.Error())

		return ""
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {

		fmt.Println("os.OpenFile err:%v" + err.Error())

		return ""
	}

	_, err = f.Write(b)

	if err != nil {

		fmt.Println("file.Write err:%v" + err.Error())

		return ""
	}

	return name
}

//将文件分行转换成米的md5

func TransformMD5(SrcPath string, Md5Path string) (bool, error) {

	fi, err := os.OpenFile(SrcPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return false, errors.New("Error: %s\n" + err.Error())
	}
	fileMd5, err := os.OpenFile(Md5Path, os.O_RDWR, os.ModePerm)

	if err != nil {
		return false, errors.New("md5文件创建出错" + err.Error())
	}

	defer fi.Close()
	defer fileMd5.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {

			return true, nil
		}
		h := md5.New()
		h.Write([]byte(string(a))) // 需要加密的字符串
		cipherStr := h.Sum(nil)
		x := hex.EncodeToString(cipherStr) + "\n"
		_, err := fileMd5.WriteString(x)

		if err != nil {
			return false, errors.New("md5文件创建出错" + err.Error())
		}
	}
	return false, nil
}

//分割文件并返回地址列表
func SplitFiles(file *os.File, size int, path string, name string) ([]string, error) {
	finfo, err := file.Stat()

	if err != nil {
		return nil, errors.New("打开文件失败" + err.Error())
	}

	//默认每次拷贝10m
	bufsize := 1024 * 1024 * 10
	if size < bufsize {
		bufsize = size
	}

	buf := make([]byte, bufsize)

	num := (int(finfo.Size()) + size - 1) / size
	paths := make([]string, 0)
	for i := 0; i < num; i++ {

		copylen := 0

		newfilename := path + name + "_" + "(" + strconv.Itoa(i) + ")" + ".txt"

		paths = append(paths, newfilename)
		newfile, err := os.Create(newfilename)

		if err != nil {
			return nil, errors.New("分割创建file失败--file:" + newfilename + err.Error())
		}

		for copylen < size {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				break
			}

			if n <= 0 {
				break
			}
			//写文件
			w_buf := buf[:n]
			newfile.Write(w_buf)
			copylen += n
		}
	}

	return paths, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

//逐行读取数据
func ReadLine(fileName string) ([]string, error) {
	fi, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	strs := make([]string, 0)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strs = append(strs, string(a))
	}

	return strs, nil
}

func WriteFileBase64(name string, data []byte) string {

	b := base64.StdEncoding.EncodeToString(data)

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {

		fmt.Println("os.OpenFile err:%v" + err.Error())

		return ""
	}

	_, err = f.Write([]byte(b))

	if err != nil {

		fmt.Println("file.Write err:%v" + err.Error())

		return ""
	}

	return name
}

func WriteFileString(name string, data string) string {

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {

		fmt.Println("os.OpenFile err:%v" + err.Error())

		return ""
	}

	_, err = f.Write([]byte(data))

	if err != nil {

		fmt.Println("file.Write err:%v" + err.Error())

		return ""
	}

	return name
}

//分割文件并返回地址列表,可自选后缀简单的文件txt，text，别的不行
func SplitFilesWithSuffix(file *os.File, size int, path string, name string, suffix string) ([]string, error) {
	finfo, err := file.Stat()

	if err != nil {
		return nil, errors.New("打开文件失败" + err.Error())
	}

	//默认每次拷贝10m
	bufsize := 1024 * 1024 * 10
	if size < bufsize {
		bufsize = size
	}

	buf := make([]byte, bufsize)

	num := (int(finfo.Size()) + size - 1) / size
	paths := make([]string, 0)
	for i := 0; i < num; i++ {

		copylen := 0

		newfilename := path + name + "_" + "(" + strconv.Itoa(i) + ")" + suffix

		paths = append(paths, newfilename)
		newfile, err := os.Create(newfilename)

		if err != nil {
			return nil, errors.New("分割创建file失败--file:" + newfilename + err.Error())
		}

		for copylen < size {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				break
			}

			if n <= 0 {
				break
			}
			//写文件
			w_buf := buf[:n]
			newfile.Write(w_buf)
			copylen += n
		}
	}

	return paths, nil
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
