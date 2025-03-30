package stu

import (
	"bytes"
	"strconv"
)

func JoinIntArrayToString(arr []int, delimiter string) string {
	var buffer bytes.Buffer
	var l = len(arr)

	for i := range l - 1 {
		buffer.WriteString(strconv.Itoa(arr[i]))
		buffer.WriteString(delimiter)
	}
	buffer.WriteString(strconv.Itoa(arr[l]))

	return buffer.String()
}

func JoinUInt64ArrayToString(arr []uint64, base int, delimiter string) string {
	var buffer bytes.Buffer
	var l = len(arr)

	for i := range l - 1 {
		buffer.WriteString(strconv.FormatUint(arr[i], base))
		buffer.WriteString(delimiter)
	}
	buffer.WriteString(strconv.FormatUint(arr[l], base))

	return buffer.String()
}

func JoinWithSurroundings(arr []string, delimiter string, prefix string, suffix string) string {
	var buffer bytes.Buffer
	var l = len(arr)

	for i := range l - 1 {
		buffer.WriteString(prefix)
		buffer.WriteString(arr[i])
		buffer.WriteString(suffix)
		buffer.WriteString(delimiter)
	}
	buffer.WriteString(prefix)
	buffer.WriteString(arr[l-1])
	buffer.WriteString(suffix)

	return buffer.String()
}
