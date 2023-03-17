package utils

import (
	"data-graph-backend/pkg/dataStructers"
	"errors"
	"testing"
)

func TestValidateDates(t *testing.T) {
	testTable := []struct {
		name          string
		minDate       string
		maxDate       string
		expectedError error
	}{
		{
			name:          "Correct dates formats",
			minDate:       "1990-01-01T00:00:00Z",
			maxDate:       "2000-01-01T00:00:00Z",
			expectedError: nil,
		},
		{
			name:          "Incorrect minDate formats",
			minDate:       "1990-01-01T00",
			maxDate:       "2000-01-01T00:00:00Z",
			expectedError: errors.New("Date value is incorrect: " + "1990-01-01T00" + ". Error: parsing time \"1990-01-01T00\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \":\""),
		},
		{
			name:          "Incorrect maxDate formats",
			minDate:       "1990-01-01T00:00:00Z",
			maxDate:       "2000-01-01T00",
			expectedError: errors.New("Date value is incorrect: " + "2000-01-01T00" + ". Error: parsing time \"2000-01-01T00\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \":\""),
		},
		{
			name:    "Empty strings",
			minDate: "",
			maxDate: "",
			expectedError: errors.New("Date value is incorrect: " + "" + ". Error: parsing time \"\" " +
				"as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\""),
		},
		{
			name:          "Min bigger than max date",
			minDate:       "2000-01-01T00:00:00Z",
			maxDate:       "1990-01-01T00:00:00Z",
			expectedError: errors.New("Date value is incorrect: MinDate (" + "2000-01-01T00:00:00Z" + ") is later than MaxDate (" + "1990-01-01T00:00:00Z" + ")"),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateDates(test.minDate, test.maxDate)
			if err != nil && test.expectedError == nil {
				t.Errorf("Actual error = %s\nExpected error = nil", err.Error())
				return
			}
			if err == nil && test.expectedError != nil {
				t.Errorf("Actual error = nil\nExpected error = %s", test.expectedError)
				return
			}
			if err != nil && test.expectedError != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Actual error = %s\nExpected error = %s", err.Error(), test.expectedError)
					return
				}
			}
		})
	}
}

func TestValidateFilterCompany(t *testing.T) {
	testTable := []struct {
		name           string
		companyFilters dataStructers.CompanyFilters
		expectedError  error
	}{
		{
			name: "Correct companyFilters",
			companyFilters: dataStructers.CompanyFilters{
				CompanyName:    "Company Name",
				Departments:    []int{1, 2, 3},
				Ceo:            "Ceo name",
				MinDate:        "1990-01-01T00:00:00Z",
				MaxDate:        "2000-01-01T00:00:00Z",
				StartStaffSize: 100,
				EndStaffSize:   10000,
			},
			expectedError: nil,
		},
		{
			name: "Incorrect companyFilters",
			companyFilters: dataStructers.CompanyFilters{
				CompanyName:    "Company Name",
				Departments:    []int{1, 2, 3},
				Ceo:            "Ceo name",
				MinDate:        "1990-01-01T00:00:00Z",
				MaxDate:        "2000-01-01T00:00:00Z",
				StartStaffSize: 10000,
				EndStaffSize:   100,
			},
			expectedError: errors.New("company Staff size has incorrect format"),
		},
		{
			name: "Incorrect companyFilters",
			companyFilters: dataStructers.CompanyFilters{
				CompanyName:    "Company Name",
				Departments:    []int{1, 2, 3},
				Ceo:            "Ceo name",
				MinDate:        "1990-01-01T00:00:00Z",
				MaxDate:        "1000-01-01T00:00:00Z",
				StartStaffSize: 100,
				EndStaffSize:   10000,
			},
			expectedError: errors.New("Wrong date format: Date value is incorrect: MinDate (1990-01-01T00:00:00Z) is later than MaxDate (1000-01-01T00:00:00Z)"),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateFilterCompany(test.companyFilters)
			if err != nil && test.expectedError == nil {
				t.Errorf("Actual error = %s\nExpected error = nil", err.Error())
				return
			}
			if err == nil && test.expectedError != nil {
				t.Errorf("Actual error = nil\nExpected error = %s", test.expectedError)
				return
			}
			if err != nil && test.expectedError != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Actual error = %s\nExpected error = %s", err.Error(), test.expectedError)
					return
				}
			}
		})
	}
}

func TestValidateFilterProduct(t *testing.T) {
	testTable := []struct {
		name           string
		productFilters dataStructers.ProductFilters
		expectedError  error
	}{
		{
			name: "Correct productFilters",
			productFilters: dataStructers.ProductFilters{
				ProductName: "Product Name",
				MinDate:     "1990-01-01T00:00:00Z",
				MaxDate:     "2000-01-01T00:00:00Z",
				Departments: []int{1, 2, 3},
				IsVerified:  true,
			},
			expectedError: nil,
		},
		{
			name: "Incorrect productFilters",
			productFilters: dataStructers.ProductFilters{
				ProductName: "Product Name",
				MinDate:     "1990-01-01T00:00:00Z",
				MaxDate:     "1000-01-01T00:00:00Z",
				Departments: []int{1, 2, 3},
				IsVerified:  true,
			},
			expectedError: errors.New("Wrong date format: Date value is incorrect: MinDate (1990-01-01T00:00:00Z) is later than MaxDate (1000-01-01T00:00:00Z)"),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateFilterProduct(test.productFilters)
			if err != nil && test.expectedError == nil {
				t.Errorf("Actual error = %s\nExpected error = nil", err.Error())
				return
			}
			if err == nil && test.expectedError != nil {
				t.Errorf("Actual error = nil\nExpected error = %s", test.expectedError)
				return
			}
			if err != nil && test.expectedError != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Actual error = %s\nExpected error = %s", err.Error(), test.expectedError)
					return
				}
			}
		})
	}
}

func TestDeleteEmpty(t *testing.T) {
	testTable := []struct {
		name        string
		str         []string
		expectedStr []string
	}{
		{
			name: "Correct delete empty",
			str: []string{
				"st rin       gs ",
				"se  co  nd",
				"third"},
			expectedStr: []string{
				"strings",
				"second",
				"third",
			},
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			str := DeleteEmpty(test.str)
			if str == nil && test.expectedStr != nil {
				t.Errorf("Actual str = nil\nExpected str = %s", test.expectedStr)
				return
			}
			if str != nil && test.expectedStr == nil {
				t.Errorf("Actual str = %s\nExpected error = nil", str)
				return
			}
			if len(str) != len(test.expectedStr) {
				t.Errorf("Actual len = %d\nExpected len = %d", len(str), len(test.expectedStr))
				return
			}
		})
	}
}
