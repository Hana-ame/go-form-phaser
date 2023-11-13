package formphaser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestJson(t *testing.T) {
	m := make(map[string]any)
	err := json.Unmarshal([]byte(`{"a":1,"b":"string"}`), &m)
	fmt.Println(err)
}

func TestStackI(t *testing.T) {
	var s Stack
	var e string

	s.Push("123")
	fmt.Printf("%+v", s)
	s.Push("456")
	fmt.Printf("%+v", s)
	e = s.Pop()
	fmt.Printf("%+v", e)
	fmt.Printf("%+v", s)
	e = s.Pop()
	fmt.Printf("%+v", e)
	fmt.Printf("%+v", s)
	e = s.Pop()
	fmt.Printf("%+v", e)
	fmt.Printf("%+v", s)

}

func TestDecoder(t *testing.T) {
	type C struct {
		A int    `formphaser:"aaa"`
		B string `formphaser:"bbb"`
	}
	type O struct {
		A int      `formphaser:"aaa"`
		B string   `formphaser:"bbb"`
		C C        `formphaser:"ccc"`
		D []string `formphaser:"ddd"`
	}
	o := O{}
	Unmarshal(func(s string) []string { return nil }, &o)
	fmt.Printf("%+v\n", o)

}

func TestTuple(t *testing.T) {
	type C struct {
		A int    `formphaser:"aaa"`
		B string `formphaser:"bbb"`
	}
	type O struct {
		A int      `formphaser:"aaa"`
		B string   `formphaser:"bbb"`
		C C        `formphaser:"ccc"`
		D []string `formphaser:"ddd"`
	}
	o := O{}
	rv := reflect.ValueOf(o)
	tp := tuple{
		// reflect.TypeOf(v),
		reflect.StructField{},
		rv,
	}
	fmt.Printf("%+v\n", tp)
	fmt.Printf("%+v\n", tp.Field(0))
	fmt.Printf("%+v\n", tp.Field(1))
	fmt.Printf("%+v\n", tp.Field(2))
	fmt.Printf("%+v\n", tp.Field(3))
	fmt.Printf("%+v\n", tp)
}
