package giphy

import (
	"fmt"
	"strings"
)

// Search Search all GIPHY GIFs for a word or phrase. Punctuation will be stripped and ignored.
// Use a plus or url encode for phrases. Example paul+rudd, ryan+gosling or american+psycho.
func (c *Client) Search(args []string) (Search, error) {
	argsStr := strings.Join(args, " ")

	path := fmt.Sprintf("/gifs/search?limit=%v&q=%s", c.Limit, argsStr)
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
