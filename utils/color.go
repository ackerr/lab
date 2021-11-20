package utils

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/muesli/termenv"
)

var term = termenv.ColorProfile()

var Color = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgWhite,
}

// PrintlnWithColor fix the color of the output on Windows.
func PrintlnWithColor(a ...interface{}) {
	fmt.Fprintln(colorable.NewColorableStdout(), a...)
}

func RandomColor(in string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(Color))
	return color.New(Color[index]).Sprintf(in)
}

func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func ColorBg(val, color string) string {
	return termenv.String(val).Background(term.Color(color)).String()
}
