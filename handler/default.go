package handler

import (
	"html/template"
	"net/http"

	"lheinrich.de/secpass/conf"
	"lheinrich.de/secpass/shorts"
	"lheinrich.de/secpass/spuser"
)

// Data to pass into template
type Data struct {
	User      string
	Lang      string
	Special   int
	LoggedOut bool
	TwoFactor TwoFactorData
	Passwords map[string]string
}

// TwoFactorData data for two-factor authentication
type TwoFactorData struct {
	Image                string
	Secret               string
	Disabled             bool
	OneTimePasswordWrong bool
}

var (
	// define functions
	funcs = template.FuncMap{
		"config": func(keys ...string) string {
			return conf.Config[keys[0]][keys[1]]
		},
		"lang": func(keys ...interface{}) string {
			return (*conf.Lang[keys[0].(string)])[keys[1].(string)]
		},
		"languages": func() []string {
			languages := []string{}
			for lang := range conf.Lang {
				languages = append(languages, lang)
			}

			return languages
		}}

	// define templates
	tpl *template.Template
)

// name function
func name(w http.ResponseWriter, r *http.Request) {
	// execute template
	shorts.Check(tpl.ExecuteTemplate(w, "name", nil), false)
}

// LoadTemplates parse
func LoadTemplates() {
	// parse templates and check error
	var err error
	tpl, err = template.New("").Funcs(funcs).ParseGlob(conf.Config["webserver"]["templatesDirectory"])
	shorts.Check(err, true)
}

// redirect to location
func redirect(w http.ResponseWriter, location string) {
	w.Header().Set("location", location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// cookiesExist check whether the cookies exist
func cookiesExist(r *http.Request, names ...string) bool {
	for _, name := range names {
		_, err := r.Cookie(name)
		if err != nil {
			return false
		}
	}
	return true
}

// getLang get language or default language
func getLang(r *http.Request) string {
	lang, err := r.Cookie("secpass_lang")
	if err != nil {
		return conf.Config["app"]["defaultLanguage"]
	}
	return lang.Value
}

// checkSession check whether logged in and return username or empty string
func checkSession(r *http.Request) string {
	if cookiesExist(r, "secpass_uuid", "secpass_name") {
		cookieUUID, _ := r.Cookie("secpass_uuid")
		cookieName, _ := r.Cookie("secpass_name")
		uuid := cookieUUID.Value
		name := cookieName.Value
		user := spuser.Sessions[uuid].User

		if user == name {
			return user
		}
	}
	return ""
}

// cookie return cookie value
func cookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	shorts.Check(err, false)

	return cookie.Value
}
