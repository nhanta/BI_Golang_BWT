package ConverttoByte

import "strconv"

func InttoByte (integer int) byte {
	val := strconv.Itoa(integer)
	i, _ := strconv.Atoi(val)
	byteI := byte(i)
	return byteI
}

