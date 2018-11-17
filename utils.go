package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}


func IsByteEmpty(b []byte) bool {
	return len(b) == 0
}

func NilPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}