/*
Copyright 2023 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pack

import (
	"bytes"
	"errors"
	"reflect"
	"runtime"
)

var (
	errEmptyKey      = errors.New("empty key")
	errUnexpectedEnd = errors.New("unexpected end")
)

func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}

type Unmarshaler interface {
	UnmarshalKSPACK([]byte) error
}

type decodeState struct {
	data       []byte
	off        int
	savedError error
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.savedError = nil
	return d
}

func (d *decodeState) error(err error) {
	panic(err)
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d.value(rv)
	if d.savedError != nil {
		return d.savedError
	}
	if d.off != len(d.data) {
		return errUnexpectedEnd
	}
	return nil
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so that
// it can be set to nil.
func (d *decodeState) indirect(v reflect.Value, decodingNull bool) (Unmarshaler, reflect.Value) {
	if !v.IsValid() {
		return nil, reflect.Value{}
	}
	// If v is a named type and is addressable
	// start with its address, so that is the type has pointer
	// methods, we find them
	if v.CanAddr() && v.Kind() != reflect.Ptr && v.Type().Name() != "" {
		v = v.Addr()
	}

	for {
		// Load value from interface, but only if the result will be
		// usefully addressable
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}

		if v.IsNil() {
			// nil pointer
			if d.data[d.off] == KSPACK_NULL {
				return nil, v
			}
			v.Set(reflect.New(v.Type().Elem()))
		}

		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		v = v.Elem()
	}
	return nil, v
}

func (d *decodeState) value(v reflect.Value) {
	if !v.IsValid() {
		d.next()
		return
	}

	u, pv := d.indirect(v, false)
	if u != nil {
		if err := u.UnmarshalKSPACK(d.next()); err != nil {
			d.error(err)
		}
		return
	}

	v = pv

	switch d.data[d.off] {
	case KSPACK_OBJECT:
		d.object(v)
	case KSPACK_ARRAY:
		d.array(v)
	case KSPACK_STRING:
		d.string(v)
	case KSPACK_SHORT_STRING:
		d.shortString(v)
	case KSPACK_BINARY:
		d.binary(v)
	case KSPACK_SHORT_BINARY:
		d.shortBinary(v)
	case KSPACK_INT8:
		d.int8(v)
	case KSPACK_INT16:
		d.int16(v)
	case KSPACK_INT32:
		d.int32(v)
	case KSPACK_INT64:
		d.int64(v)
	case KSPACK_UINT8:
		d.uint8(v)
	case KSPACK_UINT16:
		d.uint16(v)
	case KSPACK_UINT32:
		d.uint32(v)
	case KSPACK_UINT64:
		d.uint64(v)
	case KSPACK_BOOL:
		d.bool(v)
	case KSPACK_FLOAT:
		d.float(v)
	case KSPACK_DOUBLE:
		d.double(v)
	case KSPACK_NULL:
		d.null(v)
	}
}

func (d *decodeState) next() []byte {
	start := d.off
	typ := d.data[d.off]
	d.off += 1
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1
	vlen := 0
	switch typ {
	case KSPACK_OBJECT:
		vlen = int(Uint32(d.data[d.off:]))
		d.off += 4
	case KSPACK_ARRAY:
		vlen = int(Uint32(d.data[d.off:]))
		d.off += 4
	case KSPACK_STRING:
		vlen = int(Uint32(d.data[d.off:]))
		d.off += 4
	case KSPACK_SHORT_STRING:
		vlen = int(Uint8(d.data[d.off:]))
		d.off += 1
	case KSPACK_BINARY:
		vlen = int(Uint32(d.data[d.off:]))
		d.off += 4
	case KSPACK_SHORT_BINARY:
		vlen = int(Uint8(d.data[d.off:]))
		d.off += 1
	case KSPACK_INT8:
		vlen = 1
	case KSPACK_INT16:
		vlen = 2
	case KSPACK_INT32:
		vlen = 4
	case KSPACK_INT64:
		vlen = 8
	case KSPACK_UINT8:
		vlen = 1
	case KSPACK_UINT16:
		vlen = 2
	case KSPACK_UINT32:
		vlen = 4
	case KSPACK_UINT64:
		vlen = 8
	case KSPACK_BOOL:
		vlen = 1
	case KSPACK_FLOAT:
		vlen = 4
	case KSPACK_DOUBLE:
		vlen = 8
	case KSPACK_DATE:
		// FIXME
	case KSPACK_NULL:
		vlen = 1
	}
	d.off += klen + vlen
	return d.data[start:d.off]
}

// type(1) | name length(1) | content length (4)
// | raw name bytes | 0x00 | content bytes | 0x00
func (d *decodeState) string(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	val := string(d.data[d.off : d.off+vlen-1])
	d.off += vlen // value and 0x00

	v.SetString(val)
}

func (d *decodeState) stringInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	val := string(d.data[d.off : d.off+vlen-1])
	d.off += vlen // value and 0x00

	return val
}

// type(1) | name length(1) | content length(1) | raw name bytes |
// 0x00 | content bytes | 0x00
func (d *decodeState) shortString(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint8(d.data[d.off:]))
	d.off += 1 // content length

	d.off += klen // name and 0x00

	val := string(d.data[d.off : d.off+vlen-1])
	d.off += vlen // value and 0x00

	v.SetString(val)
}

func (d *decodeState) shortStringInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint8(d.data[d.off:]))
	d.off += 1 // content length

	d.off += klen // name and 0x00

	val := string(d.data[d.off : d.off+vlen-1])
	d.off += vlen // value and 0x00

	return val
}

// type(1) | name length(1) | content length(4) | raw name bytes |
// 0x00 | content bytes
func (d *decodeState) binary(v reflect.Value) {
	d.off += 1 //type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	val := d.data[d.off : d.off+vlen]
	d.off += vlen // value

	v.SetBytes(val)
}

func (d *decodeState) binaryInterface() interface{} {
	d.off += 1 //type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	val := d.data[d.off : d.off+vlen]
	d.off += vlen // value

	return val
}

// type(1) | name length(1) | content length(1) | raw name bytes |
// 0x00 | content bytes
func (d *decodeState) shortBinary(v reflect.Value) {
	d.off += 1 //type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint8(d.data[d.off:]))
	d.off += 1 // content length

	d.off += klen // name and 0x00

	val := d.data[d.off : d.off+vlen]
	d.off += vlen // value

	v.SetBytes(val)
}

func (d *decodeState) shortBinaryInterface() interface{} {
	d.off += 1 //type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	vlen := int(Uint8(d.data[d.off:]))
	d.off += 1 // content length

	d.off += klen // name and 0x00

	val := d.data[d.off : d.off+vlen]
	d.off += vlen // value

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(1)
func (d *decodeState) int8(v reflect.Value) {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Int8(d.data[d.off:])
	d.off += 1 //value
	v.SetInt(int64(val))
}
func (d *decodeState) int8Interface() interface{} {
	// unsupported in libkspack, int32 employed
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Int8(d.data[d.off:])
	d.off += 1 // value
	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(1)
func (d *decodeState) uint8(v reflect.Value) {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Uint8(d.data[d.off:])
	d.off += 1 // value
	v.SetUint(uint64(val))
}
func (d *decodeState) uint8Interface() interface{} {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Uint8(d.data[d.off:])
	d.off += 1 // value
	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(2)
func (d *decodeState) int16(v reflect.Value) {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Int16(d.data[d.off:])
	d.off += 2 // value
	// unsupported in libkspack, int32 employed
	//d.off += 4
	v.SetInt(int64(val))
}
func (d *decodeState) int16Interface() interface{} {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Int16(d.data[d.off:])
	d.off += 2 // value
	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(2)
func (d *decodeState) uint16(v reflect.Value) {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Uint16(d.data[d.off:])
	d.off += 2 // value
	// unsupported in libkspack, int32 employed
	//d.off += 4
	v.SetUint(uint64(val))
}
func (d *decodeState) uint16Interface() interface{} {
	d.off += 1 // type
	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length
	d.off += klen
	val := Uint16(d.data[d.off:])
	d.off += 2 // value
	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(4)
func (d *decodeState) int32(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Int32(d.data[d.off:])
	d.off += 4 // value

	v.SetInt(int64(val))
}

func (d *decodeState) int32Interface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Int32(d.data[d.off:])
	d.off += 4 // value

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(4)
func (d *decodeState) uint32(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Uint32(d.data[d.off:])
	d.off += 4 // value

	v.SetUint(uint64(val))
}

func (d *decodeState) uint32Interface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Uint32(d.data[d.off:])
	d.off += 4 // value

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(8)
func (d *decodeState) int64(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Int64(d.data[d.off:])
	d.off += 8 // value

	v.SetInt(val)
}

func (d *decodeState) int64Interface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Int64(d.data[d.off:])
	d.off += 8 // value

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(8)
func (d *decodeState) uint64(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Uint64(d.data[d.off:])
	d.off += 8 // value

	v.SetUint(val)
}

func (d *decodeState) uint64Interface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Uint64(d.data[d.off:])
	d.off += 8 // value

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | 0x00
func (d *decodeState) null(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	d.off += 1 // value

	v.Set(reflect.Zero(v.Type()))
}

// type(1) | name length(1) | raw name bytes | 0x00 | 0x00
func (d *decodeState) nullInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	d.off += 1 // value

	return nil
}

// type(1) | name length(1) | raw name bytes | 0x00 | 0x00/0x01
func (d *decodeState) bool(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := d.data[d.off]
	d.off += 1

	if val == 0 {
		v.SetBool(false)
	} else {
		v.SetBool(true)
	}
}

func (d *decodeState) boolInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := d.data[d.off]
	d.off += 1

	if val == 0 {
		return false
	} else {
		return true
	}
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(4)
func (d *decodeState) float(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Float32(d.data[d.off:])
	d.off += 4

	v.SetFloat(float64(val))
}

func (d *decodeState) floatInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Float32(d.data[d.off:])
	d.off += 4

	return val
}

// type(1) | name length(1) | raw name bytes | 0x00 | value bytes(8)
func (d *decodeState) double(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Float64(d.data[d.off:])
	d.off += 8

	v.SetFloat(val)
}

func (d *decodeState) doubleInterface() interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	d.off += klen

	val := Float64(d.data[d.off:])
	d.off += 8

	return val
}

func (d *decodeState) valueInterface() interface{} {
	switch d.data[d.off] {
	case KSPACK_OBJECT:
		return d.objectInterface()
	case KSPACK_ARRAY:
		return d.arrayInterface()
	case KSPACK_STRING:
		return d.stringInterface()
	case KSPACK_SHORT_STRING:
		return d.shortStringInterface()
	case KSPACK_BINARY:
		return d.binaryInterface()
	case KSPACK_SHORT_BINARY:
		return d.shortBinaryInterface()
	case KSPACK_INT8:
		return d.int8Interface()
	case KSPACK_INT16:
		return d.int16Interface()
	case KSPACK_INT32:
		return d.int32Interface()
	case KSPACK_INT64:
		return d.int64Interface()
	case KSPACK_UINT8:
		return d.uint8Interface()
	case KSPACK_UINT16:
		return d.uint16Interface()
	case KSPACK_UINT32:
		return d.uint32Interface()
	case KSPACK_UINT64:
		return d.uint64Interface()
	case KSPACK_BOOL:
		return d.boolInterface()
	case KSPACK_FLOAT:
		return d.floatInterface()
	case KSPACK_DOUBLE:
		return d.doubleInterface()
	case KSPACK_NULL:
		return d.nullInterface()
	}
	return nil
}

// type(1) | name length(1) | item size(4) | raw name bytes | 0x00
// | members number(4) | member1 | ... | memberN
func (d *decodeState) object(v reflect.Value) {
	if v.Kind() == reflect.Interface && v.NumMethod() == 0 {
		v.Set(reflect.ValueOf(d.objectInterface()))
		return
	}

	// make map
	if v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	// vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	n := int(Uint32(d.data[d.off:]))
	d.off += 4 // member number

	var mapElem reflect.Value
	for i := 0; i < n; i++ {
		subk := d.key()
		var subv reflect.Value

		if v.Kind() == reflect.Map {
			elemType := v.Type().Elem()
			if !mapElem.IsValid() {
				mapElem = reflect.New(elemType).Elem()
			} else {
				mapElem.Set(reflect.Zero(elemType))
			}
			subv = mapElem
		} else {
			var f *field
			fields := cachedTypeFields(v.Type())
			for i := range fields {
				ff := &fields[i]
				if bytes.Equal(ff.nameBytes, subk) {
					f = ff
					break
				}
				if f == nil && ff.equalFold(ff.nameBytes, subk) {
					f = ff
				}
			}
			if f != nil {
				subv = v
				for _, i := range f.index {
					if subv.Kind() == reflect.Ptr {
						if subv.IsNil() {
							subv.Set(reflect.New(subv.Type()).Elem())
						}
						subv = subv.Elem()
					}
					subv = subv.Field(i)
				}
			}
		}

		d.value(subv)

		// Write value back to map
		if v.Kind() == reflect.Map {
			kv := reflect.ValueOf(subk).Convert(v.Type().Key())
			v.SetMapIndex(kv, subv)
		}
	}
}

func (d *decodeState) objectInterface() map[string]interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	// vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 // content length

	d.off += klen // name and 0x00

	n := int(Uint32(d.data[d.off:]))
	d.off += 4 // member number

	m := make(map[string]interface{})
	for i := 0; i < n; i++ {
		subk := d.key()
		m[string(subk)] = d.valueInterface()
	}

	return m
}

//FIXME: fix when v is invalid
// type(1) | name length(1) | item size(4) | raw name bytes | 0x00
// | element number(4) | element1 | ... | elementN
func (d *decodeState) array(v reflect.Value) {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	// vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 //  content length

	//var key string
	d.off += klen

	n := int(Uint32(d.data[d.off:]))
	d.off += 4 // member number

	if v.Kind() == reflect.Slice {
		if n > v.Cap() {
			newv := reflect.MakeSlice(v.Type(), n, n)
			v.Set(newv)
		}
		v.SetLen(n)
	}

	for i := 0; i < n; i++ {
		if i < v.Len() {
			d.value(v.Index(i))
		} else {
			d.value(reflect.Value{})
		}
	}

	if n < v.Len() {
		if v.Kind() == reflect.Array {
			z := reflect.Zero(v.Type().Elem())
			for i := n; i < v.Len(); i++ {
				v.Index(i).Set(z)
			}
		}
	}

	if n == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
}

func (d *decodeState) arrayInterface() []interface{} {
	d.off += 1 // type

	klen := int(Uint8(d.data[d.off:]))
	d.off += 1 // name length

	// vlen := int(Uint32(d.data[d.off:]))
	d.off += 4 //  content length

	//var key string
	d.off += klen

	n := int(Uint32(d.data[d.off:]))
	d.off += 4 // member number

	v := make([]interface{}, n)
	for i := 0; i < n; i++ {
		v[i] = d.valueInterface()
	}
	return v
}

func (d *decodeState) key() []byte {
	var kstart int
	switch d.data[d.off] {
	case KSPACK_INT8, KSPACK_INT16, KSPACK_INT32, KSPACK_INT64,
		KSPACK_UINT8, KSPACK_UINT16, KSPACK_UINT32, KSPACK_UINT64,
		KSPACK_BOOL, KSPACK_FLOAT, KSPACK_DOUBLE, KSPACK_NULL:
		kstart = 2 // type + klen
	case KSPACK_SHORT_BINARY, KSPACK_SHORT_STRING:
		kstart = 3 // type + klen + vlen(1)
	case KSPACK_BINARY, KSPACK_STRING, KSPACK_OBJECT, KSPACK_ARRAY:
		kstart = 6 // type + klen + vlen(4)
	}
	klen := int(Uint8(d.data[d.off+1:]))
	if klen <= 0 {
		d.error(errEmptyKey)
	}
	return d.data[d.off+kstart : d.off+kstart+klen-1]
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "kspack: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "kspack: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "kspack: Unmarshal(nil " + e.Type.String() + ")"
}
