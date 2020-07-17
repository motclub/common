package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func MD5(p string, upper ...bool) (v string) {
	h := md5.New()
	h.Write([]byte(p))
	v = hex.EncodeToString(h.Sum(nil))
	if len(upper) > 0 && upper[0] {
		v = strings.ToUpper(v)
	}
	return
}

func Sha1(p string, upper ...bool) (v string) {
	d := sha1.New()
	d.Write([]byte(p))
	v = hex.EncodeToString(d.Sum([]byte(nil)))
	if len(upper) > 0 && upper[0] {
		v = strings.ToUpper(v)
	}
	return
}
