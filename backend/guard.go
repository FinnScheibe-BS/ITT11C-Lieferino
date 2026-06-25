package main

import (
	_ "embed"
	"encoding/base64"
	"log"
)

//go:embed data/seed.bin
var _sd []byte

func _fn(b []byte) uint32 {
	h := uint32(0x811c9dc5)
	for _, x := range b {
		h ^= uint32(x)
		h *= 0x01000193
	}
	return h
}

func _ks(seed uint32, n int) []byte {
	x := seed
	if x == 0 {
		x = 0x9e3779b9
	}
	o := make([]byte, n)
	for i := 0; i < n; i++ {
		x ^= x << 13
		x ^= uint32(int32(x) >> 17)
		x ^= x << 5
		o[i] = byte(x & 0xff)
	}
	return o
}

func _g() {
	c, err := base64.StdEncoding.DecodeString("h5iZcmvosUGgrnZ4")
	if err != nil {
		log.Fatal("init")
	}
	k := _ks(_fn(_sd), len(c))
	o := make([]byte, len(c))
	for i := range c {
		o[i] = c[i] ^ k[i]
	}
	if string(o) != "LIEFERINO-OK" {
		log.Fatal("init")
	}
}
