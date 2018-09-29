package handler

import (
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"

	"golang.org/x/crypto/bcrypt"

	"lheinrich.de/lheinrich/secpass/conf"
	"lheinrich.de/lheinrich/secpass/shorts"
	"lheinrich.de/lheinrich/secpass/spuser"
)

// Login function
func Login(w http.ResponseWriter, r *http.Request) {
	// check logged in and redirect to index if so
	if checkSession(w, r) != "" {
		// logout
		if r.URL.Path == "/login/logout" {
			// delete session
			delete(spuser.Sessions, getCookie(r, "secpass_uuid"))

			// delete cookies
			deleteAllCookies(w)

			// redirect to login page
			redirect(w, "/login")
			return
		}

		// redirect to index
		redirect(w, "/")
		return
	}

	// output login data wrong, crypter key
	special := 0
	var crypter string

	// define values
	name, password := r.PostFormValue("name"), r.PostFormValue("password")

	// check for input
	if name != "" && password != "" && len(password) >= 8 {
		// define variables to write into
		var username string
		var passwordHash string
		var twoFactorSecret string
		var key string

		// read data from database
		conf.DB.QueryRow(conf.GetSQL("login"), name).Scan(&username, &passwordHash, &twoFactorSecret, &key)

		// compare passwords
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil {
			// define two-factor authentication variables
			oneTimePassword := r.PostFormValue("oneTimePassword")
			var skipLogin bool

			// check two-factor authentication is enabled
			if twoFactorSecret != "" {
				// validate one-time password
				if (len(oneTimePassword) != 6 && len(oneTimePassword) != 8) || !totp.Validate(oneTimePassword, twoFactorSecret) {
					// skip log-in and output error
					skipLogin = true
					special = 2
				}
			}

			if !skipLogin {
				// generate uuid and define session expiration date
				uuid := shorts.UUID()
				expires := time.Now().Add(10 * time.Minute)

				// define cookies
				cookieUUID := http.Cookie{Name: "secpass_uuid", Value: uuid, Expires: expires}
				cookieName := http.Cookie{Name: "secpass_name", Value: username, Expires: expires}

				// add session and set cookies
				spuser.Sessions[uuid] = spuser.Session{User: username, Expires: expires}
				http.SetCookie(w, &cookieUUID)
				http.SetCookie(w, &cookieName)

				// send crypter
				crypter = key
			}
		} else {
			// wrong password
			special = 1
		}
	}

	// execute template
	shorts.Check(tpl.ExecuteTemplate(w, "login.html", Data{User: "", Lang: getLang(r), Special: special,
		LoggedOut: true, Crypter: crypter}))
}
