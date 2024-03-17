package httpAuth

import (
	"net/http"
	"time"
)

const (
	CookieName = "session_id"
)

func (h AuthHandler) makeHTTPCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:  CookieName,
		Value: sessionID,
		Expires: time.Now().
			AddDate(int(h.cookieSettings.ExpireDate.Years),
				int(h.cookieSettings.ExpireDate.Months),
				int(h.cookieSettings.ExpireDate.Days)),
		Secure:   false,
		HttpOnly: h.cookieSettings.HttpOnly,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
}

func GetCookie(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}
