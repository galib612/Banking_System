package api

import (
	"errors"
	"net/http"

	authdb "github.com/galib612/simplebank/authdb/sqlc"
	middleware "github.com/galib612/simplebank/middleware"
	"github.com/gin-gonic/gin"

	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
)

type CreateuserRequest struct {
	UserName string `json:"username" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateuserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := authdb.CreateUserParams{
		Username: req.UserName,
		Passwd:   req.Passwd,
	}

	err := server.authdbStore.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "User Singup is Done!")

}

type CreateloginRequest struct {
	UserName string `json:"username" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req CreateloginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := authdb.GetUserParams{
		Username: req.UserName,
		Passwd:   req.Passwd,
	}

	err := server.authdbStore.GetUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "User is Authenticated")

	// Generate JWT token
	generatedToken, err := middleware.GenerateJwtToken(middleware.Payload{Username: arg.Username})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	} else if len(generatedToken) == 0 {
		ctx.JSON(http.StatusInternalServerError, errResponse(errors.New("GENERATE TOKEN IS EMPTY")))
		return
	}
	// Set Cookies

	ctx.SetCookie("token", generatedToken, 3600, "/", "localhost", false, true)

}
