package database

// GetImageByID returns the path to an image by id
func (db *Database) GetImageByID(id string) (ImageMetaData, error) {
	var image ImageMetaData

	if db.SQL.Find(&image, id).RecordNotFound() {
		return ImageMetaData{}, ErrNoImagesFound
	}

	return image, nil
}

// GetImageByPath returns an image by path
func (db *Database) GetImageByPath(path string) (ImageMetaData, error) {
	var image ImageMetaData

	if db.SQL.Where("path = ?", path).First(&image).RecordNotFound() {
		return ImageMetaData{}, ErrNoImagesFound
	}

	return image, nil
}
