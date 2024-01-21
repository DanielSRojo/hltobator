package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Game struct {
	Name string `json:"game_name"`
	ID   uint   `json:"game_id"`
	Main uint   `json:"comp_main"`
}

type Response struct {
	Count uint   `json:"count"`
	Data  []Game `json:"data"`
}

func (g *Game) Duration() uint {
	return g.Main / 3600
}

func main() {
	url := "https://howlongtobeat.com/api/search"
	gameName := os.Args[1]
	payload := []byte(fmt.Sprintf(`{"searchType":"games","searchTerms":["%s"],"searchPage":1,"size":20,"searchOptions":{"games":{"userId":0,"platform":"","sortCategory":"popular","rangeCategory":"main","rangeTime":{"min":null,"max":null},"gameplay":{"perspective":"","flow":"","genre":""},"rangeYear":{"min":"","max":""},"modifier":""},"users":{"sortCategory":"postcount"},"lists":{"sortCategory":"follows"},"filter":"","sort":0,"randomizer":0}}`, gameName))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Host", "howlongtobeat.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "es")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://howlongtobeat.com/")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprint(len(payload)))
	req.Header.Set("Origin", "https://howlongtobeat.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("TE", "trailers")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	var body []byte

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading decompressed response:", err)
			return
		}
	default:
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
	}

	//	fmt.Println(string(body))

	r := &Response{}
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, g := range r.Data {
		fmt.Printf("%s: %dh\n", g.Name, g.Duration())
	}

}
