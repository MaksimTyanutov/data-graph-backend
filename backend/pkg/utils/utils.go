package utils

import (
	"data-graph-backend/pkg/dataStructers"
	"errors"
	"log"
	"time"
)

func ValidateDates(MinDate string, MaxDate string) error {
	var min, max time.Time
	var err error
	if min, err = time.Parse(time.RFC3339, MinDate); err != nil {
		return errors.New("Date value is incorrect: " + MinDate + ". Error: " + err.Error())
	}
	if max, err = time.Parse(time.RFC3339, MaxDate); err != nil {
		return errors.New("Date value is incorrect: " + MaxDate + ". Error: " + err.Error())
	}
	if min.After(max) {
		return errors.New("Date value is incorrect: MinDate (" + MinDate + ") is later than MaxDate (" + MaxDate + ")")
	}
	return nil
}

func ValidateFilterCompany(companyFilters dataStructers.CompanyFilters) error {
	err := ValidateDates(companyFilters.MinDate, companyFilters.MaxDate)
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
	err := ValidateDates(productFilters.MinDate, productFilters.MaxDate)
	if err != nil {
		err_ := errors.New("Wrong date format: " + err.Error())
		log.Print(err_)
		return err_
	}
	return nil
}

func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
