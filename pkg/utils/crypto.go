// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/sha3"
)

func BoxEnc(b []byte, pk *[32]byte) []byte {
	plainB, err := box.SealAnonymous(nil, b, pk, rand.Reader)
	if err != nil {
		panic(err)
	}
	return plainB
}

func BoxDec(b []byte, pk *[32]byte, sk *[32]byte) ([]byte, bool) {
	return box.OpenAnonymous(nil, b, pk, sk)
}

func Hash(b []byte) []byte {
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, b)
	return hash
}
