package main

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	redBoldColor    = color.New(color.FgRed, color.Bold)
	blueBoldColor   = color.New(color.FgBlue, color.Bold)
	yellowBoldColor = color.New(color.FgHiYellow, color.Bold)
	greenBoldColor  = color.New(color.FgGreen, color.Bold)
	whiteBoldColor  = color.New(color.FgWhite, color.Bold)
	cyanBoldColor   = color.New(color.FgCyan, color.Bold)

	y = yellowBoldColor
	r = redBoldColor
	b = blueBoldColor
	g = greenBoldColor
	w = whiteBoldColor
	c = cyanBoldColor

	stdoutReader = bufio.NewScanner(os.Stdin)
)

func logo() {
	println()
	redBoldColor.Println(`
    ____           __  _____       _     __         
   /  _/___  _____/ /_/ ___/____  (_)___/ /__  _____
   / // __ \/ ___/ __/\__ \/ __ \/ / __  / _ \/ ___/
 _/ // / / (__  ) /_ ___/ / /_/ / / /_/ /  __/ /    
/___/_/ /_/____/\__//____/ .___/_/\__,_/\___/_/     
                        /_/`)
	println()
}

func check(description string, e error, exit bool, pcolor *color.Color, scolor *color.Color) {
	if e != nil {
		errorPrint(description, scolor, pcolor)
		if exit {
			Print("Exiting...", pcolor, scolor, true)
			os.Exit(0)
		}
	}
}

func YesOrNo(input string) bool {
	if strings.ToLower(input) == "y" || strings.ToLower(input) == "ye" || strings.ToLower(input) == "yes" {
		return true
	} else {
		return false
	}
}

func simplifyInput(input, expected string) bool {
	if strings.ToLower(input) == expected {
		return true
	} else {
		return false
	}
}

func clear() {
	if runtime.GOOS == "windows" {
		WinClear()
	} else {
		unixClearTerminal()
	}
}

func unixClearTerminal() {
	print("\033[H\033[2J")
}

func unixLastDeleteLine() {
	print("\033[F")
	print("\033[K")
}

func WinClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func userInput(Enter string, pcolor *color.Color, scolor *color.Color, inputColor *color.Color) (string, error) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(Enter + ": ")
	result, err := input(inputColor)

	if err != nil {
		return "", err
	}

	// if result == "" {
	// 	return result, errors.New("input is empty")
	// }

	return result, err

}

func userInputWithDescription(Enter string, Description string, Default string, pcolor *color.Color, scolor *color.Color, inputColor *color.Color) (string, error) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(Enter + " (")
	scolor.Print(Description)
	pcolor.Print("): ")
	result, err := input(inputColor)

	if err != nil {
		return "", err
	}

	if result == "" {
		return Default, nil
	}

	return result, err

}

func userInputNum(Enter string, Default int, pcolor *color.Color, scolor *color.Color, inputColor *color.Color) (int, error) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("] ")
	pcolor.Print(Enter + ": ")

	return inputInt(Default, inputColor)
}

func userInputNumWithDescription(Enter string, Description string, Default int, pcolor *color.Color, scolor *color.Color, inputColor *color.Color) (int, error) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(Enter + " (")
	scolor.Print(Description)
	pcolor.Print("): ")

	return inputInt(Default, inputColor)
}

func errorPrint(e string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("!")
	scolor.Print("]")
	print(" ")
	pcolor.Print(e)
	print("\n")

}

func inputInt(Default int, color *color.Color) (int, error) {

	for {

		TMP, err := input(color)

		if err != nil {
			return 0, err
		}

		if _, err := strconv.Atoi(TMP); err == nil && TMP != "0" && !strings.Contains(TMP, "-") {
			_int64, _ := strconv.ParseInt(TMP, 0, 64)
			return int(_int64), nil
		}

		if TMP == "" {
			return Default, nil
		}

		redBoldColor.Print("[")
		cyanBoldColor.Print("!")
		redBoldColor.Print("] ")
		cyanBoldColor.Println("Enter a valid number")

	}

}

func printWithDescription(out string, Description string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(out + " (")
	scolor.Print(Description)
	pcolor.Println("): ")

}

//Print is colored customized print function, instead of standard print function.
func Print(out string, pcolor *color.Color, scolor *color.Color, newline bool) {
	end := "\n"
	if !newline {
		end = ""
	}

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(out)
	print(end)

}

func printDelete(out string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(out)
	print("\r")

}

func statusSuccessPrint(out, value string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	greenBoldColor.Print(out + ": ")
	pcolor.Print(value)
	print("\n")

}

func statusPrint(out, value string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	pcolor.Print(out + ": ")
	scolor.Print(value)
	print("\n")

}

func printSuccess(out string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	greenBoldColor.Print(out)
	print("\n")

}

func printSuccessWithDescription(Enter string, Description string, pcolor *color.Color, scolor *color.Color) {

	scolor.Print("[")
	pcolor.Print("+")
	scolor.Print("]")
	print(" ")
	greenBoldColor.Print(Enter + " (")
	scolor.Print(Description)
	greenBoldColor.Print("): ")
	print("\n")

}

func input(col *color.Color) (string, error) {

	var userInput string
	if col != nil {

		col.Set()
		defer color.Unset()

	}

	stdoutReader.Scan()
	if err := stdoutReader.Err(); err != nil {
		return "", err
	}
	userInput = stdoutReader.Text()
	return strings.Replace(userInput, "\n", "", -1), nil

}
