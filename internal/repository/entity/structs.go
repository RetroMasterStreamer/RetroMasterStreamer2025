package entity

type User struct {
	Name          string `bson:"name" json:"name"`
	Alias         string `bson:"alias" json:"alias"`
	AvatarYT      string `bson:"avatar_yt" json:"avatar_yt"`
	Password      string `bson:"password,omitempty" json:"-"`
	ReferenceText string `bson:"reference_text" json:"reference_text"`
	UserRef       string `bson:"user_ref" json:"user_ref"`
	AboutMe       string `bson:"about_me" json:"about_me"`
	RRSS          []RRSS `bson:"rrss" json:"RRSS"`
}

type NewUser struct {
	Name          string `bson:"name" json:"name"`
	Alias         string `bson:"alias" json:"alias"`
	AvatarYT      string `bson:"avatar_yt" json:"avatar_yt"`
	Password      string `bson:"password,omitempty" json:"password"`
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
	ID      string `bson:"id" json:"id"`
	Title   string `bson:"title" json:"title"`
	Content string `bson:"content" json:"content"`
	URL     string `bson:"url" json:"url"`
	Type    string `bson:"type" json:"type"`
	Author  string `bson:"author" json:"author"`
	Date    string `bson:"date" json:"date"`
}
