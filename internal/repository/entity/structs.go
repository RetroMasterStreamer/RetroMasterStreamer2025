package entity

type User struct {
	Name          string `bson:"name" json:"name"`
	Alias         string `bson:"alias" json:"alias"`
	Password      string `bson:"password,omitempty" json:"-"`
	ReferenceText string `bson:"reference_text" json:"reference_text"`
	UserRef       string `bson:"user_ref" json:"user_ref"`
	AboutMe       string `bson:"about_me" json:"about_me"`
	RRSS          []RRSS `bson:"rrss" json:"RRSS"`
}

type RRSS struct {
	Type string `bson:"type" json:"type"`
	URL  string `bson:"url" json:"URL"`
}

type UserOnline struct {
	Alias        string `bson:"alias" json:"alias"`
	SessionToken string `bson:"sesion" json:"-"`
	Hash         string `bson:"hash" json:"hash"`
	Online       bool   `bson:"online" json:"online"`
}

type PostNew struct {
	Title string `bson:"title" json:"title"`
	Text  string `bson:"text" json:"text"`
	URL   string `bson:"url" json:"URL"`
}
