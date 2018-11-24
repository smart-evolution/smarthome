package utils

import (
    "fmt"
    "time"
    "net/http"
    "crypto/sha1"
)

// CreateSessionID - creates a new session ID
func CreateSessionID(user string, pass string, time string) string {
    val := []byte(user + pass + time)
    h := sha1.New()
    h.Write(val)

    return fmt.Sprintf("%x", h.Sum(nil))
}

// GetSessionID - get user session ID
func GetSessionID(r *http.Request) (string, error) {
    sessionCookie, err := r.Cookie("sid")

    if err != nil {
        return "", err
    }

    return sessionCookie.Value, nil
}

// ClearSession - remove session cookie
func ClearSession(w http.ResponseWriter) {
    cookie := http.Cookie {
        Path: "/",
        Name: "sid",
        Expires: time.Now().Add(-100 * time.Hour),
        MaxAg
    http.SetCookie(w, &cookie)
}
