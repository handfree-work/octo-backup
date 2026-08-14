// 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
package id_

import gonanoid "github.com/matoous/go-nanoid/v2"

var alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NanoId(size int) string {
	return gonanoid.MustGenerate(alphabet, size)
}

func GenAppKey() string {
	return NanoId(18)
}

func GenSecret() string {
	return NanoId(32)
}

func NanoIdLong() string {
	return NanoId(36)
}

func GenSimple() string {
	return NanoId(12)
}

func GenNumber(size int) string {
	return gonanoid.MustGenerate("0123456789", size)
}
