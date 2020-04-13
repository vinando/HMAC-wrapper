package wrapper

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HMACWrapper struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
}

func Init(clientID string, clientSecret string, baseURL string) *HMACWrapper {
	wrapper := &HMACWrapper{}
	wrapper.ClientID = clientID
	wrapper.ClientSecret = clientSecret
	wrapper.BaseURL = baseURL
	return wrapper
}

var client = &http.Client{}

// param endpoint must start with a '/'
func (wp *HMACWrapper) DoGet(endpoint string, headers map[string]string) (resp interface{}, err error) {
	signature, t := wp.constructSignature("GET " + endpoint + " HTTP/1.1")
	req, err := http.NewRequest("GET", wp.BaseURL+endpoint, nil)
	if err != nil {
		fmt.Println("Error initiating request")
		return nil, err
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	req.Header.Set("Date", t)
	req.Header.Set("Authorization", fmt.Sprintf("hmac username=\"%s\", algorithm=\"hmac-sha256\", headers=\"date request-line\", signature=\"%s\"", wp.ClientID, signature))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		fmt.Println("Error decode json body")
		return nil, err
	}
	return resp, nil
}

func (wp *HMACWrapper) DoPost(endpoint string, body []byte, headers map[string]string) (resp interface{}, err error) {
	signature, t := wp.constructSignature("POST " + endpoint + " HTTP/1.1")
	digest, err := wp.constructDigest(body)
	req, err := http.NewRequest("POST", wp.BaseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error initiating request")
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Date", t)
	req.Header.Set("Digest", "SHA-256="+digest)
	req.Header.Set("Authorization", fmt.Sprintf("hmac username=\"%s\", algorithm=\"hmac-sha256\", headers=\"date request-line\", signature=\"%s\"", wp.ClientID, signature))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request")
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Fatalf("Error decode json body")
	}
	return resp, nil
}

func (wp *HMACWrapper) constructSignature(reqLine string) (signature string, t string) {
	t = time.Now().UTC().Format(time.RFC1123)
	str := "date: " + t + "\n" + reqLine
	h := hmac.New(sha256.New, []byte(wp.ClientSecret))
	h.Write([]byte(str))

	signature = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return
}

func (wp *HMACWrapper) constructDigest(body []byte) (string, error) {
	hh := sha256.New()
	hh.Write(body)
	digest := base64.StdEncoding.EncodeToString(hh.Sum(nil))
	return digest, nil
}
