# go-form-phaser

从`key[key][key]`like的键名传到**struct**结构

**struct**为主

1. 怎么access pointer里的内容啊
2. 遍历keys
3. 输出应该有的key
4. 查找

## json.Unmarshal

非常参考json的做法

1. 要检查是不是pointer

### note: 如何使得自定义struct能够marshal
```go
// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
```