package posts

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	// "time"

	"github.com/jinzhu/gorm"
	"github.com/qorpress/l10n"
	"github.com/qorpress/publish2"
	qor_seo "github.com/qorpress/seo"
	"github.com/qorpress/slug"
	"github.com/qorpress/sorting"
	"github.com/qorpress/validations"
	"github.com/qorpress/media/media_library"

	"github.com/qorpress/qorpress-example/pkg/models/seo"
)

type Post struct {
	gorm.Model
	l10n.Locale
	sorting.SortingDESC
	// ID           uint     `gorm:"primary_key"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	Name                  string
	NameWithSlug          slug.Slug `l10n:"sync"`
	Featured              bool
	Code                  string       `l10n:"sync"`
	CategoryID            uint         `l10n:"sync"`
	Category              Category     `l10n:"sync"`
	Collections           []Collection `l10n:"sync" gorm:"many2many:post_collections;"`
	MainImage             media_library.MediaBox
	Images         media_library.MediaBox
	Description           string           `gorm:"type:longtext"`
	PostProperties     PostProperties `sql:"type:text"`
	Seo                   qor_seo.Setting
	publish2.Version
	publish2.Schedule
	publish2.Visible
}

type PostVariation struct {
	gorm.Model
	PostID *uint
	Post   Post
	Featured          bool
	Images            media_library.MediaBox
}

func (post Post) GetID() uint {
	return post.ID
}

func (post Post) GetSEO() *qor_seo.SEO {
	return seo.SEOCollection.GetSEO("Post Page")
}

func (post Post) DefaultPath() string {
	defaultPath := "/"
	return defaultPath
}

func (post Post) MainImageURL(styles ...string) string {
	style := "main"
	if len(styles) > 0 {
		style = styles[0]
	}

	if len(post.MainImage.Files) > 0 {
		return post.MainImage.URL(style)
	}

	return "/images/default_post.png"
}

func (post Post) Validate(db *gorm.DB) {
	if strings.TrimSpace(post.Name) == "" {
		db.AddError(validations.NewError(post, "Name", "Name can not be empty"))
	}

	if strings.TrimSpace(post.Code) == "" {
		db.AddError(validations.NewError(post, "Code", "Code can not be empty"))
	}
}

type PostImage struct {
	gorm.Model
	Title        string
	Category     Category
	CategoryID   uint
	SelectedType string
	File         media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}

func (postImage PostImage) Validate(db *gorm.DB) {
	if strings.TrimSpace(postImage.Title) == "" {
		db.AddError(validations.NewError(postImage, "Title", "Title can not be empty"))
	}
}

func (postImage *PostImage) SetSelectedType(typ string) {
	postImage.SelectedType = typ
}

func (postImage *PostImage) GetSelectedType() string {
	return postImage.SelectedType
}

func (postImage *PostImage) ScanMediaOptions(mediaOption media_library.MediaOption) error {
	if bytes, err := json.Marshal(mediaOption); err == nil {
		return postImage.File.Scan(bytes)
	} else {
		return err
	}
}

func (postImage *PostImage) GetMediaOption() (mediaOption media_library.MediaOption) {
	mediaOption.Video = postImage.File.Video
	mediaOption.FileName = postImage.File.FileName
	mediaOption.URL = postImage.File.URL()
	mediaOption.OriginalURL = postImage.File.URL("original")
	mediaOption.CropOptions = postImage.File.CropOptions
	mediaOption.Sizes = postImage.File.GetSizes()
	mediaOption.Description = postImage.File.Description
	return
}

type PostProperties []PostProperty

type PostProperty struct {
	Name  string
	Value string
}

func (postProperties *PostProperties) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, postProperties)
	case string:
		if v != "" {
			return postProperties.Scan([]byte(v))
		}
	default:
		return errors.New("not supported")
	}
	return nil
}

func (postProperties PostProperties) Value() (driver.Value, error) {
	if len(postProperties) == 0 {
		return nil, nil
	}
	return json.Marshal(postProperties)
}
