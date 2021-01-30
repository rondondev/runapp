package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/rondondev/runapp/db/sqlc"

	"github.com/gin-gonic/gin"
)

type listTrainingRequest struct {
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`
}

func (server *Server) listTrainingsByUser(ctx *gin.Context) {
	var u idRequest
	if err := ctx.ShouldBindUri(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req listTrainingRequest
	period := true
	start := ctx.Query("start_date")
	end := ctx.Query("end_date")
	if start == "" && end == "" {
		period = false
	} else if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var trainings []db.Training
	var err error
	if period {
		// we can ignore the errors because the values were already validated
		start, _ := time.Parse("2006-01-02", req.StartDate)
		end, _ := time.Parse("2006-01-02", req.EndDate)
		arg := db.ListTrainingsByUserInPeriodParams{
			UserID: u.ID,
			Date:   start,
			Date_2: end,
		}

		trainings, err = server.store.ListTrainingsByUserInPeriod(ctx, arg)
	} else {
		trainings, err = server.store.ListTrainingsByUser(ctx, u.ID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trainings)
}

func (server *Server) getTraining(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	training, err := server.store.GetTraining(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, training)
}

type createTrainingRequest struct {
	UserID    int64             `json:"user_id" binding:"required"`
	Date      string            `json:"date" binding:"required,datetime=2006-01-02"`
	Sport     db.TrainingSport  `json:"sport" binding:"required"`
	Type      *string           `json:"type"`
	Intensity *string           `json:"intensity"`
	Details   string            `json:"details" binding:"required"`
	Status    db.TrainingStatus `json:"status"`
}

func (r *createTrainingRequest) toDB() (db.CreateTrainingParams, error) {
	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return db.CreateTrainingParams{}, err
	}
	arg := db.CreateTrainingParams{
		UserID:  r.UserID,
		Date:    d,
		Sport:   r.Sport,
		Details: r.Details,
	}
	if r.Type != nil {
		arg.Type.SetValid(*r.Type)
	}
	if r.Intensity != nil {
		arg.Intensity.SetValid(*r.Intensity)
	}
	if r.Status == "" {
		arg.Status = db.TrainingStatusNew
	} else {
		arg.Status = r.Status
	}

	return arg, nil
}

func (server *Server) createTraining(ctx *gin.Context) {
	var req createTrainingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the user exists
	_, err := server.store.GetUser(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("invalid user_id")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg, err := req.toDB()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	training, err := server.store.CreateTraining(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, training)
}

func (server *Server) deleteTraining(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the training exists
	training, err := server.store.GetTraining(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete the training
	err = server.store.DeleteTraining(ctx, training.ID)
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

type updateTrainingRequest struct {
	Date      string            `json:"date" binding:"required,datetime=2006-01-02"`
	Sport     db.TrainingSport  `json:"sport" binding:"required"`
	Type      *string           `json:"type"`
	Intensity *string           `json:"intensity"`
	Details   string            `json:"details" binding:"required"`
	Status    db.TrainingStatus `json:"status" binding:"required,oneof=new notified overdue done done_feedback"`
}

func (r *updateTrainingRequest) toDB(id int64) (db.UpdateTrainingParams, error) {
	d, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return db.UpdateTrainingParams{}, err
	}
	arg := db.UpdateTrainingParams{
		ID:      id,
		Date:    d,
		Sport:   r.Sport,
		Details: r.Details,
		Status:  r.Status,
	}
	if r.Type != nil {
		arg.Type.SetValid(*r.Type)
	}
	if r.Intensity != nil {
		arg.Intensity.SetValid(*r.Intensity)
	}

	return arg, nil
}

func (server *Server) updateTraining(ctx *gin.Context) {
	var u idRequest
	if err := ctx.ShouldBindUri(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateTrainingRequest
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

	updated, err := server.store.UpdateTraining(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
