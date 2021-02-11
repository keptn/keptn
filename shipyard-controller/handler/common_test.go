package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetErrorResponses(t *testing.T) {

	t.Run("TestSetNotFoundErrorResponse_withoutMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetNotFoundErrorResponse(errors.New("an-error"), c)

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusNotFound, r.Code)
		assert.Equal(t, "an-error", *createdErr.Message)
	})

	t.Run("TestSetNotFoundErrorResponse_withMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetNotFoundErrorResponse(errors.New("an-error"), c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusNotFound, r.Code)
		assert.Equal(t, "message: an-error", *createdErr.Message)
	})

	t.Run("TestSetNotFoundErrorResponse_withoutError", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetNotFoundErrorResponse(nil, c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusNotFound, r.Code)
		assert.Equal(t, "message", *createdErr.Message)
	})

	t.Run("TestSetInternalServerErrorResponse_withoutMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetInternalServerErrorResponse(errors.New("an-error"), c)

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		assert.Equal(t, "an-error", *createdErr.Message)
	})

	t.Run("TestSetInternalServerErrorResponse_withMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetInternalServerErrorResponse(errors.New("an-error"), c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		assert.Equal(t, "message: an-error", *createdErr.Message)
	})

	t.Run("TestSetInternalServerErrorResponse_withoutError", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetInternalServerErrorResponse(nil, c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		assert.Equal(t, "message", *createdErr.Message)
	})

	t.Run("SetBadRequestErrorResponse_withoutMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetBadRequestErrorResponse(errors.New("an-error"), c)

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusBadRequest, r.Code)
		assert.Equal(t, "an-error", *createdErr.Message)
	})

	t.Run("SetBadRequestErrorResponse_withMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetBadRequestErrorResponse(errors.New("an-error"), c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusBadRequest, r.Code)
		assert.Equal(t, "message: an-error", *createdErr.Message)
	})

	t.Run("SetBadRequestErrorResponse_withoutError", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetBadRequestErrorResponse(nil, c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusBadRequest, r.Code)
		assert.Equal(t, "message", *createdErr.Message)
	})

	t.Run("SetConflictErrorResponse_withoutMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetConflictErrorResponse(errors.New("an-error"), c)

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Equal(t, "an-error", *createdErr.Message)
	})

	t.Run("SetConflictErrorResponse_withMessage", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetConflictErrorResponse(errors.New("an-error"), c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Equal(t, "message: an-error", *createdErr.Message)
	})

	t.Run("SetConflictErrorResponse_withoutError", func(t *testing.T) {
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		SetConflictErrorResponse(nil, c, "message")

		createdErr := &models.Error{}
		responseBytes, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(responseBytes, createdErr)

		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Equal(t, "message", *createdErr.Message)
	})
}

func TestSetInternalServerErrorResponse(t *testing.T) {

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	SetInternalServerErrorResponse(errors.New("an-error"), c)

	createdErr := &models.Error{}
	responseBytes, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(responseBytes, createdErr)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "an-error", *createdErr.Message)

	// --
	r = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(r)
	SetInternalServerErrorResponse(errors.New("an-error"), c, "message")

	createdErr = &models.Error{}
	responseBytes, _ = ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(responseBytes, createdErr)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "message: an-error", *createdErr.Message)

	// --

	r = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(r)
	SetInternalServerErrorResponse(nil, c, "message")

	createdErr = &models.Error{}
	responseBytes, _ = ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(responseBytes, createdErr)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "message", *createdErr.Message)
}
