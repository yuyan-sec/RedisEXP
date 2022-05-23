package pkg

import (
	"log"
)

const (
	pINFO    = "[+] "
	pSUCCESS = "[*] "
	pErr     = "[-] "
)

func Info(format string) {
	log.Println(pINFO, format)
}

func Err(format error) {
	log.Println(pErr, format)
}

func Success(format interface{}) {
	log.Println(pSUCCESS, format)
}
