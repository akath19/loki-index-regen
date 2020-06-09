package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/yellowpepper/loki-index-regen/inspect"
)

type chunksParser struct {
	LokiAddress string
}

type logEntry struct {
	Timestamp time.Time `json:"ts"`
	Line      string    `json:"line"`
}

type promtailStream struct {
	Labels  string      `json:"labels"`
	Entries []*logEntry `json:"entries"`
}

type promtailMessage struct {
	Streams []promtailStream `json:"streams"`
}

func (c *chunksParser) ParseAndSaveChunk(path string) {
	header, chunk, err := inspect.ParseFile(path, true, true, false)

	if err != nil {
		log.Fatal(err)
	}

	err = c.sendToLoki(header, chunk)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(path)

	if err != nil {
		log.Fatal(err)
	}
}
func (c *chunksParser) sendToLoki(headers *inspect.ChunkHeader, chunk *inspect.LokiChunk) error {
	var streams []promtailStream

	streams = append(streams, promtailStream{
		Labels:  headers.Metric.String(),
		Entries: getEntries(chunk),
	})

	message := promtailMessage{
		Streams: streams,
	}

	jsonMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	res, httpErr := http.Post(c.LokiAddress+"/api/prom/push", "application/json", bytes.NewReader(jsonMessage))

	if httpErr != nil {
		return err
	}

	if res.StatusCode != 204 {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		return fmt.Errorf("unexpected HTTP status code: %d, message: %s", res.StatusCode, string(body))
	} else {
		log.Println("Chunk sent to Loki successfully")
	}

	return nil
}

func getEntries(chunk *inspect.LokiChunk) []*logEntry {
	var entries []*logEntry

	for _, block := range chunk.Blocks {
		for _, entry := range block.Entries {
			entries = append(entries, &logEntry{
				Timestamp: time.Unix(entry.Timestamp/1000000000, 0),
				Line:      entry.Line,
			})
		}
	}

	return entries
}
