package msg

import "fmt"

const (
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"
	Bright  = "\x1b[1m"
	Reset   = "\x1b[0m"
)

var ColorMap map[string]string
var ColorMapKeys []string

// Err is a wrapper around the Msg() function, which is a wrapper around fmt.Printf()
// The function prints the text, but using the terminal color code for Red.
func Err(f string, a ...interface{}) {
	initMap()
	Msg(Red, f, a...)
}

// Info is a wrapper around the Msg() function, which is a wrapper around fmt.Printf()
// The function prints the text, but using the terminal color code for Green.
func Info(f string, a ...interface{}) {
	initMap()
	Msg(Green, f, a...)

}

// Msg is a wrapper around fmt.Printf() that takes a color code sequence as the first parameter.
// It also prints the Reset string at the end of the function.
func Msg(color string, f string, a ...interface{}) {
	initMap()

	// Make sure we always reset
	defer fmt.Print(Reset)

	// Dark Mode FTW
	if color == Black {
		fmt.Print(Bright)
	}
	fmt.Print(color)
	fmt.Printf(f, a...)
}

func initMap() {
	if len(ColorMap) == 0 {
		ColorMap = make(map[string]string, 10)
		ColorMap["Black"] = Black
		ColorMap["Red"] = Red
		ColorMap["Green"] = Green
		ColorMap["Yellow"] = Yellow
		ColorMap["Blue"] = Blue
		ColorMap["Magenta"] = Magenta
		ColorMap["Cyan"] = Cyan
		ColorMap["White"] = White
		ColorMap["Bright"] = Bright
		ColorMap["Reset"] = Reset
	}
	if len(ColorMapKeys) == 0 {
		ColorMapKeys = make([]string, 10, 10)
		ColorMapKeys[0] = "Black"
		ColorMapKeys[1] = "Red"
		ColorMapKeys[2] = "Green"
		ColorMapKeys[3] = "Yellow"
		ColorMapKeys[4] = "Blue"
		ColorMapKeys[5] = "Magenta"
		ColorMapKeys[6] = "Cyan"
		ColorMapKeys[7] = "White"
		ColorMapKeys[8] = "Bright"
		ColorMapKeys[9] = "Reset"

	}
}
