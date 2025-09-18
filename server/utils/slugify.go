package utils

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

func Slugify(s string) string {
	s = strings.ToLower(s)

	p := pinyin.NewArgs()
	pinyins := pinyin.LazyPinyin(s, p)
	s = strings.Join(pinyins, " ")

	s = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.Join(strings.Fields(s), "-")
	return s
}
