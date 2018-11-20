package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"bytes"

	"golang.org/x/crypto/ripemd160"

)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte  // 只是 x，y 数组
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}


// 为什么不把这个钱包地址放到 Wallet 里面去，每次要算一遍？
func (w Wallet) getAddress() []byte {
	pubKeyhash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{version}, pubKeyhash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}


func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	NilPanic(err)

	publicRIPED160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPED160
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(apppend([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	sendcodSHA := sha256.Sum256(firstSHA[:])

	return sendcodSHA[:addressChecksumLen]
}


func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	NilPanic(err)

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
