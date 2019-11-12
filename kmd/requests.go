package kmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"io/ioutil"
	"net/http"
)

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
	var respData interface{}
	respData = &jsonSuccessResponse{}
	if err := json.Unmarshal(respTx, respData); err != nil {
		fmt.Println("err01: ", err)
		respData = &jsonErrResponse{}
		err := json.Unmarshal(respTx, respData)
		if err != nil {
			return nil, err
		}
	}
	return respData, nil
}
