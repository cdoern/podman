package reports

type ScpReport struct {
	Id  string `json:id` //nolint
	Err error  `json:Err, omitempty`
}
