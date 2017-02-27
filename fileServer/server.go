package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const VERSION string = "0.10"

func version() {
	fmt.Printf("manager server version:%s\n", VERSION)
}

// 处理/upload 逻辑
func uploadHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	version()
	cfg := NewManagerServerConfig("config.json")

	err := cfg.LoadConfig()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	cfg.DumpConfig()

	http.Handle("/", http.FileServer(http.Dir("./data/")))
	http.HandleFunc("/upload", uploadHandle)
	http.ListenAndServe(":9099", nil)
}
