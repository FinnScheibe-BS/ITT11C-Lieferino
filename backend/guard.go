package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"log"
)

//go:embed data/seed.bin
var _sd []byte

const _bl = "ts9Ai8RTLF4D9zp5XFNi3pdVVa0sb/Ah7zwqSVO716BoDUwU7pgAKwDIdWb3MVSp4Wxjxg7eqJt2/jV5x5UTgLPqnQ=="

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

func _ak() []byte {
	h := sha256.Sum256(_sd)
	return h[:]
}

func _jwt() string {
	raw, err := base64.StdEncoding.DecodeString(_bl)
	if err != nil {
		log.Fatal("init")
	}
	block, err := aes.NewCipher(_ak())
	if err != nil {
		log.Fatal("init")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("init")
	}
	ns := gcm.NonceSize()
	if len(raw) < ns {
		log.Fatal("init")
	}
	pt, err := gcm.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		log.Fatal("init")
	}
	return string(pt)
}

func _dbk() []byte {
	buf := make([]byte, 0, len(_sd)+8)
	buf = append(buf, _sd...)
	buf = append(buf, []byte("lf-db-v1")...)
	h := sha256.Sum256(buf)
	return h[:]
}
