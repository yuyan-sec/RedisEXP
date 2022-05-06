package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

var (
	filePath string
	addr     string
)

func init() {
	flag.StringVar(&filePath, "f", "", "Filename")
	flag.StringVar(&addr, "u", "", "URL")
	flag.Parse()
}

func main() {
	if filePath == "" {
		return
	}

	err := doUpload(addr, filePath)
	if err != nil {
		fmt.Printf("upload file [%s] error: %s", filePath, err)
		return
	}
	fmt.Printf("upload file [%s] ok\n", filePath)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createReqBody(filePath string) (string, io.Reader, error) {
	var err error
	pr, pw := io.Pipe()
	bw := multipart.NewWriter(pw) // body writer
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}

	go func() {
		defer f.Close()

		_, fileName := filepath.Split(filePath)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes("file"), escapeQuotes(fileName)))
		h.Set("Content-Type", "image/png")
		fw1, _ := bw.CreatePart(h)
		var buf = make([]byte, 1024)
		cnt, _ := io.CopyBuffer(fw1, f, buf)
		log.Printf("copy %d bytes from file %s in total\n", cnt, fileName)
		bw.Close()
		pw.Close()
	}()
	return bw.FormDataContentType(), pr, nil
}

func doUpload(addr, filePath string) error {
	// create body
	contType, reader, err := createReqBody(filePath)
	if err != nil {
		return err
	}

	log.Printf("createReqBody ok\n")
	url := fmt.Sprintf("http://%s/upload", addr)
	req, err := http.NewRequest("POST", url, reader)

	//add headers
	req.Header.Add("Content-Type", contType)

	client := &http.Client{}
	log.Printf("upload %s...\n", filePath)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request send error:", err)
		return err
	}
	resp.Body.Close()
	log.Printf("upload %s ok\n", filePath)
	return nil
}
