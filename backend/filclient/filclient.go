package filclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type FilClient struct {
	filecoinApi string
}

func NewFileClient(filecoinApi string) *FilClient {
	return &FilClient{
		filecoinApi: filecoinApi,
	}
}

func (fc *FilClient) rpcCall(ctx context.Context, rpcMetchod string, rpcParams ...interface{}) (result json.RawMessage, err error) {
	logger.Debug(fc.filecoinApi, ", method=", rpcMetchod, ", rpcParams=", rpcParams)
	rpcId := rand.Int()

	rpcJsonMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  rpcMetchod,
		"params":  rpcParams,
		"id":      rpcId,
	}

	mJson, err := json.Marshal(rpcJsonMsg)
	if err != nil {
		logger.Errorf("RpcCall marshal param failed: %v", err)
		return
	}
	contentReader := bytes.NewReader(mJson)
	request, err := http.NewRequestWithContext(ctx, "POST", fc.filecoinApi, contentReader)
	if err != nil {
		logger.Errorf("RpcCall http request failed: %v", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Errorf("RpcCall client do failed: %v", err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("RpcCall ioutil failed: %v", err)
		return
	}

	var jsonResp clientResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		logger.WithField("body", string(body)).Errorf("filclient unmarshal json failed: %v", err)
		return
	}

	if jsonResp.Error != nil {
		return nil, errors.New(jsonResp.Error.Error())
	}

	return jsonResp.Result, nil
}

type clientResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	ID      int64           `json:"id"`
	Error   *respError      `json:"error,omitempty"`
}

type respError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *respError) Error() string {
	if e.Code >= -32768 && e.Code <= -32000 {
		return fmt.Sprintf("RPC error (%d): %s", e.Code, e.Message)
	}
	return e.Message
}
