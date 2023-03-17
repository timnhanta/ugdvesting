package types

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	HedgehogBaseUrlTest = "https://localhost:52884/gridspork/vesting-storage/"
)

func HegdehogRequestGetVestingByAddr(addr string) Vesting {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(HedgehogBaseUrlTest + addr)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response Vesting
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		panic(err)
	}

	return response
}
