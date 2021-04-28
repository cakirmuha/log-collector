package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log_collector "github.com/cakirmuha/log-collector"
)

var client = http.Client{Timeout: 30 * time.Second}

const apiLogURL = "http://localhost:8484/log"

func importDummyEvents() {
	file, _ := ioutil.ReadFile("files/dummyEvents.json")

	var data []log_collector.ApiResponse
	{
		_ = json.Unmarshal(file, &data)

		fmt.Println(len(data))

		for i, d := range data {
			fmt.Println(i)
			b, _ := json.Marshal(d)

			payload := strings.NewReader(string(b))
			req, _ := http.NewRequest("POST", apiLogURL, payload)

			response, _ := client.Do(req)
			respBody, _ := ioutil.ReadAll(response.Body)
			var res log_collector.ApiResponse
			_ = json.Unmarshal(respBody, &res)

			defer response.Body.Close()
			if !res.Success {
				fmt.Println("Error:", res.Error, d)
			}
		}
	}
}
