package ocrmutils

import (
	"fmt"
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"math/big"
	"strconv"
)

const (
	ibannumber_min_size     = 15
	ibannumber_max_size     = 34
	ibannumber_magic_number = 97
)

func CheckIban(iban string) error {
	if iban[0:2] != "KZ" {
		return ocrmerrors.New(ocrmerrors.INVALID, "Iban must begin at KZ", "Iban должен начинаться с сиволов KZ")
	}
	if len(iban) < ibannumber_min_size {
		return ocrmerrors.New(ocrmerrors.INVALID, fmt.Sprintf("Iban lenth must be more than %d", ibannumber_min_size), fmt.Sprintf("Длина Iban должно быть больше %d символов", ibannumber_min_size))
	}
	if len(iban) > ibannumber_max_size {
		return ocrmerrors.New(ocrmerrors.INVALID, fmt.Sprintf("Iban lenth must be less than %d", ibannumber_max_size), fmt.Sprintf("Длина Iban должно быть меньше %d символов", ibannumber_max_size))
	}

	iban = iban[4:] + iban[:4]
	mods := ""

	for _, v := range iban {
		i := int(v)
		// Check if c is characters A-Z (codepoint 65 - 90 in ASCII)
		if i > 64 && i < 91 {

			//In ascii characters starts 65 and ends 90
			//To get A=10, B=11, C=12 etc...
			//Minus 55 from ascii code
			i -= 55
			mods += strconv.Itoa(i)
		} else {
			//Or get own value as int value of character (to be digits)
			mods += string(i)
		}
	}

	bigVal, succees := new(big.Int).SetString(mods, 10)
	if !succees {
		return ocrmerrors.New(ocrmerrors.INVALID, "Iban check digits validation failed", "Введенный Iban не верный")
	}

	modVal := new(big.Int).SetInt64(ibannumber_magic_number)
	resVal := new(big.Int).Mod(bigVal, modVal)

	if resVal.Int64() != 1 {
		return ocrmerrors.New(ocrmerrors.INVALID, "Iban check digits validation failed", "Введенный Iban не верный")
	}
	return nil
}
