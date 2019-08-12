package database

import (
	"errors"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// log "github.com/sirupsen/logrus"
)

const (
	bleveIndex = "scifgif.bleve"
	sqliteFile = "scifgif.db"
)

var (
	// ErrNoImagesFound represents a failure to find images.
	ErrNoImagesFound = errors.New("no images found")
	// ErrNoASCIIFound represents a failure to find ascii.
	ErrNoASCIIFound = errors.New("no ascii found")
)

// Database database object
type Database struct {
	SQL *gorm.DB
	IDX bleve.Index
}

// ImageMetaData image meta-data object
type ImageMetaData struct {
	ID string `json:"id,omitempty" gorm:"primary_key"`
	// ID     string `json:"id,omitempty" gorm:"primary_key"`
	Source string `json:"source,omitempty"`
	Name   string `json:"name,omitempty"`
	Title  string `json:"title,omitempty"`
	Text   string `json:"text,omitempty"`
	Path   string `json:"path,omitempty"`
}

// ASCIIData ascii-emoji object
type ASCIIData struct {
	ID       string `json:"id,omitempty" gorm:"primary_key"`
	Source   string `json:"source,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	Emoji    string `json:"emoji,omitempty"`
}

// Open will open a Database or create it if it doesn't exist
func Open() (*Database, error) {

	var err error
	var sql *gorm.DB
	var index bleve.Index

	if _, err = os.Stat(bleveIndex); os.IsNotExist(err) {
		index, err = bleve.New(bleveIndex, bleve.NewIndexMapping())
		if err != nil {
			return nil, err
		}
		index.Close()
	}

	index, err = bleve.Open(bleveIndex)
	if err != nil {
		return nil, err
	}

	sql, err = gorm.Open("sqlite3", sqliteFile)
	if err != nil {
		return nil, err
	}

	sql.AutoMigrate(&ImageMetaData{})
	sql.AutoMigrate(&ASCIIData{})

	return &Database{sql, index}, nil
}

// Close closes a Database
func (db *Database) Close() {
	db.IDX.Close()
	db.SQL.Close()
}

// WriteImageToDatabase upserts image metadata into Database
func (db *Database) WriteImageToDatabase(image ImageMetaData, itype string) error {

	// add to sqlite
	db.SQL.Create(&image)
	// add to bleve
	err := db.IDX.Index(image.ID, image)
	if err != nil {
		return err
	}

	return nil
}

// WriteASCIIToDatabase upserts ascii metadata into Database
func (db *Database) WriteASCIIToDatabase(ascii ASCIIData) error {

	// add to sqlite
	db.SQL.Create(&ascii)
	// add to bleve
	err := db.IDX.Index(ascii.ID, ascii)
	if err != nil {
		return err
	}

	return nil
}

// Finalize makes index read only optimized
func Finalize() error {
	return nil
}
