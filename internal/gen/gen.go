package gen

import (
	"crypto/sha256"
	rand2 "math/rand/v2"
)

const sm = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

func GeneratePassword(pw, site string, version string) string {
	hash := sha256.Sum256([]byte(pw + site + version))
	rnd := rand2.NewChaCha8(hash)
	buf := make([]byte, 16)
	n, err := rnd.Read(buf)
	if err != nil {
		panic(err)
	}
	for i, b := range buf[:n] {
		if int(b) >= len(sm) {
			buf[i] = sm[int(b)%len(sm)]
			continue
		}
		buf[i] = sm[b]
	}
	return string(buf)
}
