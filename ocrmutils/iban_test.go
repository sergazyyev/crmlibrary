package ocrmutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	iban = "KZ909470398926234633"
	iniban = "KZ909470398926234611"
)

func TestCheckIban(t *testing.T) {
	assert.NoError(t, CheckIban(iban))
	assert.Error(t, CheckIban(iniban))
}
