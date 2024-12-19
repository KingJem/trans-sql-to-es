package main

import (
	"encoding/json"
	"fmt"
	"github.com/cch123/elasticsql"
	"net/http"
)

// TransRes 存储转换结果
type TransRes struct {
	Sql     string `json:"sql"`
	Es      string `json:"es"`
	Message string `json:"message"`
}

// TransReq 存储转换请求
type TransReq struct {
	Sql string `json:"sql"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var transreq TransReq
	var transres TransRes
	// 从请求体中解码 JSON 数据
	err := json.NewDecoder(r.Body).Decode(&transreq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 将 SQL 转换为 Elasticsearch DSL
	dsl, _, err := elasticsql.Convert(transreq.Sql)
	if err != nil {
		transres.Message = err.Error()
		transres.Sql = transreq.Sql
		transres.Es = ""
	} else {
		transres.Es = dsl
		transres.Sql = transreq.Sql
		transres.Message = "success"
	}

	// 将结果转换为 JSON 数据
	jsonData, err := json.Marshal(transres)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 发送 HTTP 响应
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 打印 JSON 数据
	fmt.Println(string(jsonData))
}

func main() {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":9898", nil)
}
