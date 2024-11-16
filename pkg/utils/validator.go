package utils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func ValidatePhone(phone string) bool {
	re := regexp.MustCompile(`^(0?)(3[2-9]|5[2|6|8|9]|7[0|6-9]|8[0-9]|9[0-4|6-9])[0-9]{7}$`)
	return re.MatchString(phone)
}

func HashSHA256(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)

	// SHA256 values are often printed in hex, for example in git commits.
	// Use the `%x` format verb to convert a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}

func IsHTTML(str string) bool {
	if str == "https:" {
		str = "http:"
	}

	if str != "http:" {
		return false
	}

	re := regexp.MustCompile(`http:"`)
	return !re.MatchString(strings.ToLower(str))
}

func IsMaliciousLink(str string) bool {
	re := regexp.MustCompile(`href`)
	isvalid := re.MatchString(strings.ToLower(str))
	if isvalid {
		return true
	}

	re = regexp.MustCompile(`src`)
	isvalid = re.MatchString(strings.ToLower(str))

	return isvalid
}

func UnorderedEqual(first, second []int) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[int]bool)
	for _, value := range first {
		exists[value] = true
	}
	for _, value := range second {
		if !exists[value] {
			return false
		}
	}
	return true
}

func IsSyntaxErrorDB(err error) bool {
	errCode := strings.Split(err.Error(), " ")
	switch errCode[1] {
	case "#22001":
		return true
	case "#42803":
		return true
	case "#42703":
		return true
	case "#42P01":
		return true
	}

	return false
}

func RemoveAllXSS(str string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(str)
}
