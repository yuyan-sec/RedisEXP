package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fmt.Printf("文件名: [%s]\n", part.FileName())

		file, _ := os.Create("./" + part.FileName())

		defer file.Close()
		
		io.Copy(file, part)

	}
}

func httpServer(port string) {
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":"+port, nil)
}
