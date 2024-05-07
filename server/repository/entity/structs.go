package entity

type User struct {
	Name          string
	Alias         string
	Password      string
	ReferenceText string
	UserRef       string
}

type RRSS struct {
	Type string
	URL  string
}

type PostNew struct {
	Title string
	Text  string
	URL   string
}
