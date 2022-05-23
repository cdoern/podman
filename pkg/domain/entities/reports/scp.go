package reports

type ScpReport struct {
	Id     string `json:"Id"` //nolint
	Source string `json:"Source_Image"`
	Dest   string `json:"Destination_Tag"`
	Err    error  `json:"Err,omitempty"`
}
