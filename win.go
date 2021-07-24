// +build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// MessageBox win32 API MessageBoxW
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

func messageBoxPlain(title, caption string, flags uint) int {
	const NULL = 0
	return MessageBox(NULL, caption, title, flags)
}

func setWindowSize(x int, y int) {
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("mode con:cols=%d lines=%d", x, y))
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//KernelLib ...
var KernelLib syscall.Handle
var setConsoleTitleWProc uintptr

func freeKernelLib() {
	syscall.FreeLibrary(KernelLib)
}

func loadKernelAndProc() error {
	KernelLib, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return err
	}
	setConsoleTitleWProc, err = syscall.GetProcAddress(KernelLib, "SetConsoleTitleW")
	if err != nil {
		return err
	}
	return nil
}

func setConsoleTitle(title string) (int, error) {
	r, _, err := syscall.Syscall(setConsoleTitleWProc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func setMaxStdio() (int, error) {

	var err error
	var crt syscall.Handle
	crt, err = syscall.LoadLibrary("msvcrt.dll")
	if err != nil {
		crt, err = syscall.LoadLibrary("crtdll.dll")
		if err != nil {
			crt, err = syscall.LoadLibrary("crt.dll")
			if err != nil {
				return 0, err
			}
		}
	}
	defer syscall.FreeLibrary(crt)
	_getmaxstdioProc, err := syscall.GetProcAddress(crt, "_getmaxstdio")
	_setmaxstdioProc, err := syscall.GetProcAddress(crt, "_setmaxstdio")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(_getmaxstdioProc, 0, 0, 0, 0)
	maxstdio := int(ret)

	if maxstdio != 2048 {
		ret, _, _ := syscall.Syscall(_setmaxstdioProc, 0, 2048, 0, 0)
		if int(ret) == 2048 {
			return int(ret), nil
		}
	}

	return maxstdio, nil

}
