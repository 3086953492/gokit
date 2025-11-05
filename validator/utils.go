package validator

import (
	"regexp"
	"strings"
)

// CamelToSnake 驼峰命名转蛇形命名
// UsernameUnique -> username_unique
// XMLParser -> xml_parser
func CamelToSnake(s string) string {
	// 处理连续大写字母：XMLParser -> XmlParser
	re1 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	s = re1.ReplaceAllString(s, "${1}_${2}")

	// 处理普通驼峰：userName -> user_name
	re2 := regexp.MustCompile(`([a-z\d])([A-Z])`)
	s = re2.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(s)
}

