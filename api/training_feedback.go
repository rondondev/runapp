package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/rondondev/runapp/db/sqlc"

	"github.com/gin-gonic/gin"
)

type listTrainingFeedbackRequest struct {
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`
}

func (server *Server) listTrainingFeedbacksByUser(ctx *gin.Context) {
	var u idRequest
	if err := ctx.ShouldBindUri(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req listTrainingFeedbackRequest
	period := true
	start := ctx.Query("start_date")
	end := ctx.Query("end_date")
	if start == "" && end == "" {
		period = false
	} else if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var feedbacks []db.TrainingFeedback
	var err error
	if period {
		// we can ignore the errors because the values were already validated
		start, _ := time.Parse("2006-01-02", req.StartDate)
		end, _ := time.Parse("2006-01-02", req.EndDate)
		arg := db.ListTrainingFeedbacksByUserInPeriodParams{
			UserID: u.ID,
			Date:   start,
			Date_2: end,
		}

		feedbacks, err = server.store.ListTrainingFeedbacksByUserInPeriod(ctx, arg)
	} else {
		feedbacks, err = server.store.ListTrainingFeedbacksByUser(ctx, u.ID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, feedbacks)
}

func (server *Server) getTrainingFeedback(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	feedback, err := server.store.GetTrainingFeedback(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, feedback)
}

type createTrainingFeedbackRequest struct {
	BorgScale      int32            `json:"borg_scale" binding:"required,min=6,max=20"`
}

func (r *createTrainingFeedbackRequest) toDB(trainingID int64) (db.CreateTrainingFeedbackParams, error) {
	arg := db.CreateTrainingFeedbackParams{
		TrainingID: trainingID,
		BorgScale: r.BorgScale,
	}
	return arg, nil
}

func (server *Server) createTrainingFeedback(ctx *gin.Context) {
	var t idRequest
	if err := ctx.ShouldBindUri(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req createTrainingFeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the training exists
	_, err := server.store.GetTraining(ctx, t.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("invalid training_id")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg, err := req.toDB(t.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	training, err := server.store.CreateTrainingFeedback(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, training)
}

func (server *Server) deleteTrainingFeedback(ctx *gin.Context) {
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

	// Check if the feedback exists
	_, err = server.store.GetTrainingFeedback(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Delete the feedback
	err = server.store.DeleteTrainingFeedback(ctx, training.ID)
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

type updateTrainingFeedbackRequest struct {
	createTrainingFeedbackRequest
}

func (r *updateTrainingFeedbackRequest) toDB(trainingID int64) (db.UpdateTrainingFeedbackParams, error) {
	arg := db.UpdateTrainingFeedbackParams{
		TrainingID: trainingID,
		BorgScale: r.BorgScale,
	}
	return arg, nil
}

func (server *Server) updateTrainingFeedback(ctx *gin.Context) {
	var t idRequest
	if err := ctx.ShouldBindUri(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateTrainingFeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if the training exists
	training, err := server.store.GetTraining(ctx, t.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg, err := req.toDB(training.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updated, err := server.store.UpdateTrainingFeedback(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
