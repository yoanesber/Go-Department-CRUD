package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yoanesber/Go-Department-CRUD/controller"
	"github.com/yoanesber/Go-Department-CRUD/model"
)

// SetupSampleDepartment creates a sample department object for testing purposes
// It returns a Department struct with sample data
func GetSampleDepartment() model.Department {
	return model.Department{
		ID:          "d001",
		DeptName:    "HR",
		Active:      true,
		CreatedBy:   1,
		CreatedDate: time.Now(),
		UpdatedBy:   1,
		UpdatedDate: time.Now(),
	}
}

// SetupSampleDepartments creates a slice of sample department objects for testing purposes
// It returns a slice of Department structs with sample data
func GetSampleDepartments() []model.Department {
	return []model.Department{
		{
			ID:          "d001",
			DeptName:    "HR",
			Active:      true,
			CreatedBy:   1,
			CreatedDate: time.Now(),
			UpdatedBy:   1,
			UpdatedDate: time.Now(),
		},
		{
			ID:          "d002",
			DeptName:    "IT",
			Active:      true,
			CreatedBy:   1,
			CreatedDate: time.Now(),
			UpdatedBy:   1,
			UpdatedDate: time.Now(),
		},
	}
}

// MockService is a mock implementation of the DepartmentService interface for testing purposes.
type MockService struct{}

// Mock implementation of the DepartmentService.GetAllDepartments method
// This method returns a list of departments for testing purposes
func (m *MockService) GetAllDepartments() ([]model.Department, error) {
	return GetSampleDepartments(), nil
}

// Mock implementation of the DepartmentService.GetDepartmentByID method
// This method returns a single department for testing purposes
func (m *MockService) GetDepartmentByID(id string) (model.Department, error) {
	return GetSampleDepartment(), nil
}

// Mock implementation of the DepartmentService.CreateDepartment method
// This method creates a new department for testing purposes
func (m *MockService) CreateDepartment(department model.Department) (model.Department, error) {
	return department, nil
}

// Mock implementation of the DepartmentService.UpdateDepartment method
// This method updates an existing department for testing purposes
func (m *MockService) UpdateDepartment(id string, department model.Department) (model.Department, error) {
	return department, nil
}

// Mock implementation of the DepartmentService.DeleteDepartment method
// This method deletes a department for testing purposes
func (m *MockService) DeleteDepartment(id string) (bool, error) {
	return true, nil
}

// SetupRouter initializes the Gin router and sets up the routes for department management
// It uses the MockService for testing purposes
func SetupRouter() *gin.Engine {
	// Define a mock service for testing
	mock := &MockService{}
	ctrl := &controller.DepartmentController{Service: mock}

	// Create a new Gin router instance
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	// Set up the API version group
	// This group contains all the routes for the API version 1
	v1 := r.Group("/api/v1")
	{
		// Routes for department management
		// These routes handle CRUD operations for departments
		dept := v1.Group("/departments")
		{
			dept.GET("", ctrl.GetAllDepartments)
			dept.GET("/:id", ctrl.GetDepartmentByID)
			dept.POST("", ctrl.CreateDepartment)
			dept.PUT("/:id", ctrl.UpdateDepartment)
			dept.DELETE("/:id", ctrl.DeleteDepartment)
		}
	}

	return r
}

func ConvertHttpResponseToDepartment(t *testing.T, resp *httptest.ResponseRecorder) (model.Department, error) {
	// Unmarshal the response body into a HttpResponse object
	// The HttpResponse object contains the data returned by the server
	var httpResponse model.HttpResponse
	err := json.Unmarshal(resp.Body.Bytes(), &httpResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Marshal the Data back to JSON
	// This is done to ensure that the Data can be converted back to JSON format
	jsonData, err := json.Marshal(httpResponse.Data)
	if err != nil {
		t.Fatalf("Failed to marshal Data back to JSON: %v", err)
	}

	// Unmarshal the JSON data into a Department object
	// This is done to ensure that the Data can be converted back to a Department object
	var dept model.Department
	err = json.Unmarshal(jsonData, &dept)
	if err != nil {
		return model.Department{}, err
	}

	return dept, nil
}

func TestGetAllDepartments(t *testing.T) {
	r := SetupRouter()

	// Create a new HTTP request to the endpoint
	// The request is a GET request to the "/departments" endpoint with no body
	req, err := http.NewRequest("GET", "/api/v1/departments", nil)
	if err != nil {
		t.Fatalf("Failed to get all departments: %v", err)
	}

	// Set the request header
	req.Header.Set("Accept", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	// Create a new HTTP response recorder to capture the response
	// The response recorder is used to simulate an HTTP response for testing purposes
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check if the response status code is 200 OK
	// This means the request was successful and the server returned the expected response
	assert.Equal(t, http.StatusOK, resp.Code, "Expected status code 200 OK")

	// Unmarshal the response body into a HttpResponse object
	// The HttpResponse object contains the data returned by the server
	var httpResponse model.HttpResponse
	err = json.Unmarshal(resp.Body.Bytes(), &httpResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check if the response data is not empty
	// This means the server returned a list of departments as expected
	assert.NotEmpty(t, httpResponse.Data, "Expected departments list to be not empty")
}

func TestGetDepartmentByID(t *testing.T) {
	r := SetupRouter()

	// Create a new HTTP request to the endpoint
	// The request is a GET request to the "/departments/{id}" endpoint with no body
	req, err := http.NewRequest("GET", "/api/v1/departments/"+GetSampleDepartment().ID, nil)
	if err != nil {
		t.Fatalf("Failed to get department by ID: %v", err)
	}

	// Set the request header
	req.Header.Set("Accept", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	// Create a new HTTP response recorder to capture the response
	// The response recorder is used to simulate an HTTP response for testing purposes
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check if the response status code is 200 OK
	// This means the request was successful and the server returned the expected response
	assert.Equal(t, http.StatusOK, resp.Code, "Expected status code 200 OK")

	// Convert the response to a Department object
	dept, err := ConvertHttpResponseToDepartment(t, resp)
	if err != nil {
		t.Fatalf("Failed to convert response to Department: %v", err)
	}

	// Check if the department ID and name match the expected values
	// This is done to ensure that the Data contains the expected values
	assert.Equal(t, GetSampleDepartment().ID, dept.ID, "Expected department ID to match")
	assert.Equal(t, GetSampleDepartment().DeptName, dept.DeptName, "Expected department name to match")
}

func TestCreateDepartment(t *testing.T) {
	r := SetupRouter()

	// Sample department data
	newDept := model.Department{
		ID:        GetSampleDepartment().ID,
		DeptName:  GetSampleDepartment().DeptName,
		Active:    GetSampleDepartment().Active,
		CreatedBy: GetSampleDepartment().CreatedBy,
	}
	jsonData, _ := json.Marshal(newDept)

	// Create a new HTTP request to the endpoint
	// The request is a POST request to the "/departments" endpoint with the department data in the body
	req, err := http.NewRequest("POST", "/api/v1/departments", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create department: %v", err)
	}

	// Set the request header
	// The request header specifies the content type of the request body
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	// Create a new HTTP response recorder to capture the response
	// The response recorder is used to simulate an HTTP response for testing purposes
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check if the response status code is 201 Created
	// This means the request was successful and the server created a new department
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Convert the response to a Department object
	createdDept, err := ConvertHttpResponseToDepartment(t, resp)
	if err != nil {
		t.Fatalf("Failed to convert response to Department: %v", err)
	}

	// Check if the created department ID matches the expected ID
	// This is done to ensure that the created department has the expected ID
	assert.Equal(t, newDept.ID, createdDept.ID, "Expected created department ID to match")
	assert.Equal(t, newDept.DeptName, createdDept.DeptName, "Expected created department name to match")
}

func TestUpdateDepartment(t *testing.T) {
	r := SetupRouter()

	// Sample update data
	updateDept := model.Department{
		DeptName:  GetSampleDepartment().DeptName,
		Active:    GetSampleDepartment().Active,
		UpdatedBy: GetSampleDepartment().UpdatedBy,
	}
	jsonData, _ := json.Marshal(updateDept)

	// Create a new HTTP request to the endpoint
	// The request is a PUT request to the "/departments/{id}" endpoint with the updated department data in the body
	req, err := http.NewRequest("PUT", "/api/v1/departments/"+GetSampleDepartment().ID, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to update department: %v", err)
	}

	// Set the request header
	// The request header specifies the content type of the request body
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	// Create a new HTTP response recorder to capture the response
	// The response recorder is used to simulate an HTTP response for testing purposes
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check if the response status code is 200 OK
	// This means the request was successful and the server updated the department
	assert.Equal(t, http.StatusOK, resp.Code)

	// Convert the response to a Department object
	updatedDept, err := ConvertHttpResponseToDepartment(t, resp)
	if err != nil {
		t.Fatalf("Failed to convert response to Department: %v", err)
	}

	// Check if the updated department ID matches the expected ID
	// This is done to ensure that the updated department has the expected ID
	assert.Equal(t, updateDept.DeptName, updatedDept.DeptName, "Expected updated department name to match")
}

func TestDeleteDepartment(t *testing.T) {
	r := SetupRouter()

	// Create a new HTTP request to the endpoint
	// The request is a DELETE request to the "/departments/{id}" endpoint with no body
	req, err := http.NewRequest("DELETE", "/api/v1/departments/"+GetSampleDepartment().ID, nil)
	if err != nil {
		t.Fatalf("Failed to delete department: %v", err)
	}

	// Set the request header
	// The request header specifies the content type of the request body
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", accessToken))

	// Create a new HTTP response recorder to capture the response
	// The response recorder is used to simulate an HTTP response for testing purposes
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check if the response status code is 200 OK
	// This means the request was successful and the server deleted the department
	assert.Equal(t, http.StatusOK, resp.Code)
}
