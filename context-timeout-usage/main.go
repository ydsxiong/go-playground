package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type CallResponse struct {
	Resp *Response
	Err  error
}

type Response struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var mockRemoteAPIServiceEndpoint = "https://jsonplaceholder.typicode.com/todos/1"

func callAPIService(logger log.Logger) (apiResponse *Response, err error) {
	defer func(begin time.Time) {
		_ = logger.Log(
			"method", "callAPIService",
			"uid", apiResponse.UserID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	apiResponse = &Response{}

	resp, err := http.Get(mockRemoteAPIServiceEndpoint)
	if err != nil {
		err = fmt.Errorf("unexpected error from http call")
		return
	}

	defer resp.Body.Close()

	byedata, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("Error in reading response")
		return
	}

	err = json.Unmarshal(byedata, apiResponse)

	if err != nil {
		err = fmt.Errorf("error in unmarshalling response")
	}

	return
}

func wrapUpAPICall(logger log.Logger) <-chan *CallResponse {
	respChan := make(chan *CallResponse, 1)
	go func() {
		res, err := callAPIService(logger)
		respChan <- &CallResponse{res, err}
	}()
	return respChan
}

func callForHttpResponseWithinTimeLimit(logger log.Logger, timeLimitMS int) (res *Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMS)*time.Millisecond)
	defer cancel()

	select { // whichever case returning response quicker:
	case <-ctx.Done():
		return nil, fmt.Errorf("Time's up!!!")
	case respChan := <-wrapUpAPICall(logger):
		return respChan.Resp, respChan.Err
	}
}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	res, err := callForHttpResponseWithinTimeLimit(logger, 250)

	if err != nil {
		fmt.Printf("error received: %v", err)
	} else {
		fmt.Printf("response received: %v", res)
	}
}
