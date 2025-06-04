package llmServer

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"io"
)

type RpcRequest struct {
    Data map[string]interface{} `json:"Data"`
}

type RpcResponse struct {
    Result interface{} `json:"Result"`
    Status int `json:"Status"`
}

func CallRPC(method string, data map[string]interface{}) (interface{}, error) {
	if data == nil {
		data = make(map[string]interface{})
	}
    reqBody, _ := json.Marshal(RpcRequest{Data: data})
    resp, err := http.Post(fmt.Sprintf("http://localhost:5974/%s", method), "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

	check, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Raw response: %s\n", string(check))

    var respParsed RpcResponse
	// fmt.Println(json.NewDecoder(resp.Body).Decode(&respParsed))
    if err := json.Unmarshal(check, &respParsed); err != nil {
        return nil, err
    }
    return respParsed.Result, nil
}