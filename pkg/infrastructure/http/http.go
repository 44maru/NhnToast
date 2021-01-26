package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const Token = "9d36f54badf54162a51e678dd52c5e29"

func Delete(url string, header map[string]string) error {
	req, err := createRequest("DELETE", url, header, nil, nil)
	if err != nil {
		log.Printf("ERROR : failed create http request. %s\n", err.Error())
		return err
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		log.Println("client.Do() error")
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || 300 <= res.StatusCode {
		log.Printf("http status error. status %d\n", res.StatusCode)
		return fmt.Errorf("http status error. status %d", res.StatusCode)
	}

	return nil
}

func Get(url string, header, queryParam map[string]string) ([]byte, error) {
	req, err := createRequest("GET", url, header, queryParam, nil)
	if err != nil {
		log.Printf("ERROR : failed create http request. %s\n", err.Error())
		return nil, err
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		log.Println("client.Do() error")
		return nil, err
	}
	defer res.Body.Close()

	return analyzeHttpResponse(res)
}

func Post(url string, jsonBytes []byte, header map[string]string) ([]byte, error) {
	return postOrPut("POST", url, jsonBytes, header)
}

func Put(url string, jsonBytes []byte, header map[string]string) ([]byte, error) {
	return postOrPut("PUT", url, jsonBytes, header)
}

func createRequest(httpMethod, url string, header, queryParam map[string]string, bodyBuffer *bytes.Buffer) (*http.Request, error) {
	var req *http.Request
	var err error
	if bodyBuffer == nil {
		req, err = http.NewRequest(httpMethod, url, nil)
	} else {
		req, err = http.NewRequest(httpMethod, url, bodyBuffer)
	}

	if err != nil {
		return nil, err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	params := req.URL.Query()
	for key, value := range queryParam {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()

	return req, nil
}

func postOrPut(postOrPut, url string, jsonBytes []byte, header map[string]string) ([]byte, error) {
	//log.Printf("HTTP %s to %s\n", postOrPut, url)
	//log.Printf("HTTP Body is below.\n%s\n", string(jsonBytes))

	req, err := createRequest(postOrPut, url, header, nil, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("ERROR : failed create http request. %s\n", err.Error())
		return nil, err
	}

	res, err := new(http.Client).Do(req)
	if err != nil {
		log.Println("client.Do() error")
		return nil, err
	}
	defer res.Body.Close()

	return analyzeHttpResponse(res)
}

func analyzeHttpResponse(res *http.Response) ([]byte, error) {

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Response read error.")
		return bodyBytes, err
	}

	if res.StatusCode < 200 || 300 <= res.StatusCode {
		log.Printf("http status error. status %d\n", res.StatusCode)
		return bodyBytes, fmt.Errorf("http status error. status %d. resposne body is below.\n%s\n", res.StatusCode, string(bodyBytes))
	}

	var buf bytes.Buffer
	if len(bodyBytes) > 0 {
		err = json.Indent(&buf, bodyBytes, "", "  ")
		if err != nil {
			log.Println("Response JSON format error.")
			return bodyBytes, err
		}
	}

	return bodyBytes, nil
}
