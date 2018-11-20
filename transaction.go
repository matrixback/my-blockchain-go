package main

/*
注意锁定和解锁的过程，锁定时，将一个 PubKeyHash 加入某个 output，
然后到了解锁，必须要有 PubKey 来证明可以解锁这个交易，并且还要有个 signature,
来证明这个 input 有效。

注意，一笔交易就是输入和输出，对输出没有什么要求，只要总是不大于输入。
而对输入要求比较高，input 要引用一些 output，那么必须证明自己是这些 output 的拥有者。
证明之一是 PubKey 符合，第二是，你拥有这个 PubKey 的私钥。怎么证明你拥有私钥？用签名。
把一些 data 进行签名，然后写入 input，并且可以用这个 PubKey 解开。
*/

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
)

const subsidy = 10

type Transaction struct {
	ID 		[]byte
	Vin 	[]TXInput
	Vout 	[]TXOutput
}


// 注意 output 和 input，一个是 pubkeyHash, 一个是 pubkey
type TXOutput struct {
	Value 		 int
	PubKeyHash   []byte // 公钥的哈希，只有特定的公钥才能解锁。
	                   // 用 hash 而不直接用 key，多了一层保障。不能很明显的看出一个这笔 out 属于某个 address
}

// 一笔交易的第几个输出
type TXInput struct {
	Txid      []byte
	Vout 	  int
	PubKey    []byte  // 用这个公钥来解锁
	Signature []byte  // 签名用来解锁 input 引用的out，然后将值保存在 output 中
	                  // 矿工先验证这个 input 是否能解锁其引用的 output，如果成功，则加入
	                  // 此笔交易
}


// 判断这个公钥可以解开 output 的公钥哈希值，即锁定脚本
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}


func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-1]
	out.PubKeyHash = pubKeyHash
}


func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}


func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// 实际上 coinbase 的input 没什么用，有 output 就可以呢
	txin := TXInput{[]byte{}, -1, data} // 不需要解锁脚本，写一段文字即可。解锁脚本主要是用于对
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

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 &&
		   len (tx.Vin[0].Txid) == 0 &&
		   tx.Vin[0].Vout == -1
}

// 下面的两个函数功能差不多，一个判断是某个 input 可以解锁一个字符串，
// 一个判断某个output 是否可以被一个字符串解锁
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}


func NewUTXOTransaction(from ,to string, amount int, bc *Blockchain)*Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	for txid, outs:= range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc-amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}


// prevTXs 是所有的前面的交易？
// 对这一笔交易进行签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		NilPanic(err)

		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

// 返回值为什么不是指针
// 反正只是一个签名，所以只保留每笔交易的核心值
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	// vout 原样输入
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
