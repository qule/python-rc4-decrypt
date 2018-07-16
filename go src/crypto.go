package common

import (
	"encoding/hex"
	"errors"
)

var rc4NomalKey []uint8

var defaultRc4 *RC4Cipher

func init() {
	rc4NomalKey = make([]uint8, 256)
	for i := 0; i < 256; i++ {
		rc4NomalKey[i] = uint8(9*i + 7)
	}
	defaultRc4 = NewRC4Cipher()
}

type RC4Cipher struct {
	s    [256]uint8
	i, j uint8
}

func NewRC4Cipher() *RC4Cipher {
	var c RC4Cipher
	copy(c.s[:], rc4NomalKey)
	return &c
}

func (c *RC4Cipher) SetKey(key []byte) error {
	k := len(key)
	if k < 1 || k > 256 {
		return errors.New("eroor RC4Cipher key size")
	}
	copy(c.s[:], rc4NomalKey)
	c.i, c.j = 0, 0
	var j uint8 = 0
	for i := 0; i < 256; i++ {
		j += c.s[i] + key[i%k]
		c.s[i], c.s[j] = c.s[j], c.s[i]
	}
	return nil
}

func (c *RC4Cipher) AddKey(key []byte) error {
	k := len(key)
	if k < 1 || k > 256 {
		return errors.New("eroor RC4Cipher key size")
	}
	var j uint8 = 0
	for i := 0; i < 256; i++ {
		j += c.s[i] + key[i%k]
		c.s[i], c.s[j] = c.s[j], c.s[i]
	}
	return nil
}

func (c *RC4Cipher) StaticXorStream(dst, src []byte) {
	i, j := c.i, c.j
	for k, v := range src {
		i += 1
		j += c.s[i]
		dst[k] = v ^ c.s[c.s[i]+c.s[j]]
	}
}

func (c *RC4Cipher) XorStream(dst, src []byte) {
	i, j := c.i, c.j
	for k, v := range src {
		i += 1
		j += c.s[i]
		c.s[i], c.s[j] = c.s[j], c.s[i]
		dst[k] = v ^ c.s[c.s[i]+c.s[j]]
	}
	c.i, c.j = i, j
}

func ResetRc4Key(key []byte) error {
	return defaultRc4.SetKey(key)
}

func StaticXorStream(dst, src []byte) {
	if dst == nil {
		dst = make([]byte, len(src))
	}
	defaultRc4.StaticXorStream(dst, src)
}

func Rc4Encode(src []byte) string {
	dst := make([]byte, len(src))
	defaultRc4.StaticXorStream(dst, src)
	return hex.EncodeToString(dst)
}

func Rc4Decode(src string) ([]byte, error) {
	srcByte, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(srcByte))
	defaultRc4.StaticXorStream(dst, srcByte)
	return dst, nil
}
