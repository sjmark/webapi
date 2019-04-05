package utils

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ***************************************************************************
// 压缩文件
func Compress(files []*os.File, dest string) error {

	// files 文件数组(可以是不同DIR下的文件或者文件夹)
	// dest  压缩文件存放地址

	d, _ := os.Create(dest)

	defer d.Close()

	w := zip.NewWriter(d)

	defer w.Close()

	for _, file := range files {

		err := compress(file, "", w)

		if err != nil {
			return err
		}
	}

	return nil
}
func compress(file *os.File, prefix string, zw *zip.Writer) error {

	info, err := file.Stat()

	if err != nil {
		return err
	}

	if info.IsDir() {

		prefix = prefix + "/" + info.Name()

		fileInfos, err := file.Readdir(-1)

		if err != nil {
			return err
		}

		for _, fi := range fileInfos {

			f, err := os.Open(file.Name() + "/" + fi.Name())

			if err != nil {
				return err
			}

			err = compress(f, prefix, zw)

			if err != nil {
				return err
			}

		}

	} else {

		header, err := zip.FileInfoHeader(info)

		header.Name = prefix + "/" + header.Name

		if err != nil {
			return err
		}

		writer, err := zw.CreateHeader(header)

		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)

		file.Close()

		if err != nil {
			return err
		}

	}

	return nil
}

// ***************************************************************************
func DeCompress(zipffile, dest string) error { // 解压

	reader, err := zip.OpenReader(zipffile)

	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {

		rc, err := file.Open()

		if err != nil {
			return err
		}

		defer rc.Close()

		filename := dest + file.Name

		err = os.MkdirAll(getDir(filename), 0755)

		if err != nil {
			return err
		}

		w, err := os.Create(filename)

		if err != nil {
			return err
		}

		defer w.Close()

		_, err = io.Copy(w, rc)

		if err != nil {
			return err
		}

		w.Close()

		rc.Close()

	}

	return nil
}
func DeCompressToString(zipffile string) (string, error) {

	reader, err := zip.OpenReader(zipffile)

	if err != nil {
		return "", err
	}

	defer reader.Close()

	var buffer bytes.Buffer

	for _, file := range reader.File {

		rc, err := file.Open()

		if err != nil {
			return "", err
		}

		defer rc.Close()

		_, err = io.Copy(&buffer, rc)

		if err != nil {
			return "", err
		}

		rc.Close()

	}

	return buffer.String(), nil
}
func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {

	rs := []rune(str)

	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// ***************************************************************************
func ZipCompress(files map[string]*os.File, dest string) error {

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	for name, file := range files {

		f, err := w.Create(name)

		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(file)

		if err != nil {
			return err
		}

		_, err = f.Write([]byte(bytes))

		if err != nil {
			return err
		}
	}

	err := w.Close()

	if err != nil {
		return err
	}

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	buf.WriteTo(f)

	return nil
}

// ***************************************************************************
//压缩文件Src到Dst
func CompressFile(Dst string, Src string) error {
	newfile, err := os.Create(Dst)
	if err != nil {
		return err
	}
	defer newfile.Close()

	file, err := os.Open(Src)
	if err != nil {
		return err
	}

	zw := gzip.NewWriter(newfile)

	filestat, err := file.Stat()
	if err != nil {
		return err
	}

	zw.Name = filestat.Name()
	zw.ModTime = filestat.ModTime()
	_, err = io.Copy(zw, file)
	if err != nil {
		return err
	}

	zw.Flush()

	if err := zw.Close(); err != nil {
		return err
	}
	return nil
}

