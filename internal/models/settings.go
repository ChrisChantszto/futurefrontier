package models

type SMTPSettings struct {
	AllowCustom bool       `bson:"allowCustom" json:"allowCustom"`
	Selected    string     `bson:"selected" json:"selected"` // "customer" | "company"
	Company     SMTPConfig `bson:"company" json:"company"`
	Customer    *SMTPConfig `bson:"customer,omitempty" json:"customer,omitempty"`
}

type SMTPConfig struct {
	Host      string `bson:"host" json:"host"`
	Port      int    `bson:"port" json:"port"`
	Username  string `bson:"username" json:"username"`
	Password  string `bson:"password" json:"password"`
	FromName  string `bson:"fromName" json:"fromName"`
	FromEmail string `bson:"fromEmail" json:"fromEmail"`
	UseTLS    bool   `bson:"useTLS" json:"useTLS"`
}

type Settings struct {
	ID           string       `bson:"_id" json:"id"`
	Theme        string       `bson:"theme" json:"theme"` // light|dark|system
	FaviconURL   string       `bson:"faviconUrl" json:"faviconUrl"`
	Languages    []string     `bson:"languages" json:"languages"`
	PageLimit    int          `bson:"pageLimit" json:"pageLimit"`
	UserLimit    int          `bson:"userLimit" json:"userLimit"`
	LanguageLimit int         `bson:"languageLimit" json:"languageLimit"`
	SMTP         SMTPSettings `bson:"smtp" json:"smtp"`
}