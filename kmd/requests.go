package kmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"io/ioutil"
	"log"
	"net/http"
)

// 发送交易上链
func SendTx(tx *types.Transaction, url string) (interface{}, error) {
	jsonTx, _ := json.Marshal(tx)
	reqData := &jsonRequest{
		Method:  "tx_sendTx",
		Version: "2.0",
		Id:      1,
		Payload: []json.RawMessage{jsonTx},
	}
	jsonReqData, _ := json.Marshal(reqData)
	reader := bytes.NewReader(jsonReqData)
	// post
	resp, err := http.Post(url, "application/json;charset=UTF-8", reader)
	if err != nil {
		fmt.Println("发送交易post出错：", err)
		return nil, err
	}
	respTx, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respon := new(jsonSuccessResponse)
	err = json.Unmarshal(respTx, respon)
	if err != nil {
		log.Println("unmarshal error:", err)
		return "", nil
	}
	return respon.Result.(string), nil
}

// getBalance
func getBalance(lemoAddress string, url string) (string, error) {
	jsonlemoAdd, err := json.Marshal(lemoAddress)
	if err != nil {
		log.Println("json 102 marshal error:", err)
	}
	data := &jsonRequest{
		Version: "2.0",
		Id:      1,
		Method:  "account_getBalance",
		Payload: []json.RawMessage{jsonlemoAdd},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("json 112 marshal error:", err)
		return "", err
	}
	reader := bytes.NewReader(jsonData)
	resp, err := http.Post(url, "application/json;charset=UTF-8", reader)
	if err != nil {
		log.Println("post error:", err)
		return "", err
	}
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response error:", err)
		return "", err
	}
	defer resp.Body.Close()
	respon := new(jsonSuccessResponse)
	err = json.Unmarshal(byteResp, respon)
	if err != nil {
		log.Println("unmarshal error:", err)
		return "", nil
	}
	return respon.Result.(string), nil
}

func getAccount(lemoAddress string, url string) (map[string]interface{}, error) {
	jsonlemoAdd, err := json.Marshal(lemoAddress)
	if err != nil {
		log.Println("json 102 marshal error:", err)
	}
	data := &jsonRequest{
		Version: "2.0",
		Id:      1,
		Method:  "account_getAccount",
		Payload: []json.RawMessage{jsonlemoAdd},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("json 112 marshal error:", err)
		return nil, err
	}
	reader := bytes.NewReader(jsonData)
	resp, err := http.Post(url, "application/json;charset=UTF-8", reader)
	if err != nil {
		log.Println("post error:", err)
		return nil, err
	}
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	respon := new(jsonSuccessResponse)
	err = json.Unmarshal(byteResp, respon)
	if err != nil {
		log.Println("unmarshal error:", err)
		return nil, nil
	}

	return respon.Result.(map[string]interface{}), nil
}
