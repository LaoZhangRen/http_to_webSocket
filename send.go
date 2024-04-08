package main

import "net/http"

func send(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//接收get patient_id参数
	inpatientNo := r.FormValue("inpatient_no")
	//接收get message参数
	message := r.FormValue("message")
	// 接收get code参数
	code := r.FormValue("code")
	//根据patient_id找到对应的client
	client := getClient(inpatientNo, hub)

	//如果client存在，则将message发送给client
	if client != nil {
		client.send <- body(code, message)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("send success"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("send fail"))
	}
}

func body(code string, message string) []byte {
	return []byte(`{"status":"success","message":"` + message + `","code":` + code + `}`)
}
