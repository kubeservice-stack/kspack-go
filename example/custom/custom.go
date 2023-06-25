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

package example

import (
	"fmt"

	codec "github.com/kubeservice-stack/kspack-go"
)

type CustomStruct struct {
	S string
	N int
}

func main() {
	a, err := codec.PluginInstance(codec.KSPACK).Marshal(&CustomStruct{S: "Hello", N: 100})
	if err != nil {
		panic(err)
	}

	var v CustomStruct
	err = codec.PluginInstance(codec.KSPACK).Unmarshal(a, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", v)
}
