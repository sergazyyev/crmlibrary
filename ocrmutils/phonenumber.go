package ocrmutils

import "strings"

func FormatPhoneNumber(phoneNumber string) string {
	phoneNumber = GetDigitsFromString(phoneNumber)
	if len(phoneNumber) == 0 {
		return ""
	}
	if phoneNumber[0:1] == "+" {
		phoneNumber = strings.Replace(phoneNumber, "+", "", 1)
	}
	if phoneNumber[0:1] == "8" {
		phoneNumber = ReplaceCharInStringAtIndex(phoneNumber, '7', 0)
	}
	return phoneNumber
}
