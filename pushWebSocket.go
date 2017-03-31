package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var image []byte

type CoreFile struct {
}

func (corefile CoreFile) GetFileContent(filename string) []byte {

	var err error
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return content
}

func handleFile(w http.ResponseWriter, r *http.Request) {

	nbytes, nchunks := int64(0), int64(0)

	websocket, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	conn, bwriter, okhijack := websocket.Hijack()

	if okhijack != nil {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
	}

	defer conn.Close()

	file, fopenOk := os.Open("app.txt")

	if fopenOk != nil {
		http.Error(w, "error reading file", http.StatusInternalServerError)
	}

	println("reading file")

	reader := bufio.NewReader(file)

	buff := make([]byte, 1024)

	for {

		println("reading file")
		//n, e := reader.Read(buff[:cap(buff)])
		n, e := reader.Read(buff)
		//buff := buff[:n]

		println("read ok ")

		if n != 0 {
			if e == nil {
				c := string(buff)
				bwriter.WriteString(c)
				bwriter.Flush()
				println("flusing" + c)
				continue
			}
			if e == io.EOF {
				break
			}
		}

		nchunks++
		nbytes += int64(len(buff))

		if e != nil && e != io.EOF {
			log.Fatal(e)
		}
	}
}

func handlerHtml(w http.ResponseWriter, r *http.Request) {
	websocket, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	f := CoreFile{}
	fcontent := f.GetFileContent("app.txt")

	conn, bwriter, okhijack := websocket.Hijack()

	if okhijack != nil {
		http.Error(w, "opsopsosp", http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	c := string(fcontent)

	bwriter.WriteString(c)
	bwriter.Flush()
}

func main() {
	http.HandleFunc("/", handleFile)
	fmt.Println("start http listening :8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
