package handler

import (
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"

	"golang.org/x/crypto/bcrypt"

	"lheinrich.de/secpass/conf"
	"lheinrich.de/secpass/shorts"
	"lheinrich.de/secpass/spuser"
)

// Login function
func Login(w http.ResponseWriter, r *http.Request) {
	// check logged in and redirect to index if so
	if checkSession(r) != "" {
		// logout
		if r.URL.Path == "/login/logout" {
			// delete cookie
			delete(spuser.Sessions, cookie(r, "secpass_uuid"))

			// redirect to login page
			redirect(w, "/login")
			return
		}

		// redirect to index
		redirect(w, "/")
		return
	}

	// output login data wrong
	special := 0

	// define values
	name, password := r.PostFormValue("name"), r.PostFormValue("password")

	// check for input
	if name != "" && password != "" && len(password) >= 8 {
		// read password hash from database and check for error
		var username string
		var passwordHash string
		err := conf.DB.QueryRow(conf.GetSQL("login"), name).Scan(&username, &passwordHash)
		shorts.Check(err, true)

		// compare passwords
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil {
			// define two-factor authentication variables
			oneTimePassword := r.PostFormValue("oneTimePassword")
			twoFactorSecret := spuser.TwoFactorSecret(username)
			var skipLogin bool

			// check two-factor authentication is enabled
			if twoFactorSecret != "" {
				// validate one-time password
				if len(oneTimePassword) < 6 || len(oneTimePassword) > 8 || !totp.Validate(oneTimePassword, twoFactorSecret) {
					// skip log-in and output error
					skipLogin = true
					special = 2
				}
			}

			if !skipLogin {
				// generate uuid and define session expiration date
				uuid := shorts.UUID()
				expires := time.Now().Add(time.Hour)

				// define cookies
				cookieUUID := http.Cookie{Name: "secpass_uuid", Value: uuid, Expires: expires}
				cookieName := http.Cookie{Name: "secpass_name", Value: username, Expires: expires}
				cookieHash := http.Cookie{Name: "secpass_hash", Value: shorts.Hash(password), Expires: expires}

				// add session and set cookies
				spuser.Sessions[uuid] = spuser.Session{User: username, Expires: expires}
				http.SetCookie(w, &cookieUUID)
				http.SetCookie(w, &cookieName)
				http.SetCookie(w, &cookieHash)

				// redirect to index
				w.Header().Set("location", "/")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
		} else {
			special = 1
		}
	}

	// execute template
	shorts.Check(tpl.ExecuteTemplate(w, "login.html", Data{User: "", Lang: getLang(r), Special: special, LoggedOut: true}), false)
}
