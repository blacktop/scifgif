package database

// UpdateKeywords adds new keywords to an image's search text
func (db *Database) UpdateKeywords(image ImageMetaData) error {

	// remove from bleve
	db.IDX.Delete(image.ID)
	// add back to bleve
	err := db.IDX.Index(image.ID, image)
	if err != nil {
		return err
	}
	// update sqlite
	db.SQL.Save(&image)

	return nil
}
