package models

type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int64  `json:"population"`
}

func (c Country) Validate() bool {
	if c.Name == "" || c.Capital == "" || c.Currency == "" {
		return false
	}
	return true
}
