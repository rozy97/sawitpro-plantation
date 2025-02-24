package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/helper"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// // This is just a test endpoint to get you started. Please delete this endpoint.
// // (GET /hello)
// func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
// 	var resp generated.HelloResponse
// 	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
// 	return ctx.JSON(http.StatusOK, resp)
// }

// Endpoint Create /estate
// (POST /estate)
func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.CreateEstateRequest
	// Bind request body to struct
	if err := ctx.Bind(&req); err != nil || req.Length <= 0 || req.Width <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid Request Body"})
	}

	id, err := s.Repository.CreateEstate(ctx.Request().Context(), repository.Estate{Length: req.Length, Width: req.Width})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusCreated, generated.CreateResponse{Id: id})
}

// Get Estate Drone Plan
// (GET /estate/{id}/drone-plan)
func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id string, params generated.GetEstateIdDronePlanParams) error {

	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "Estate not found"})
		default:
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}

	trees, err := s.Repository.GetEstateTrees(ctx.Request().Context(), id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	maxDistance := 0
	countFirstRest := false
	var resp generated.GetEstateDronePlanResponse
	if params.MaxDistance != nil && *params.MaxDistance > 0 {
		maxDistance = *params.MaxDistance
		countFirstRest = true
		x, y := 0, 0
		resp = generated.GetEstateDronePlanResponse{
			Distance: 0,
			Rest: &struct {
				X *int `json:"x,omitempty"`
				Y *int `json:"y,omitempty"`
			}{X: &x, Y: &y},
		}
	}

	if len(trees) == 0 {
		return ctx.JSON(http.StatusOK, resp)
	}

	statsHelper := helper.Stats{
		Estate:         estate,
		Trees:          trees,
		CountFirstRest: countFirstRest,
		MaxDistance:    maxDistance,
	}

	statsHelper.CalculateTotalDistance()
	resp.Distance = statsHelper.Distance
	if countFirstRest {
		resp.Rest.X = &statsHelper.Rest.X
		resp.Rest.Y = &statsHelper.Rest.Y
	}

	return ctx.JSON(http.StatusOK, resp)
}

// Get Estate Stats
// (GET /estate/{id}/stats)
func (s *Server) GetEstateIdStats(ctx echo.Context, id string) error {
	err := uuid.Validate(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid Estate ID"})
	}

	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "Estate not found"})
		default:
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}

	stats, err := s.Repository.GetEstateStats(ctx.Request().Context(), estate.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, generated.GetEstateStatsResponse{
		Count:  stats.TotalTrees,
		Max:    stats.MaxHeight,
		Min:    stats.MinHeight,
		Median: stats.Median,
	})
}

// Create Tree Within Estate
// (POST /estate/{id}/tree)
func (s *Server) PostEstateIdTree(ctx echo.Context, id string) error {
	err := uuid.Validate(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid Estate ID"})
	}

	var req generated.CreateTreeRequest
	// Bind request body to struct
	if err := ctx.Bind(&req); err != nil || req.X <= 0 || req.Y <= 0 || req.Height <= 0 || req.Height > 30 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid Request Body"})
	}

	estate, err := s.Repository.GetEstateByID(ctx.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "Estate not found"})
		default:
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
		}
	}

	if req.X > estate.Length || req.Y > estate.Width {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid Request Body"})
	}

	treeID, err := s.Repository.CreateTree(ctx.Request().Context(), repository.Tree{
		EstateID: estate.ID,
		X:        req.X,
		Y:        req.Y,
		Height:   req.Height,
	})
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "trees_estate_id_x_y_idx"` {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Tree already exist"})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(201, generated.CreateResponse{Id: treeID})
}
