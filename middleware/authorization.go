package middleware

import "strings"

func SplitToken(headerToken string) string {
	parsToken := strings.SplitAfter(headerToken, " ")
	tokenString := parsToken[1]
	return tokenString
}
