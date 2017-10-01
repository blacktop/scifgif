package ascii

import (
	"encoding/json"
	"io/ioutil"

	"github.com/blacktop/scifgif/elasticsearch"
)

// GetAllASCIIEmoji loads all ascii-emojis into elasticsearch
func GetAllASCIIEmoji() error {
	file, err := ioutil.ReadFile("ascii/emoji.json")
	if err != nil {
		return err
	}

	emojis := make([]elasticsearch.ASCIIData, 0)
	err = json.Unmarshal(file, &emojis)
	if err != nil {
		return err
	}

	for _, e := range emojis {
		// index into elasticsearch
		elasticsearch.WriteASCIIToDatabase(elasticsearch.ASCIIData{
			ID:       e.ID,
			Source:   "ascii",
			Keywords: e.Keywords,
			Emoji:    e.Emoji,
		})
	}
	return nil
}
