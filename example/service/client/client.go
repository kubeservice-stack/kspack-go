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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	codec "github.com/kubeservice-stack/kspack-go"
)

type ClientData struct {
	Siblings int
	Name     string
	Phone    string
	Money    float64
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		port := r.URL.Query().Get("port")

		resp, err := http.Get("http://" + ip + ":" + port + "/")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		var c ClientData
		err = codec.PluginInstance(codec.KSPACK).Unmarshal(body, &c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonResp, _ := json.Marshal(c)
			w.Write(jsonResp)
			return
		}
	})

	s := &http.Server{
		Addr:           "0.0.0.0:18081",
		Handler:        mux,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20, // 1048576; 1MiB
	}
	s.ListenAndServe()
}
