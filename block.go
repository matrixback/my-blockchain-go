package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {
	Timestamp 		int64
	Transactions    []*Transaction
	PrevBlockHash	[]byte
	Hash			[]byte
	Nonce 			int
}

// 创建一个块的时候计算 pow
func NewBlock(transactions []*Transaction, PrevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		transactions,
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

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
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
	NilPanic(err)
	return &block
}

// 在Pow 中把所有交易的hash取出来
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}