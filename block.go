package main

import (
	"time"
	"crypto/sha256"
	"strconv"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp 		int64
	Data			[]byte
	PrevBlockHash	[]byte
	Hash			[]byte
	Nonce 			int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		timestamp},
		[]byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// 创建一个块的时候计算 pow
func NewBlock(data string, PrevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		PrevBlockHash,
		[]byte{},
		0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	// 此处用 hash[:]，不直接用 hash 是为了安全，防止计算出的数据在某个地方被释放
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	// gob 解码时需要一个可读对象，通过 bytes.NewReader 转换一次
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
