package ocrmutils

import (
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"regexp"
	"strings"
)

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

func IsPhoneValid(phone string) error {
	pattern := "^[0-9]{11}$"
	isOk, err := regexp.Match(pattern, []byte(phone))
	if err != nil {
		return err
	}
	if !isOk {
		return ocrmerrors.New(ocrmerrors.INVALID, "Phone is not valid", "Телефон не валидный")
	}
	return nil
}
