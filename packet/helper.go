package packet

import "github.com/wizcas/mudever.svc/telbyte"

func insert(dst []byte, b byte, pos int) []byte {
	dst = append(dst, 0)
	copy(dst[pos+1:], dst[pos:])
	dst[pos] = b
	return dst
}

func escapeData(data []byte) []byte {
	result := data[:]
	for i := 0; i < len(result); i++ {
		b := result[i]
		if b == byte(telbyte.IAC) {
			result = insert(result, byte(telbyte.IAC), i)
			i++
		}
	}
	return result
}
