package pkg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 150 * time.Millisecond,
}

type SuccessResponse struct {
	Result   *http.Response `json:"result"`
	Endpoint string         `json:"endpoint"`
	RTT      int            `json:"rtt"`
	Body     []byte         `json:"-"`
}

func Shotgun(endpoints []string, mainRequest *http.Request) (SuccessResponse, error) {
	resultCh := make(chan SuccessResponse, len(endpoints))
	reqBody, err := io.ReadAll(mainRequest.Body)
	if err != nil {
		return <-resultCh, err
	}
	fmt.Printf("Shotgun %v\n", string(reqBody))
	for _, endpoint := range endpoints {
		go makeRequest(endpoint, mainRequest, reqBody, resultCh)
	}
	return <-resultCh, nil
}

func makeRequest(endpoint string, req *http.Request, reqBodyBytes []byte, successCh chan SuccessResponse) {
	// Create a new request with the same method and body as the original request
	newReq, err := http.NewRequest(req.Method, endpoint, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return
	}
	// Copy headers from the original request to the new request
	newReq.Header = make(http.Header)
	for key, values := range req.Header {
		newReq.Header[key] = values
	}
	//newReq.Header["Content-Type"] = []string{"application/json"}
	// Send request
	startTime := time.Now()
	resp, err := client.Do(newReq)
	endTime := time.Now()

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		successCh <- SuccessResponse{
			Result:   resp,
			Endpoint: endpoint,
			RTT:      int(endTime.Sub(startTime).Milliseconds()),
			Body:     body,
		}
	}
}
