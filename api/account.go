package api

import (
	"database/sql"
	"net/http"

	db "github.com/galib612/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"

	_ "github.com/swaggo/swag/example/celler/httputil"
	_ "github.com/swaggo/swag/example/celler/model"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR RUP"`
}

//Info -- when using gin, everything we do inside a handler will involve the context object
//        it provide lots of covenient methods to read inputs parameters and write responses.

// createAccount godoc
// @Summary      Creare an account
// @Description  Create New Account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        AccountReq  body		CreateAccountRequest  true  "Create account"
// @Success      200  {object}  db.Account
// @Failure      400  {string}  string    "error"
// @Failure      404  {string}  string    "error"
// @Failure      500  {string}  string    "error"
// @Security     BasicAuth
// @Router       /accounts/create/ [post]
func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.BindJSON(&req); err != nil {
		// Second Object is json object that want to send to client
		// So here we need a function to convert this err into key value object so that gin an
		// serialize it to JSOn before returning to the client.
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	account, err := server.postgresStore.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
	//ctx.SetCookie() https://github.com/gin-gonic/gin#set-and-get-a-cookie

}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount godoc
// @Summary      Get an account
// @Description  Get Account by id
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path		int		true  "Get Account Req"
// @Success      200  {object}  db.Account
// @Failure      400  {string}  string    "error"
// @Failure      404  {string}  string    "error"
// @Failure      500  {string}  string    "error"
// @Router       /accounts/get/{id} [get]
func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.postgresStore.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type ListAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

// listAccounts godoc
// @Summary      List accounts
// @Description  List Account Based on page_no and page_id
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        listAccountReq    query     ListAccountRequest  true  "List Account Request"
// @Success      200  {array}   db.Account
// @Failure      400  {string}  string "error"
// @Failure      404  {string}  string "error"
// @Failure      500  {string}  string "error"
// @Router       /accounts/list/ [get]
func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, err := server.postgresStore.ListAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)

}

type UpdateAccountRequest struct {
	ID      int64 `json:"id" binding:"required"`
	Balance int64 `json:"balance" binding:"required"`
}

// updateAccount godoc
// @Summary      Update an account
// @Description  Update by json account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        updateAccountReq	body	UpdateAccountRequest	true  "Update Account Req"
// @Success      200      {object}  db.Account
// @Failure      400  {string}  string "error"
// @Failure      404  {string}  string "error"
// @Failure      500  {string}  string "error"
// @Router       /accounts/update/ [post]
func (server *Server) updateAccount(ctx *gin.Context) {
	var req UpdateAccountRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}
	account, err := server.postgresStore.UpdateAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteAccount godoc
// @Summary      Delete an Account
// @Description  Delete Account by id
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path		int		true  "Delete Account Req"
// @Success      200  {string}  string	  "Ok"
// @Failure      400  {string}  string    "error"
// @Failure      404  {string}  string    "error"
// @Failure      500  {string}  string    "error"
// @Router       /accounts/delete/{id} [get]
func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	affectedrow, err := server.postgresStore.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	if affectedrow > 0 {
		msg := "Account is Deleted Successfully"
		ctx.JSON(http.StatusOK, msg)
	} else {
		msg := "Account Not found"
		ctx.JSON(http.StatusOK, msg)
	}
}
