package ocrmutils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestFormatPhoneNumber(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{
			name: "plus",
			value: "+77078338948",
		},
		{
			name: "skobka",
			value: "+ 7 (707) 83 38 94 8 ",
		},
		{
			name: "probel",
			value: "8 (707) - 833 - 89 - 4 - 8 ",
		},
		{
			name: "vosem",
			value: "8 707 833 89-48",
		},
		{
			name: "online",
			value: "+7(707)-833-8948",
		},
	}
	for _, cas := range testCases {
		t.Run(cas.name, func(t *testing.T) {
			assert.Equal(t, "77078338948", FormatPhoneNumber(cas.value))
		})
	}
}

func TestFormatPhoneNumber2(t *testing.T) {
	phone := "+7(701)-981-0400"
	re := regexp.MustCompile("[0-9]+")
	phone = strings.Join(re.FindAllString(phone, -1), "")
	t.Log(phone)
}
