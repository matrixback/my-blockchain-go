1. []byte

go 的 byte 是什么类型？字符类型？

uint8 类型，二进制的01映射为 ASCII 码。

既然 string 就是一系列字节，而 []byte 也可以表达一系列字节，那么实际运用中应当如何取舍？

- string 值不可修改，所以每次修改时会创建一个新的 string，要想原地修改就只能用 []byte
- string 值不可为 nil，所以如果你想要通过返回 nil 表达点什么额外的含义，就只能用 []byte
- []byte 可以是数组也可以是切片，Go 语言的切片这么灵活，想要用切片的特性就只能用 []byte
- 如果你要调用第三方库，如果第三方库主要用的是 string，那么为了避免频繁的进行类型转换，
那你就可以用 string。毕竟 Go 是强类型语言，类型转换会极大地降低代码的可读性。[]byte 也是同理。
需要进行字符串处理的时候，因为 []byte 支持切片，并且可以原地修改，所以 string 更快一些，
所以注重性能的地方你可以用 []byte

conclude: byte 修改起来比较灵活，string 可看做是 const 语义。

在 block-go 中，底层的存储都用 []byte，而每个交易的 data 数据为 string。

    bc.AddBlock("Send 1 BTC to Ivan")

这样的话，计算工作量证明也比较好算

    func (pow *ProofOfWork) prepareData(nonce int) []byte {
        data := bytes.Join(
            [][]byte{
                pow.block.PrevBlockHash,
                pow.block.Data,
                IntToHex(pow.block.Timestamp),
                IntToHex(int64(targetBits)),
                IntToHex(int64(nonce)),
            },
            []byte{},
        )

        return data
    }

2. 取出任何值的二进制，然后转化为 bytes

    import (
        "bytes"
        "encoding/binary"
        "log"
    )

    func IntToHex(num int64) []byte {
        buff := new(bytes.Buffer)
        err := binary.Write(buff, binary.BigEndian, num)
        if err != nil {
            log.Panic(err)
        }

        return buff.Bytes()
    }