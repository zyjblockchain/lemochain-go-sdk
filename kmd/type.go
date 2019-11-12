package kmd

import "encoding/json"

type jsonRequest struct {
	Method  string            `json:"method"`
	Version string            `json:"jsonrpc"`
	Id      uint64            `json:"id"`
	Payload []json.RawMessage `json:"params,omitempty"`
}

type jsonSuccessResponse struct {
	Version string      `json:"jsonrpc"`
	Id      uint64      `json:"id"`
	Result  interface{} `json:"result"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type jsonErrResponse struct {
	Version string    `json:"jsonrpc"`
	Id      uint64    `json:"id"`
	Error   jsonError `json:"error"`
}
