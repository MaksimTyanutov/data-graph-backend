package utils

import (
	"data-graph-backend/pkg/dataStructers"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func ToFile(str string, name string) {
	file, err := os.Create(name + ".json")

	if err != nil {
		fmt.Println("Unable to create file:", err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println("Unable to create file:", err.Error())
	}
}

func ValidateDates(filters dataStructers.Filters) error {
	if _, err := time.Parse(time.RFC3339, filters.CompanyFilters.MinDate); err != nil {
		return errors.New("Date value is incorrect: " + filters.CompanyFilters.MinDate + ". Error: " + err.Error())
	}
	if _, err := time.Parse(time.RFC3339, filters.CompanyFilters.MaxDate); err != nil {
		return errors.New("Date value is incorrect: " + filters.CompanyFilters.MinDate + ". Error: " + err.Error())
	}
	if _, err := time.Parse(time.RFC3339, filters.ProductFilters.MinDate); err != nil {
		return errors.New("Date value is incorrect: " + filters.CompanyFilters.MinDate + ". Error: " + err.Error())
	}
	if _, err := time.Parse(time.RFC3339, filters.ProductFilters.MaxDate); err != nil {
		return errors.New("Date value is incorrect: " + filters.CompanyFilters.MinDate + ". Error: " + err.Error())
	}
	return nil
}

func ValidateFilter(filters dataStructers.Filters) error {
	err := ValidateDates(filters)
	if err != nil {
		err_ := errors.New("Wrong date format: " + err.Error())
		log.Print(err_)
		return err_
	}
	if filters.CompanyFilters.StartStaffSize > filters.CompanyFilters.EndStaffSize || filters.CompanyFilters.StartStaffSize < 0 || filters.CompanyFilters.EndStaffSize < 0 {
		err_ := errors.New("company Staff size has incorrect format")
		log.Print(err_)
		return err_
	}
	return nil
}
