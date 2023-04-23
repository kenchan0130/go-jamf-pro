package classic

type Icon struct {
	ID       *int    `xml:"id,omitempty"`
	Filename *string `xml:"filename,omitempty"`
	URI      *string `xml:"uri,omitempty"`
}
