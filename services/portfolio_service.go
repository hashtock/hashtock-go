package services

// import (
// 	"net/http"

// 	"github.com/hashtock/hashtock-go/core"
// 	"github.com/hashtock/hashtock-go/models"
// )

// func Portfolio(req *http.Request, r render.Render) {
// 	profile, _ := models.GetProfile(req)
// 	shares, err := models.GetProfileShares(req, profile)
// 	if err != nil {
// 		r.JSON(core.ErrToErrorer(err))
// 		return
// 	}

// 	r.JSON(http.StatusOK, shares)
// }

// func PortfolioTagInfo(req *http.Request, params martini.Params, r render.Render) {
// 	hash_tag_name := params["tag"]

// 	profile, err := models.GetProfile(req)
// 	if err != nil {
// 		r.JSON(core.ErrToErrorer(err))
// 		return
// 	}

// 	share, err := models.GetProfileShareByTagName(req, profile, hash_tag_name)
// 	if err != nil {
// 		r.JSON(core.ErrToErrorer(err))
// 		return
// 	}

// 	r.JSON(http.StatusOK, share)
// }
