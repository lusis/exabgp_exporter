package utils

import (
	"regexp"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

// StripAnsi removes all ansi sequences from a string
// sometimes exabgp has bad json and this is the best route until fixes are made upstream
// from here: https://github.com/acarl005/stripansi
func StripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}
