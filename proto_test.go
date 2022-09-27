package main

import (
	"hash/crc32"
	"testing"
)

func TestHash(t *testing.T) {
	crc32q := crc32.MakeTable(0xD5828281)
	t.Logf("%v", crc32.Checksum([]byte("1He3llo1 world"), crc32q))
}
