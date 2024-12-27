package main

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetStudentNameByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance of DBInterface
	mockDB := NewMockDBInterface(ctrl)

	// Define test cases
	tests := []struct {
		name          string
		id            int
		mockBehavior  func(mockDB *MockDBInterface)
		expectedName  string
		expectedError error
	}{
		{
			name: "Success - Valid ID",
			id:   101,
			mockBehavior: func(mockDB *MockDBInterface) {
				mockDB.EXPECT().
					First(gomock.Any(), 101).
					DoAndReturn(func(out interface{}, where ...interface{}) *gorm.DB {
						// Simulate setting the student data
						student := out.(*Student)
						student.ID = 101
						student.Name = "Iman"
						return &gorm.DB{} // Return a success response
					})
			},
			expectedName:  "Iman",
			expectedError: nil,
		},
		{
			name: "Failure - Record Not Found",
			id:   102,
			mockBehavior: func(mockDB *MockDBInterface) {
				mockDB.EXPECT().
					First(gomock.Any(), 102).
					DoAndReturn(func(out interface{}, where ...interface{}) *gorm.DB {
						return &gorm.DB{Error: gorm.ErrRecordNotFound}
					})
			},
			expectedName:  "",
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Failure - Database Error",
			id:   103,
			mockBehavior: func(mockDB *MockDBInterface) {
				mockDB.EXPECT().
					First(gomock.Any(), 103).
					DoAndReturn(func(out interface{}, where ...interface{}) *gorm.DB {
						return &gorm.DB{Error: errors.New("database connection failed")}
					})
			},
			expectedName:  "",
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "Failure - Invalid Data Type",
			id:   -1,
			mockBehavior: func(mockDB *MockDBInterface) {
				mockDB.EXPECT().
					First(gomock.Any(), -1).
					DoAndReturn(func(out interface{}, where ...interface{}) *gorm.DB {
						return &gorm.DB{Error: errors.New("invalid data type")}
					})
			},
			expectedName:  "",
			expectedError: errors.New("invalid data type"),
		},
		{
			name: "Failure - Nil Student Struct",
			id:   104,
			mockBehavior: func(mockDB *MockDBInterface) {
				mockDB.EXPECT().
					First(gomock.Any(), 104).
					DoAndReturn(func(out interface{}, where ...interface{}) *gorm.DB {
						return &gorm.DB{Error: errors.New("student struct is nil")}
					})
			},
			expectedName:  "",
			expectedError: errors.New("student struct is nil"),
		},
	}

	// Execute each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockDB) // Set up mock behavior

			// Call the function under test
			name, err := GetStudentNameByID(mockDB, tt.id)

			// Validate the results
			assert.Equal(t, tt.expectedName, name)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
