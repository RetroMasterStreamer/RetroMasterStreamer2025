package entity

type User struct {
	Name          string `bson:"name"`
	Alias         string `bson:"alias"`
	Password      string `bson:"password"`
	ReferenceText string `bson:"reference_text"`
	UserRef       string `bson:"user_ref"`
}

type RRSS struct {
	Type string `bson:"type"`
	URL  string `bson:"url"`
}

type PostNew struct {
	Title string `bson:"title"`
	Text  string `bson:"text"`
	URL   string `bson:"url"`
}
