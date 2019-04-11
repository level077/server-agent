package elastic_send

import (
	//"fmt"
	"net/http"
	"bytes"
)

func Bulk(body string,esHost string, esPort string) (string) {
	url := "http://" + esHost + ":" + esPort + "/_bulk"
	jsonStr := []byte(body)	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    	req.Header.Set("Content-Type", "application/x-ndjson")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
        	panic(err)
    	}
    	defer resp.Body.Close()

	return resp.Status
	/*
    	fmt.Println("response Status:", resp.Status)
    	fmt.Println("response Headers:", resp.Header)
    	body, _ := ioutil.ReadAll(resp.Body)
    	fmt.Println("response Body:", string(body))
	*/
}
