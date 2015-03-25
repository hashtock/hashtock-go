package models

import (
    "errors"
    "net/http"

    "github.com/gorilla/context"
    netContext "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "google.golang.org/api/plus/v1"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/core"
)

const (
    reqProfile = "reqProfile"
)

var (
    ErrUserNotFound = errors.New("User not found")
)

// User has to be registered and logged at this point!
// Auth middleware should take care of that
func GetProfile(req *http.Request) (profile *Profile, err error) {
    val, ok := context.GetOk(req, reqProfile)
    if ok && val.(*Profile) != nil {
        return val.(*Profile), nil
    }

    err = core.NewForbiddenError()
    return
}

func GetOrRegisterProfile(req *http.Request, userId string, authToken string) (profile *Profile, err error) {
    if val, ok := context.GetOk(req, reqProfile); ok {
        return val.(*Profile), nil
    }

    var gProfile *Profile
    if userId == "" {
        gProfile, err = fetchProfileFromGoogle(authToken)
        if err != nil {
            return
        }
        userId = gProfile.UserID
    }

    profile, err = getProfileForUserId(req, userId)
    if err == ErrUserNotFound {
        if gProfile != nil {
            profile, err = registerProfile(req, gProfile)
        } else {
            err = core.NewNotFoundError(err.Error())
        }
    }

    if err != nil {
        return
    }

    context.Set(req, reqProfile, profile)
    return
}

func fetchProfileFromGoogle(authToken string) (profile *Profile, err error) {
    gauthConf := conf.GetConfig().GAuthConfig()
    token := new(oauth2.Token)
    token.AccessToken = authToken

    ctx := netContext.Background()
    client := gauthConf.Client(ctx, token)
    plusService, err := plus.New(client)
    if err != nil {
        return
    }

    gperson, err := plusService.People.Get("me").Do()
    if err != nil {
        return
    }

    accountEmail := ""
    for _, email := range gperson.Emails {
        if email.Type == "account" {
            accountEmail = email.Value
        }
    }

    if accountEmail == "" {
        err = errors.New("User does not have valid email in Google profile")
        return
    }

    profile = &Profile{
        UserID: accountEmail,
    }

    return
}

func getProfileForUserId(req *http.Request, userID string) (profile *Profile, err error) {
    col := storage.Collection(ProfileCollectionName)
    defer col.Database.Session.Close()

    profile = new(Profile)
    selector := bson.M{
        "user_id": userID,
    }
    err = col.Find(selector).One(profile)
    if err == mgo.ErrNotFound {
        err = ErrUserNotFound
    }

    return
}

func registerProfile(req *http.Request, profileDetails *Profile) (profile *Profile, err error) {
    col := storage.Collection(ProfileCollectionName)
    defer col.Database.Session.Close()

    if profileDetails.UserID == "" {
        err = errors.New("Could not register user. No ID found")
        return
    }

    profileDetails.Founds = StartingFounds
    profileDetails.IsAdmin = conf.IsAdmin(profileDetails.UserID)
    err = col.Insert(profileDetails)

    return profileDetails, err
}
