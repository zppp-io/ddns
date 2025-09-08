package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"zppp.io/ddns/config"

	"go.uber.org/zap"
)

type Payload struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	config.InitEnvConfig()

	apiToken := config.GetConfig("API_TOKEN")
	zoneId := config.GetConfig("ZONE_ID")
	recordId := config.GetConfig("RECORD_ID")
	interval := config.GetConfig("INTERVAL")
	name := config.GetConfig("NAME")
	pro := config.GetConfig("PROXIED")
	intervalTime, _ := strconv.Atoi(interval)
	if intervalTime <= 0 {
		intervalTime = 1
	}
	proxied := false
	if pro == "true" {
		proxied = true
	}

	currentIp := "0.0.0.0"

	for {
		rsp, err := http.Get("http://myip.ipip.net")
		if err != nil {
			sugar.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			sugar.Error(err)
			continue
		}
		data := string(body)
		reg := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+.[0-9]+")
		ip := reg.FindString(data)
		if ip == "" {
			sugar.Error("cannot match ip")
			continue
		}
		sugar.Info("current ip is ", ip)
		if currentIp != ip {
			currentIp = ip
			payload := &Payload{
				Type:    "A",
				Name:    name,
				Content: ip,
				Ttl:     1,
				Proxied: proxied,
			}
			data, err := json.Marshal(payload)
			if err != nil {
				sugar.Fatal(err)
			}

			req, _ := http.NewRequest(http.MethodPut, "https://api.cloudflare.com/client/v4/zones/"+zoneId+"/dns_records/"+recordId, bytes.NewBuffer(data))
			req.Header.Add("Authorization", "Bearer "+apiToken)
			req.Header.Add("Content-Type", "application/json")
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				sugar.Fatal(err)
				time.Sleep(5 * time.Second)
				continue
			}
			body, _ = ioutil.ReadAll(response.Body)
			if response.StatusCode != 200 {
				sugar.Fatal("failed sync ip to cloud flare", zap.String("error", string(body)))
			} else {
				sugar.Info("sync success")
			}
		}

		time.Sleep(time.Duration(intervalTime) * time.Minute)
	}

}
