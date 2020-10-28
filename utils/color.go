package utils

import (
	"math/rand"

	"github.com/muesli/termenv"
)

var (
	term = termenv.ColorProfile()
)

// Color base 8color
// Black:   \u001b[30m
// Red:     \u001b[31m
// Green:   \u001b[32m
// Yellow:  \u001b[33m
// Blue:    \u001b[34m
// Magenta: \u001b[35m
// Cyan:    \u001b[36m
// White:   \u001b[37m
// Reset:   \u001b[0m
var Color = []string{
	"\u001b[31;1m",
	"\u001b[32;1m",
	"\u001b[33;1m",
	"\u001b[34;1m",
	"\u001b[35;1m",
	"\u001b[36;1m",
	"\u001b[37;1m",
}

func RandomColor(in string) (out string) {
	index := rand.Intn(len(Color))
	out = Color[index] + in
	return
}

func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}
