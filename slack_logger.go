package goweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type slackLogger struct {
	minLevel int
}

func (l *slackLogger) Log(level int, message interface{}) {
	if level >= l.minLevel {
		title, color := getTitleAndColor(level)
		req, err := makeSlackRequest(Map{
			"color": color,
			"title": title,
			"text":  fmt.Sprintf("%s - %s", caller(), message),
		})
		if err != nil {
			log.Println(err)
			return
		}
		go sendSlackRequest(req)
	}
}

func makeSlackRequest(m Map) (*http.Request, error) {
	payload, err := json.Marshal(Map{
		"attachments": []Map{m},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", os.Getenv("SLACK_LOG_WEBHOOK"), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func sendSlackRequest(req *http.Request) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b := new(bytes.Buffer)
		b.ReadFrom(res.Body)
		log.Printf(
			"failed to send slack message: %d - %s\n",
			res.StatusCode,
			b.String(),
		)
	}
}
