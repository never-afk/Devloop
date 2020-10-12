package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sony/sonyflake"
	"hash/crc32"
	"strconv"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{})

func UniqueId() string {
	id, _ := sf.NextID()
	return strconv.FormatUint(id, 10)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//16位md5
func M16(str string) string {
	return Md5(str)[8:24]
}

func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

func IntVal(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}
