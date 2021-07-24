// +build !windows

package main

func loadKernelAndProc()                                             {}
func freeKernelLib()                                                 {}
func setConsoleTitle(title string) (int, error)                      { return 0, nil }
func setWindowSize(x int, y int)                                     {}
func messageBoxPlain(title, caption string, flags uint) int          { return 0 }
func MessageBox(hwnd uintptr, caption, title string, flags uint) int { return 0 }
func setMaxStdio() (int, error)                                      { return 0, nil }
