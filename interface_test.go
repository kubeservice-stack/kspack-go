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

package codec_test

import (
	"bytes"
	"reflect"
	"testing"

	codec "github.com/kubeservice-stack/kspack-go"
	"github.com/stretchr/testify/assert"
)

func BenchmarkKSPack(b *testing.B) {
	va := new(int)
	*va = 1
	str := new(string)
	*str = "dongjiang"
	str1 := new(string)
	*str1 = "long string users/dongjiang/Documentsdfasdntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/kspack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/kspack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/kspack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/kspack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/kspack]asdfasntsdfasdfsdgs/go/src/github"

	a := &TV{
		F: map[string]interface{}{
			"ui64":  uint64(0xFFFFFFFFFFFFFFFF),
			"ui32":  uint32(0xFFFFFFFF),
			"bys":   bytes.Runes([]byte("dasdf")),
			"alpha": "a-z",
			"a":     1,
			"is":    true,
			"str":   str,
			"tttt":  str1,
			"dd":    float64(1.11),
			"ff":    float32(1.11),
			"b":     va,
			"c":     reflect.ValueOf(va),
			"d":     map[string]interface{}{"aa": "bb"},
			"e":     []interface{}{"aa", 1, va},
			"f":     map[string]float64{"a": float64(-45.2231)},
		},
	}
	for i := 0; i < b.N; i++ {
		codec.PluginInstance(codec.KSPACK).Marshal(a)
	}
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	aa := codec.HasRegister(codec.PACK("aaa"))
	assert.False(aa)
	aa = codec.HasRegister(codec.KSPACK)
	assert.True(aa)

	codec.Register("dfdf", codec.NewKSPack)
	aa = codec.HasRegister(codec.PACK("dfdf"))
	assert.True(aa)
	adp := codec.PluginInstance(codec.PACK("dfdf"))
	assert.NotNil(adp)

	adp = codec.PluginInstance(codec.PACK("gh"))
	assert.Nil(adp)
}

func TestRegisterPanic(t *testing.T) {
	assert := assert.New(t)
	aa := codec.HasRegister(codec.PACK("aaa"))
	assert.False(aa)
	assert.Panics(func() { codec.Register(codec.PACK("bb"), nil) })
	codec.Register("aa", codec.NewKSPack)
	assert.Panics(func() { codec.Register(codec.PACK("aa"), codec.NewKSPack) })

}
