package utils

import "math/rand"

var StatusColor = map[string]string{
	"failed":   "\u001b[31;1mfailed\u001b[0m",
	"success":  "\u001b[32;1msuccess\u001b[0m",
	"pending":  "\u001b[33;1mpending\u001b[0m",
	"running":  "\u001b[34;1mrunning\u001b[0m",
	"skipped":  "\u001b[36;1mskipped\u001b[0m",
	"created":  "\u001b[37;1mcreated\u001b[0m",
	"canceled": "\u001b[37;1mcancel\u001b[0m",
	"manual":   "\u001b[37;1mmanual\u001b[0m",
}

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

func ColorStatus(status string) string {
	s, ok := StatusColor[status]
	if ok {
		return s
	}
	return status
}
