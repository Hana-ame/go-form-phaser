package formphaser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// stack
// tested
type Stack []string

func (s *Stack) Push(e string) {
	*s = append((*s), e)
}
func (s *Stack) Pop() (e string) {
	if len(*s) == 0 {
		return
	}
	e = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return
}
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}
func (s *Stack) Length() int {
	return len(*s)
}
func (s *Stack) Peek() (e string) {
	if s.IsEmpty() {
		return
	}
	e = (*s)[s.Length()-1]
	return
}
func (s *Stack) ToSlice() []string {
	return []string(*s)
}

type Index []int

func (s *Index) Push(e int) {
	*s = append((*s), e)
}
func (s *Index) Pop() (e int) {
	if len(*s) == 0 {
		return
	}
	e = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return
}
func (s *Index) IsEmpty() bool {
	return len(*s) == 0
}
func (s *Index) Length() int {
	return len(*s)
}
func (s *Index) Peek() (e int) {
	if s.IsEmpty() {
		return
	}
	e = (*s)[s.Length()-1]
	return
}
func (s *Index) ToSlice() []int {
	return []int(*s)
}

// error
// copy form json
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "formphaser: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "formphaser: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "formphaser: Unmarshal(nil " + e.Type.String() + ")"
}

type tuple struct {
	// reflect.Type
	reflect.StructField
	reflect.Value
}

// must be struct when call it
func (t *tuple) Field(i int) tuple {
	return tuple{
		// Type:        uncompressingType(t.Type.Field(i).Type),
		StructField: uncompressingType(t.Type()).Field(i),
		Value:       uncompressingValue(t.Value.Field(i)),
	}
}
func (t *tuple) Type() reflect.Type {
	return uncompressingType(t.Value.Type())
}
func (t *tuple) Tag() string {
	return t.StructField.Tag.Get("formphaser")
}
func (t *tuple) Kind() reflect.Kind {
	return t.Value.Kind()
}

// decoder
type decoder struct {
	Stack
	getter func(string) []string
}

func (d *decoder) get(key string) []string {
	return d.getter(key)
}

func (d *decoder) setstruct(t tuple) error {
	// v := uncompressing(t.Value)
	// must call when struct

	for i := 0; i < t.NumField(); i++ {
		subt := t.Field(i)

		d.Push(subt.Tag())
		path := stack2path(d.Stack)

		fmt.Println(path)
		fmt.Println(subt.Value)
		fmt.Println(subt.Type())

		switch subt.Kind() {
		default:
			fmt.Printf("not support %+v", subt.Type())
		case reflect.Struct:
			e := d.setstruct(subt)
			if e != nil {
				return e
			}
		case reflect.Bool:
			// if get nothing
			if d.getter(path) == nil {
				d.Pop()
				continue
			}
			subt.SetBool(isTrue(d.get(path)[0]))
		case reflect.Int:
			if d.getter(path) == nil {
				d.Pop()
				continue
			}
			x, e := strconv.Atoi(d.get(path)[0])
			if e != nil {
				return e
			}
			subt.SetInt(int64(x))
		case reflect.String:
			// if d.getter(path) == nil {
			// 	d.Pop()
			// 	continue
			// }
			// t.SetString(d.get(path)[0])
			subt.SetString(path)
		case reflect.Slice:
			x := reflect.ValueOf(d.Stack.ToSlice())
			subt.Set(x)
		}

		//
		d.Pop()
	}

	// if v.Kind() == reflect.Struct {
	// 	for i := 0; i < v.NumField(); i++ {
	// 		d.setvalue(t.Field(i))
	// 	}
	// } else if v.Kind() == reflect.String {

	// } else if v.Kind() == reflect.Slice { // terminated for the timebeing

	// } else if v.Kind() == reflect.Bool {

	// } else if v.Kind() == reflect.Int {

	// } else {
	// 	return &InvalidUnmarshalError{reflect.TypeOf(v)}
	// }

	return nil
}

func (d *decoder) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	rv = uncompressingValue(rv)

	// only support struct
	if rv.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	t := tuple{
		// reflect.TypeOf(v),
		reflect.StructField{},
		rv,
	}
	d.setstruct(t)

	return nil
}

// entrance
// @input v 指向 struct 的指针
// @input getter 从 key 取出 []string 的方法
// 带 array 的直接用 []string, 注释含有`[]`就可以
// 没特别测试过, 之后添加
func Unmarshal(getter func(string) []string, v any) error {

	d := decoder{
		Stack:  Stack{},
		getter: getter,
	}

	e := d.unmarshal(v)

	return e
}

// helper

// 返回不是pointer的最终对象
func uncompressingValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}
func uncompressingType(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}

// misc
func stack2path(s Stack) (path string) {
	if s.Length() == 1 {
		return s.Peek()
	}
	for i, str := range s.ToSlice() {
		if i == 0 {
			path += str
		} else if strings.Contains(str, "[]") {
			path += str
		} else {
			path += "[" + str + "]"
		}
	}
	return
}

func isTrue(s string) bool {
	r := s == "true" || s == "1"
	return r
}
