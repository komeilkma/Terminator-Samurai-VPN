package cipher
var _key = []byte("12312300TSVPN@@!!2022")

func XOR(src []byte) []byte {
	_klen := len(_key)
	for i := 0; i < len(src); i++ {
		src[i] ^= _key[i%_klen]
	}
	return src
}
func SetKey(key string) {
	_key = []byte(key)
}

