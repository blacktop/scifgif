package database

import "github.com/jinzhu/gorm"

// GetRandomImage returns a random image path from source (xkcd/giphy)
func (db *Database) GetRandomImage(source string) (ImageMetaData, error) {

	var image ImageMetaData

	if db.SQL.Order(gorm.Expr("random()")).Where("source = ?", source).First(&image).RecordNotFound() {
		return ImageMetaData{}, ErrNoImagesFound
	}

	return image, nil
}

// GetRandomASCII returns a random ascii-emoji
func (db *Database) GetRandomASCII() (ASCIIData, error) {

	var ascii ASCIIData

	if db.SQL.Order(gorm.Expr("random()")).Where("source = ?", "ascii").First(&ascii).RecordNotFound() {
		return ASCIIData{}, ErrNoASCIIFound
	}

	return ascii, nil
}
