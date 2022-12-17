package utils

import (
	"data-graph-backend/pkg/dataStructers"
	"errors"
	"log"
	"time"
)

func ValidateDates(MinDate string, MaxDate string) error {
	if _, err := time.Parse(time.RFC3339, MinDate); err != nil {
		return errors.New("Date value is incorrect: " + MinDate + ". Error: " + err.Error())
	}
	if _, err := time.Parse(time.RFC3339, MaxDate); err != nil {
		return errors.New("Date value is incorrect: " + MinDate + ". Error: " + err.Error())
	}
	return nil
}

func ValidateFilterCompany(companyFilters dataStructers.CompanyFilters) error {
	err := ValidateDates(companyFilters.MaxDate, companyFilters.MinDate)
	if err != nil {
		err_ := errors.New("Wrong date format: " + err.Error())
		log.Print(err_)
		return err_
	}
	if companyFilters.StartStaffSize > companyFilters.EndStaffSize || companyFilters.StartStaffSize < 0 || companyFilters.EndStaffSize < 0 {
		err_ := errors.New("company Staff size has incorrect format")
		log.Print(err_)
		return err_
	}
	return nil
}

func ValidateFilterProduct(productFilters dataStructers.ProductFilters) error {
	err := ValidateDates(productFilters.MaxDate, productFilters.MinDate)
	if err != nil {
		err_ := errors.New("Wrong date format: " + err.Error())
		log.Print(err_)
		return err_
	}
	return nil
}
