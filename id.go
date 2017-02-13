package wecms

import (
	"strings"
	"crypto/rand"
	"io"
	"encoding/hex"
)

type ID string

var idRanded = rand.Reader

const RootID ID = "11111111-1111-1111-1111-111111111111"

func newIdStr() string {
	var idBytes [16]byte
	io.ReadFull(idRanded, idBytes[:])
	idBytes[6] = (idBytes[6] & 0x0f) | 0x40
	idBytes[8] = (idBytes[8] & 0x3f) | 0x80
	var strBuf [36]byte
	hex.Encode(strBuf[:], idBytes[:4])
	strBuf[8] = '-'
	hex.Encode(strBuf[9:13], idBytes[4:6])
	strBuf[13] = '-'
	hex.Encode(strBuf[14:18], idBytes[6:8])
	strBuf[18] = '-'
	hex.Encode(strBuf[19:23], idBytes[8:10])
	strBuf[23] = '-'
	hex.Encode(strBuf[24:], idBytes[10:])
	return string(strBuf)
}

// NewID create a new random ID
func NewID() ID {
	return ID(strings.ToLower(newIdStr()))
}

// Eq check one id equal to another(escape cases)
func (id ID) Eq(newId ID) bool {
	return strings.EqualFold(string(id), string(newId))
}