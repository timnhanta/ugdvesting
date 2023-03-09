package types

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	HedgehogBaseUrlTest = "https://localhost:52884/gridspork/vesting-storage/"
	HedgehogBaseUrl     = "https://daemon:52884/gridspork/vesting-storage/"
)

func HegdehogRequestGetVestingByAddr(addr string) Vesting {
	// Create a new HTTP client with custom TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send the HTTP request and get the response
	// "/gridspork/vesting-storage"
	resp, err := client.Get(HedgehogBaseUrlTest + addr)
	if err != nil {
		panic(err)
	}

	// Read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))

	var vesting Vesting
	err = json.Unmarshal([]byte(body), &vesting)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	fmt.Printf("%+v\n", vesting)

	return vesting
}
