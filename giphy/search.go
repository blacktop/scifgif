package giphy

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Search Search all GIPHY GIFs for a word or phrase. Punctuation will be stripped and ignored.
// Use a plus or url encode for phrases. Example paul+rudd, ryan+gosling or american+psycho.
func (c *Client) Search(args []string, offset int) (Search, error) {
	argsStr := strings.Join(args, "+")

	path := fmt.Sprintf("/gifs/search?limit=%v&offset=%d&rating=%s&q=%s", c.Limit, offset, c.Rating, argsStr)

	log.WithFields(log.Fields{
		"query":  argsStr,
		"limit":  c.Limit,
		"offset": offset,
		"rating": c.Rating,
	}).Debug("searching Giphy")

	req, err := c.NewRequest(path)
	if err != nil {
		return Search{}, err
	}

	var search Search
	if _, err = c.Do(req, &search); err != nil {
		return Search{}, err
	}

	return search, nil
}
