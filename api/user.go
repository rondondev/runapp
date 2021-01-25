package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/rondondev/runapp/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Type  string  `json:"type" binding:"required,oneof=admin coach athlete"`
	Name  string  `json:"name" binding:"required"`
	Email string  `json:"email" binding:"required"`
	Phone *string `json:"phone"`
	Birth *string `json:"birth" binding:"datetime=2006-01-02"`
}

func (r *createUserRequest) toDB() (db.CreateUserParams, error) {
	arg := db.CreateUserParams{
		Type:  db.UserType(r.Type),
		Name:  r.Name,
		Email: r.Email,
	}
	if r.Phone != nil {
		arg.Phone.SetValid(*r.Phone)
	}
	if r.Birth != nil {
		t, err := time.Parse("2006-01-02", *r.Birth)
		if err != nil {
			return db.CreateUserParams{}, err
		}
		arg.Birth.SetValid(t)
	}

	return arg, nil
}

type updateUserRequest struct {
	createUserRequest
	Active *bool `json:"active" binding:"required"`
}

func (r *updateUserRequest) toDB(id int64) (db.UpdateUserParams, error) {
	arg := db.UpdateUserParams{
		ID:     id,
		Type:   db.UserType(r.Type),
		Name:   r.Name,
		Email:  r.Email,
		Active: *r.Active,
	}
	if r.Phone != nil {
		arg.Phone.SetValid(*r.Phone)
	}
	if r.Birth != nil {
		t, err := time.Parse("2006-01-02", *r.Birth)
		if err != nil {
			return db.UpdateUserParams{}, err
		}
		arg.Birth.SetValid(t)
	}

	return arg, nil
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg, err := req.toDB()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=10,max=1000"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) listActiveUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListActiveUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListActiveUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) listAllUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAllUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListAllUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type deleteUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if user exists
	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete the user
	err = server.store.DeleteUser(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type updateUserIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var u updateUserIDRequest
	if err := ctx.ShouldBindUri(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the user exists
	user, err := server.store.GetUser(ctx, u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg, err := req.toDB(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updated, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
