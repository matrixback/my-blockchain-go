package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

const subsidy = 10

type Transaction struct {
	ID 		[]byte
	Vin 	[]TXInput
	Vout 	[]TXOutput
}

type TXOutput struct {
	Value 		 int
	ScriptPubKey string // 锁定脚本或者地址可以解锁
}

type TXInput struct {
	Txid      []byte
	Vout 	  int
	ScriptSig string  // 签名用来解锁 input 引用的out，然后将值保存在 output 中
	                  // 矿工先验证这个 input 是否能解锁其引用的 output，如果成功，则加入
	                  // 此笔交易
}


func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 实际上 coinbase 的input 没什么用，有 output 就可以呢
	txin := TxInput([]byte{}, -1, data) // 不需要解锁脚本，写一段文字即可。解锁脚本主要是用于对
										// Output 的解锁，没有 output，则不需要
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// 只是计算 hash 值。需要先序列化为 []byte
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash[32] byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}