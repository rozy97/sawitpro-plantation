package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_PostEstate(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{Repository: mockRepo}

	t.Run("failed test case: invalid request body", func(t *testing.T) {
		invalidBody := `{"length": -5, "width": 10}`
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader([]byte(invalidBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		if assert.NoError(t, server.PostEstate(ctx)) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		}
	})

	t.Run("failed test case: error create estate", func(t *testing.T) {
		requestBody := generated.CreateEstateRequest{Length: 10, Width: 20}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		mockRepo.EXPECT().CreateEstate(gomock.Any(), repository.Estate{Length: 10, Width: 20}).Return("", errors.New("error create estate"))

		if assert.NoError(t, server.PostEstate(ctx)) {
			assert.Equal(t, http.StatusInternalServerError, res.Code)
		}
	})

	t.Run("success test case", func(t *testing.T) {
		requestBody := generated.CreateEstateRequest{Length: 10, Width: 20}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		id := uuid.New().String()
		mockRepo.EXPECT().CreateEstate(gomock.Any(), repository.Estate{Length: 10, Width: 20}).Return(id, nil)

		if assert.NoError(t, server.PostEstate(ctx)) {
			assert.Equal(t, http.StatusCreated, res.Code)

			var responseBody map[string]string
			json.Unmarshal(res.Body.Bytes(), &responseBody)
			assert.Equal(t, id, responseBody["id"])
		}
	})
}

func Test_PostEstateIdTree(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	s := &Server{Repository: mockRepo}

	validEstateID := uuid.New().String()

	t.Run("failed test case: invalid estate ID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		req := httptest.NewRequest(http.MethodPost, "/estate/"+invalidID+"/tree", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		err := s.PostEstateIdTree(ctx, invalidID)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("failed test case: invalid request body", func(t *testing.T) {
		reqBody, _ := json.Marshal(generated.CreateTreeRequest{X: -1, Y: 5, Height: 10})
		req := httptest.NewRequest(http.MethodPost, "/estate/"+validEstateID+"/tree", bytes.NewReader(reqBody))
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		err := s.PostEstateIdTree(ctx, validEstateID)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("failed test case: estate not found", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{}, sql.ErrNoRows)

		reqBody, _ := json.Marshal(generated.CreateTreeRequest{X: 5, Y: 5, Height: 10})
		req := httptest.NewRequest(http.MethodPost, "/estate/"+validEstateID+"/tree", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.PostEstateIdTree(ctx, validEstateID)

		t.Log("Response Body:", res.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("failed test case: tree already exists", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).
			Return(repository.Estate{ID: validEstateID, Length: 10, Width: 10}, nil)
		mockRepo.EXPECT().CreateTree(gomock.Any(), gomock.Any()).
			Return("", errors.New(`pq: duplicate key value violates unique constraint "trees_estate_id_x_y_idx"`))

		reqBody, _ := json.Marshal(generated.CreateTreeRequest{X: 5, Y: 5, Height: 10})
		req := httptest.NewRequest(http.MethodPost, "/estate/"+validEstateID+"/tree", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		err := s.PostEstateIdTree(ctx, validEstateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)

		var responseBody map[string]string
		json.Unmarshal(res.Body.Bytes(), &responseBody)
		assert.Equal(t, "Tree already exist", responseBody["message"])
	})

	t.Run("success case", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).
			Return(repository.Estate{ID: validEstateID, Length: 10, Width: 10}, nil)
		mockRepo.EXPECT().CreateTree(gomock.Any(), gomock.Any()).
			Return("tree-id", nil)

		reqBody, _ := json.Marshal(generated.CreateTreeRequest{X: 5, Y: 5, Height: 10})
		req := httptest.NewRequest(http.MethodPost, "/estate/"+validEstateID+"/tree", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.PostEstateIdTree(ctx, validEstateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		var responseBody map[string]string
		json.Unmarshal(res.Body.Bytes(), &responseBody)
		assert.Equal(t, "tree-id", responseBody["id"])
	})

}

func Test_GetEstateIdStats(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	s := &Server{Repository: mockRepo}

	validEstateID := uuid.New().String()

	t.Run("failed test case: invalid estate ID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		req := httptest.NewRequest(http.MethodGet, "/estate/"+invalidID+"/stats", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdStats(ctx, invalidID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("failed test case: estate not found", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{}, sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/stats", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdStats(ctx, validEstateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("failed test case: repository error", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{ID: validEstateID}, nil)
		mockRepo.EXPECT().GetEstateStats(gomock.Any(), validEstateID).Return(repository.Stats{}, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/stats", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdStats(ctx, validEstateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("success case", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{ID: validEstateID}, nil)
		mockRepo.EXPECT().GetEstateStats(gomock.Any(), validEstateID).Return(repository.Stats{
			TotalTrees: 100,
			MaxHeight:  30,
			MinHeight:  5,
			Median:     15,
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/stats", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdStats(ctx, validEstateID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func Test_GetEstateIdDronePlan(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	s := &Server{Repository: mockRepo}

	validEstateID := uuid.New().String()

	t.Run("failed test case: estate not found", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{}, sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/drone-plan", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdDronePlan(ctx, validEstateID, generated.GetEstateIdDronePlanParams{})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("failed test case: internal server error on GetEstateByID", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{}, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/drone-plan", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdDronePlan(ctx, validEstateID, generated.GetEstateIdDronePlanParams{})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("failed test case: internal server error on GetEstateTrees", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{ID: validEstateID}, nil)
		mockRepo.EXPECT().GetEstateTrees(gomock.Any(), validEstateID).Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/drone-plan", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdDronePlan(ctx, validEstateID, generated.GetEstateIdDronePlanParams{})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("success case", func(t *testing.T) {
		mockRepo.EXPECT().GetEstateByID(gomock.Any(), validEstateID).Return(repository.Estate{ID: validEstateID}, nil)
		mockRepo.EXPECT().GetEstateTrees(gomock.Any(), validEstateID).Return(nil, sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/estate/"+validEstateID+"/drone-plan", nil)
		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)
		ctx.SetRequest(req)

		err := s.GetEstateIdDronePlan(ctx, validEstateID, generated.GetEstateIdDronePlanParams{})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
