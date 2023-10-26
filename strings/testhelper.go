package strings

import (
	"fmt"
	"strings"
)

const debug = false

type randomSource string

func (s randomSource) get(seed int, length int) string {
	const seedShift = 97
	if length == 0 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(s[(seedShift+seed+i)%len(s)])
	}
	return b.String()
}

var (
	asciiSource  = randomSource("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	asciiSources = map[string]randomSource{
		"ascii": asciiSource,
	}
	utf8Source = randomSource("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzΑαАаاבאअब你あア₽€£¥©®§—–‘’“”•…†✓★🌞❤😊👍👎🔥⚡❄🇺🇸🌸🎵⚽🌙🎂✌☮🍔🍕🎅💃👰👏👋👌✍👍🏼☠☃🐉👑🌪👽🚀🥅🌍📷✏⚓⌚🌻☂🎑💣🏔🎤🌈😺🐶🦓🦉🐵🌴🍺🎲🎾👟🚲🍦⏰🎆😍💔🎅👻🎬🌎🍷")

	allSources = map[string]randomSource{
		"ascii": asciiSource,
		"utf8":  utf8Source,
	}
)

func testName(args ...any) string {
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = fmt.Sprintf("%v", arg)
	}
	return strings.Join(argStrs, "/")
}
