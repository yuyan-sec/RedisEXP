package logger

import (
	"fmt"
	"os"
)

const (
	pINFO    = "[+] "
	pSUCCESS = "[*] "
	pErr     = "[-] "
)

func Info(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, pINFO+format+"\n", args...)
}

func Err(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, pErr+format+"\n", args...)
}

func Success(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, pSUCCESS+format+"\n", args...)
}
