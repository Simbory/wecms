package wecms

import (
	"github.com/google/uuid"
	"strings"
)

type ID string

// NewID create a new random ID
func NewID() ID {
	id,_ := uuid.NewRandom()
	return ID(id.String())
}

// Eq check one id equal to another(escape cases)
func (id ID) Eq(newId ID) bool {
	return strings.EqualFold(string(id), string(newId))
}

const RootID ID = "11111111-1111-1111-1111-111111111111"