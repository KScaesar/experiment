package api

import (
	"github.com/gin-gonic/gin"

	"experiment/swag/domain"
)

// CreateUser
// @Summary Create User
// @version 1.0
// @consumer json
// @produce json
// @param Authorization header string true "'Bearer token'"
// @param RequestBody body domain.User true "RequestBody"
// @Success 200 {object} object{code=integer,msg=string}
// @Router /api/user [POST]
func CreateUser(c *gin.Context) {
	_ = new(domain.User)
}

// GetUser
// @version 1.0
// @consumer json
// @produce json
// @param Authorization header string true "'Bearer token'"
// @param id path string true "user id"
// @param UrlQueryString query domain.UserParam false "UrlQueryString"
// @Success 200 {object} NormalResp{data=[]domain.User}
// @Router /api/user/{id} [get]
func GetUser(c *gin.Context) {
	_ = new(domain.UserParam)
}

type NormalResp struct {
	Code int
	Msg  string
	Data interface{}
}
