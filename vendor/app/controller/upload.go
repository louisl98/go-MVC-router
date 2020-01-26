package controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

// generate random number
func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// use random number as prefix for file name to avoid file overwriting
func TempFile(dir, pattern string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}
	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		suffix, prefix = pattern[:pos], pattern[pos+1:]
	} else {
		suffix = pattern
	}
	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20)
    file, handler, err := r.FormFile("upload")
    if err != nil {
        fmt.Println(err)
        return
	}
    defer file.Close()
    tempFile, err := TempFile("uploads", handler.Filename)
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    tempFile.Write(fileBytes)
}