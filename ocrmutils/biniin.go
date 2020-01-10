package ocrmutils

import (
	"fmt"
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"strconv"
	"strings"
)

func IsBinIinValid(binIin string) error {
	if len(binIin) != 12 {
		return ocrmerrors.New(ocrmerrors.INVALID, "length of bin/iin must be equal 12", "В БИН/ИИН количество символов должно быть равно 12")
	}
	ves1 := []int{1,2,3,4,5,6,7,8,9,10,11}
	ves2 := []int{3,4,5,6,7,8,9,10,11,1,2}
	var razrBinIin [12]int
	iinArr := strings.Split(binIin, "")
	sum := 0

	for i, k := range iinArr {
		val, err := strconv.Atoi(k)
		if err != nil {
			return ocrmerrors.New(ocrmerrors.INVALID, "symbols of bin/iin must be numeric", "БИН/ИИН должен состоять только из цифр")
		}
		razrBinIin[i] = val
	}

	for i1, k1 := range ves1 {
		sum = sum + k1*razrBinIin[i1]
	}

	ctrlSum := sum%11

	if ctrlSum == 10 {
		sum = 0
		for i2, k2 := range ves2 {
			sum = sum + k2*razrBinIin[i2]
		}
		ctrlSum = sum%11
	}
	eq := razrBinIin[11]==ctrlSum
	if !eq {
		return ocrmerrors.New(ocrmerrors.INVALID, fmt.Sprintf("invalid bin/iin %s", binIin) , fmt.Sprintf("БИН/ИИН %s не валидный", binIin))
	} else {
		return nil
	}
}