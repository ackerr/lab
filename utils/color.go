package utils

import "math/rand"

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
