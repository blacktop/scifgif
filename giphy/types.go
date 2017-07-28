package giphy

// Search represents a search response from the Giphy API
type Search struct {
	Data       []Gif      `json:"data"`
	Meta       Meta       `json:"meta"`
	Pagination Pagination `json:"pagination"`
}

// Gif giphy gif structure
type Gif struct {
	Type       string `json:"type,omitempty"`
	ID         string `json:"id,omitempty"`
	Slug       string `json:"slug,omitempty"`
	URL        string `json:"url,omitempty"`
	BitlyURL   string `json:"bitly_url,omitempty"`
	EmbedURL   string `json:"embed_url,omitempty"`
	Username   string `json:"username,omitempty"`
	Source     string `json:"source,omitempty"`
	Rating     string `json:"rating,omitempty"`
	Caption    string `json:"caption,omitempty"`
	ContentURL string `json:"content_url,omitempty"`

	Tags         []string `json:"tags,omitempty"`
	FeaturedTags []string `json:"featured_tags,omitempty"`

	User User `json:"user,omitempty"`

	SourceTld     string `json:"source_tld,omitempty"`
	SourcePostURL string `json:"source_post_url,omitempty"`

	UpdateDatetime   string `json:"update_datetime,omitempty"`
	CreateDatetime   string `json:"create_datetime,omitempty"`
	ImportDatetime   string `json:"import_datetime,omitempty"`
	TrendingDatetime string `json:"trending_datetime,omitempty"`

	Images Images `json:"images,omitempty"`
}

// Images giphy image structure
type Images struct {
	FixedHeight            gipyImageDataExtended       `json:"fixed_height,omitempty"`
	FixedHeightStill       gipyImageData               `json:"fixed_height_still,omitempty"`
	FixedHeightDownsampled gipyImageDataExtended       `json:"fixed_height_downsampled,omitempty"`
	FixedWidth             gipyImageDataExtended       `json:"fixed_width,omitempty"`
	FixedWidthStill        gipyImageData               `json:"fixed_width_still,omitempty"`
	FixedWidthDownsampled  gipyImageDataExtended       `json:"fixed_width_downsampled,omitempty"`
	FixedHeightSmall       gipyImageDataExtended       `json:"fixed_height_small,omitempty"`
	FixedHeightSmallStill  gipyImageData               `json:"fixed_height_small_still,omitempty"`
	FixedWidthSmall        gipyImageDataExtended       `json:"fixed_width_small,omitempty"`
	FixedWidthSmallStill   gipyImageData               `json:"fixed_width_small_still,omitempty"`
	Downsized              gipyImageDataSized          `json:"downsized,omitempty"`
	DownsizedStill         gipyImageData               `json:"downsized_still,omitempty"`
	DownsizedLarge         gipyImageDataSized          `json:"downsized_large,omitempty"`
	DownsizedMedium        gipyImageDataSized          `json:"downsized_medium,omitempty"`
	DownsizedSmall         gipyImageDataSized          `json:"downsized_small,omitempty"`
	Original               gipyImageDataExtendedFrames `json:"original,omitempty"`
	OriginalStill          gipyImageData               `json:"original_still,omitempty"`
	Looping                gipyImageLooping            `json:"looping,omitempty"`
	Preview                gipyImagePreview            `json:"preview,omitempty"`
	PreviewGif             gipyImageDataSized          `json:"preview_gif,omitempty"`
}

// User The User Object
type User struct {
}

// Meta The Meta Object
type Meta struct {
	Status     int    `json:"status,omitempty"`
	Msg        string `json:"msg,omitempty"`
	ResponseID string `json:"response_id,omitempty"`
}

// Pagination The Pagination Object
type Pagination struct {
	TotalCount int `json:"total_count,omitempty"`
	Count      int `json:"count,omitempty"`
	Offset     int `json:"offset,omitempty"`
}

type gipyImageLooping struct {
	Mp4 string `json:"mp4,omitempty"`
}

type gipyImageData struct {
	URL    string `json:"url,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

type gipyImageDataSized struct {
	URL    string `json:"url,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Size   string `json:"size,omitempty"`
}

type gipyImagePreview struct {
	Mp4     string `json:"mp4,omitempty"`
	Mp4Size string `json:"mp4_size,omitempty"`
	Width   string `json:"width,omitempty"`
	Height  string `json:"height,omitempty"`
}

type gipyImageDataExtended struct {
	URL      string `json:"url,omitempty"`
	Width    string `json:"width,omitempty"`
	Height   string `json:"height,omitempty"`
	Size     string `json:"size,omitempty"`
	Mp4      string `json:"mp4,omitempty"`
	Mp4Size  string `json:"mp4_size,omitempty"`
	Webp     string `json:"webp,omitempty"`
	WebpSize string `json:"webp_size,omitempty"`
}

type gipyImageDataExtendedFrames struct {
	URL      string `json:"url,omitempty"`
	Width    string `json:"width,omitempty"`
	Height   string `json:"height,omitempty"`
	Size     string `json:"size,omitempty"`
	Frames   string `json:"frames,omitempty"`
	Mp4      string `json:"mp4,omitempty"`
	Mp4Size  string `json:"mp4_size,omitempty"`
	Webp     string `json:"webp,omitempty"`
	WebpSize string `json:"webp_size,omitempty"`
}
