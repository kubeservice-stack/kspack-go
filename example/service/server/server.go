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

package main

import (
	"net/http"
	"time"

	codec "github.com/kubeservice-stack/kspack-go"
)

type ServerData struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a, err := codec.PluginInstance(codec.KSPACK).Marshal(&ServerData{
			Name:     "dongjiang",
			BirthDay: time.Date(2017, 7, 7, 9, 0, 0, 0, time.Local),
			Phone:    "13811111111",
			Siblings: 10,
			Spouse:   true,
			Money:    29.11,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(a)
	})

	s := &http.Server{
		Addr:           "0.0.0.0:18080",
		Handler:        mux,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20, // 1048576; 1MiB
	}
	s.ListenAndServe()
}
