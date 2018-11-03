package controllers

import (
    "log"
    "errors"
    "time"
    "net/http"
    "gopkg.in/mgo.v2/bson"
    "github.com/oskarszura/smarthome/models"
    "github.com/oskarszura/smarthome/utils"
    ctrutl "github.com/oskarszura/smarthome/controllers/utils"
    "github.com/oskarszura/gowebserver/session"
    "github.com/oskarszura/gowebserver/router"
)

// Authenticate - handle login page and login process
func Authenticate(w http.ResponseWriter, r *http.Request, opt router.UrlOptions, sm session.ISessionManager) {
    defer r.Body.Close()

    switch r.Method {
    case "GET":
        ctrutl.RenderTemplate(w, r, "login", sm)

    case "POST":
        sessionID, _ := utils.GetSessionID(r)
        isLogged := sm.IsExist(sessionID)

        if !isLogged {
            user := r.PostFormValue("username")
            password := r.PostFormValue("password")
            expiration := time.Now().Add(365 * 24 * time.Hour)

            authenticatedUser, authErr := authenticateUser(user, password)

            if authErr == nil {
                cookieValue := user + password

                cookie := http.Cookie {
                    Name: "sid",
                    Value: cookieValue,
                    Expires: expiration }

                session := sm.Create(cookieValue)
                session.Set("user", authenticatedUser)

                http.SetCookie(w, &cookie)
                http.Redirect(w, r, "/", http.StatusSeeOther)
            }
        }

        http.Redirect(w, r, "/login", http.StatusSeeOther)

    default:
    }
}

func authenticateUser(user string, password string) (models.User, error) {
    var loggedUser models.User

    ds := utils.GetDataSource()
    c := ds.C("users")

    err := c.Find(bson.M{
        "username": user,
        "password": password,
    }).One(&loggedUser)

    if err != nil {
        log.Println("User not found err=", err)
        return models.User{}, errors.New("foobar")
    }

    log.Println("Logged in as user", loggedUser)

    return loggedUser, nil
}
