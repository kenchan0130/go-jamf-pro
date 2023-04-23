package classic

type GeneralCategory struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type SelfServiceCategory struct {
	ID        *int    `xml:"id,omitempty"`
	Name      *string `xml:"name,omitempty"`
	DisplayIn *bool   `xml:"display_in,omitempty"`
	FeatureIn *bool   `xml:"feature_in,omitempty"`
}
