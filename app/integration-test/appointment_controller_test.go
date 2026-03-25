package integrationtest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"scenario-a/internal/config"
	"scenario-a/internal/config/sqldb"
	"scenario-a/internal/dep"
	"scenario-a/internal/dto/requestdto"
	"scenario-a/internal/dto/responsedto"
	"scenario-a/internal/model"
	"scenario-a/internal/route"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AppointmentControllerTestSuite struct {
	suite.Suite
	openDB func(t *testing.T) (*sql.DB, *gorm.DB)

	e    *echo.Echo
	deps *dep.Dependencies
}

func (suite *AppointmentControllerTestSuite) SetupSuite() {
	config.MustInitForTest()
	cfg := config.Get()
	suite.openDB = sqldb.MustInitForTest(cfg.Env)
}

func (suite *AppointmentControllerTestSuite) SetupTest() {
	_, db := suite.openDB(suite.T())

	suite.deps = dep.Init(db)
	suite.e = echo.New()
	route.Load(suite.e, suite.deps)
}

// Helper method to make HTTP requests
func (suite *AppointmentControllerTestSuite) makeRequest(method, path string, payload any) *httptest.ResponseRecorder {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		suite.Require().NoError(err)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if payload != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)
	return rec
}

func TestAppointmentControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AppointmentControllerTestSuite))
}

func (suite *AppointmentControllerTestSuite) Test_ListAppointments() {
	rec := suite.makeRequest(http.MethodGet, "/appointments", nil)

	suite.Require().Equal(http.StatusOK, rec.Code)

	var resp responsedto.ListAppointmentsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp.Items)
}

func (suite *AppointmentControllerTestSuite) Test_BookAppointment() {
	// Use a fixed future date that is guaranteed to be a Monday to ensure Dealership is open (Jan 3, 2050)
	startAt := time.Date(2050, time.January, 3, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		payload       requestdto.BookAppointmentRequest
		expectedCode  int
		checkResponse func(rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			payload: requestdto.BookAppointmentRequest{
				CustomerID:    1,
				DealershipID:  1,
				ServiceTypeID: 1, // 60 mins duration
				TechnicianID:  1, // Uses dealership 1, can do service type 1
				StartAt:       startAt,
			},
			expectedCode: http.StatusOK,
			checkResponse: func(rec *httptest.ResponseRecorder) {
				var resp responsedto.Appointment
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				suite.Require().NoError(err)
				suite.Require().NotEmpty(resp.UUID)
				suite.Require().Equal(string(model.AppointmentStatusCreated), resp.Status)
			},
		},
		{
			name:    "fail validation",
			payload: requestdto.BookAppointmentRequest{
				// Missing fields to trigger validation error
			},
			expectedCode:  http.StatusBadRequest,
			checkResponse: func(rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			rec := suite.makeRequest(http.MethodPost, "/appointments", tt.payload)
			suite.Require().Equal(tt.expectedCode, rec.Code, "Response: %s", rec.Body.String())
			if tt.checkResponse != nil {
				tt.checkResponse(rec)
			}
		})
	}
}

func (suite *AppointmentControllerTestSuite) Test_CancelAppointment() {
	// Book an appointment first directly via repository
	startAt := time.Date(2050, time.January, 3, 10, 0, 0, 0, time.UTC)
	appointment := &model.Appointment{
		CommonModel:   model.CommonModel{UUID: "test-cancel-uuid"},
		CustomerID:    1,
		DealershipID:  1,
		ServiceBayID:  1,
		TechnicianID:  1,
		ServiceTypeID: 1,
		Status:        model.AppointmentStatusCreated,
		StartAt:       startAt,
		EndAt:         startAt.Add(1 * time.Hour),
	}
	err := suite.deps.Repository.Appointment.Create(suite.T().Context(), appointment)
	suite.Require().NoError(err)
	suite.Require().NotZero(appointment.ID)

	// Cancel it via API using helper
	payload := requestdto.CancelAppointmentRequest{
		Description: "Changing my mind",
	}
	rec := suite.makeRequest(http.MethodPost, fmt.Sprintf("/appointments/%s/cancel", appointment.UUID), payload)

	suite.Require().Equal(http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

	var resp responsedto.Appointment
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Require().Equal(string(model.AppointmentStatusCancelled), resp.Status)
	suite.Require().Equal("Changing my mind", *resp.Description)
}

func (suite *AppointmentControllerTestSuite) Test_CompleteAppointment() {
	// Book an appointment first directly via repository
	startAt := time.Date(2050, time.January, 4, 10, 0, 0, 0, time.UTC) // Avoid overlap if run in parallel
	appointment := &model.Appointment{
		CommonModel:   model.CommonModel{UUID: "test-complete-uuid"},
		CustomerID:    1,
		DealershipID:  1,
		ServiceBayID:  1,
		TechnicianID:  1,
		ServiceTypeID: 1,
		Status:        model.AppointmentStatusCreated,
		StartAt:       startAt,
		EndAt:         startAt.Add(1 * time.Hour),
	}
	err := suite.deps.Repository.Appointment.Create(suite.T().Context(), appointment)
	suite.Require().NoError(err)
	suite.Require().NotZero(appointment.ID)

	// Complete it via API using helper
	payload := requestdto.CompleteAppointmentRequest{
		Description: "All done",
	}
	rec := suite.makeRequest(http.MethodPost, fmt.Sprintf("/appointments/%s/complete", appointment.UUID), payload)

	suite.Require().Equal(http.StatusOK, rec.Code, "Response: %s", rec.Body.String())

	var resp responsedto.Appointment
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	suite.Require().NoError(err)
	suite.Require().Equal(string(model.AppointmentStatusCompleted), resp.Status)
	suite.Require().Equal("All done", *resp.Description)
	suite.T().Logf("Response body: %s", rec.Body.String())
}

func (suite *AppointmentControllerTestSuite) Test_BookAppointment_Concurrent() {
	cfg := config.Get()
	sqlDB, gormDB := sqldb.MustInit(cfg.Env)
	defer sqlDB.Close()

	realDeps := dep.Init(gormDB)
	realE := echo.New()
	route.Load(realE, realDeps)

	// Use a fixed future Monday at 9:00 AM (Jan 3, 2050)
	startAt := time.Date(2050, time.January, 3, 9, 0, 0, 0, time.UTC)

	// Clean up after test
	defer gormDB.Exec("DELETE FROM appointments WHERE start_at = ?", startAt)

	payload := requestdto.BookAppointmentRequest{
		CustomerID:    1,
		DealershipID:  1,
		ServiceTypeID: 1, // 60 mins duration
		TechnicianID:  1, // Uses dealership 1, can do service type 1
		StartAt:       startAt,
	}

	concurrentRequests := 5

	var wg sync.WaitGroup
	var mu sync.Mutex
	statusCodes := make([]int, 0, concurrentRequests)

	// Using a barrier ensures all goroutines fire their requests at nearly the exact same time
	startBarrier := make(chan struct{})

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/appointments", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			<-startBarrier // Wait for the start signal

			realE.ServeHTTP(rec, req)

			mu.Lock()
			statusCodes = append(statusCodes, rec.Code)
			mu.Unlock()
		}()
	}

	// Release the barrier to start all concurrent requests
	close(startBarrier)
	wg.Wait()

	successCount := 0
	failCount := 0
	for _, code := range statusCodes {
		if code == http.StatusOK {
			successCount++
		} else {
			failCount++
		}
	}

	suite.Require().Equal(1, successCount, "Only one booking should succeed")
	suite.Require().Equal(concurrentRequests-1, failCount, "Other bookings should fail due to conflict")
}
