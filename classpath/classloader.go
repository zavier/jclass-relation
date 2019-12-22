package classpath

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 加载读取所有的类
func LoadAll(classPath string) [][]byte {
	pathListSeparator := string(os.PathListSeparator)
	if strings.Contains(classPath, pathListSeparator) {
		total := make([][]byte, 0)
		for _, path := range strings.Split(classPath, pathListSeparator) {
			bytes := loadClass(path)
			total = append(total, bytes...)
		}
		return total
	} else {
		return loadClass(classPath)
	}
}

func loadClass(classPath string) [][]byte {
	if strings.HasSuffix(classPath, "*") {
		return loadWildcardEntry(classPath)
	}
	if strings.HasSuffix(classPath, ".jar") || strings.HasSuffix(classPath, ".JAR") ||
		strings.HasSuffix(classPath, ".zip") || strings.HasSuffix(classPath, ".ZIP") {
		return loadZipEntry(classPath)
	}
	return loadDirEntry(classPath)
}

func loadWildcardEntry(path string) [][]byte {
	baseDir := path[:len(path)-1] // 去掉最后的通配符*

	dataList := make([][]byte, 0)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			zipData := loadZipEntry(path)
			for i := 0; i < len(zipData); i++ {
				dataList = append(dataList, zipData[i])
			}
		}
		return nil
	}
	_ = filepath.Walk(baseDir, walkFn)
	return dataList
}


func loadZipEntry(path string) [][]byte {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	r, err := zip.OpenReader(absPath)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	dataList := make([][]byte, 0)
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".class") {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				panic(err)
			}
			dataList = append(dataList, data)
		}
	}
	return dataList
}

func loadDirEntry(path string) [][]byte {
	absDir, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	dataList := make([][]byte, 0)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			loadDirEntry(filepath.Join(path, info.Name()))
		}
		if strings.HasSuffix(info.Name(), ".class") {
			fileName := filepath.Join(path, info.Name())
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.Println("read file error" + fileName + err.Error())
				return nil
			}
			dataList = append(dataList, data)
		}
		return nil
	}
	_ = filepath.Walk(absDir, walkFn)
	return dataList
}
