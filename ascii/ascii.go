package ascii

import (
	"encoding/json"
	"io/ioutil"

	"github.com/blacktop/scifgif/database"
)

// GetAllASCIIEmoji loads all ascii-emojis into database
func GetAllASCIIEmoji() error {

	// open database
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	file, err := ioutil.ReadFile("ascii/emoji.json")
	if err != nil {
		return err
	}

	emojis := make([]database.ASCIIData, 0)
	err = json.Unmarshal(file, &emojis)
	if err != nil {
		return err
	}

	for _, e := range emojis {
		// index into database
		db.WriteASCIIToDatabase(database.ASCIIData{
			ID:       e.ID,
			Source:   "ascii",
			Keywords: e.Keywords,
			Emoji:    e.Emoji,
		})
	}

	return nil
}
