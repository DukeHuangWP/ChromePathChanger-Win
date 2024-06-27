package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

const (
	ChromePathDefault = `R:\Canary\Application\chrome.exe`
	ChromePathEnv     = "ChromePath"
)

func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

// MessageBoxPlain of Win32 API.
func MessageBoxPlain(title, caption string) int {
	const (
		NULL  = 0
		MB_OK = 0
	)
	return MessageBox(NULL, caption, title, MB_OK)
}

// 請使用 go build -ldflags -H=windowsgui
//
//go:generate goversioninfo
func main() {

	chromePath := ChromePathDefault
	if _, err := os.Stat(chromePath); err != nil {
		chromePath = os.Getenv(ChromePathEnv)
		log.Println("flag " + chromePath)
	}

	if _, err := os.Stat(chromePath); err != nil {
		dir, err := os.Getwd()
		if err != nil {
			MessageBoxPlain("ChromePath", err.Error())
			return
		} else {
			chromePath = dir + string(os.PathSeparator) + "Application" + string(os.PathSeparator) + `chrome.exe`
			log.Println("pwd " + chromePath)
		}
	}

	if _, err := os.Stat(chromePath); err != nil {
		MessageBoxPlain("ChromePath", fmt.Sprintf("0.'%v' 路徑不存在!\n\n1. 擷取環境變數'%v'='%v' 路徑不存在!\n\n2. 掃描'%v'路徑不存在!", ChromePathDefault, ChromePathEnv, os.Getenv(ChromePathEnv), chromePath))
		return
	}

	userdir := ""
	paths := strings.Split(chromePath, string(os.PathSeparator))
	for index := 0; index < len(paths)-2; index++ {
		if index <= 0 {
			userdir = paths[index]
			continue
		}
		userdir = userdir + string(os.PathSeparator) + paths[index]
	}
	userdir = userdir + string(os.PathSeparator) + "User Data"

	args := os.Args
	args[0] = `--user-data-dir=` + userdir
	cmd := exec.Command(chromePath, args...)
	err := cmd.Run()
	if err != nil {
		MessageBoxPlain("ChromePath", err.Error())
		return
	}
}
