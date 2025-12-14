package validation

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

var phoneRegexp = regexp.MustCompile("^[+][0-9]{11}$")

var (
	InvalidAgeErr   = errors.New("invalid age")
	InvalidPhoneErr = errors.New("invalid phone")
	InvalidNameErr  = errors.New("invalid name")
)

type RegistrationForm struct {
	Age   uint
	Name  string
	Phone string
}

type ValidationError struct {
	ValidAge        bool
	ValidMinLenName bool
	ValidMaxLenName bool
	ValidPhone      bool
}

func (e ValidationError) Error() string {
	if err := errors.Join(e.Unwrap()...); err != nil {
		return err.Error()
	}
	return ""
}

func (e ValidationError) Unwrap() []error {
	var errs []error

	if !e.ValidPhone {
		errs = append(errs, fmt.Errorf("%w: phone is invalid", InvalidPhoneErr))
	}

	if !e.ValidMinLenName {
		errs = append(errs, fmt.Errorf("%w: name is to shoort", InvalidNameErr))
	}

	if !e.ValidMaxLenName {
		errs = append(errs, fmt.Errorf("%w: name is to long", InvalidNameErr))
	}

	if !e.ValidAge {
		errs = append(errs, fmt.Errorf("%w: age must be greater than 18", InvalidAgeErr))
	}

	return errs
}

func (f RegistrationForm) Check() error {
	var errs []error
	if f.Age < 18 {
		errs = append(errs, fmt.Errorf("%w: age must be greater than 18", InvalidAgeErr))
	}

	if utf8.RuneCountInString(f.Name) < 2 {
		errs = append(errs, fmt.Errorf("%w: name is to shoort", InvalidNameErr))
	}

	if utf8.RuneCountInString(f.Name) > 50 {
		errs = append(errs, fmt.Errorf("%w: name is to long", InvalidNameErr))
	}

	if !phoneRegexp.MatchString(f.Phone) {
		errs = append(errs, fmt.Errorf("%w: phone is invalid", InvalidPhoneErr))
	}

	return errors.Join(errs...)
}

func (f RegistrationForm) CheckV2() error {
	var err ValidationError

	if f.Age >= 18 {
		err.ValidAge = true
	}

	if utf8.RuneCountInString(f.Name) >= 2 {
		err.ValidMinLenName = true
	}

	if utf8.RuneCountInString(f.Name) <= 50 {
		err.ValidMaxLenName = true
	}

	if phoneRegexp.MatchString(f.Phone) {
		err.ValidPhone = true
	}

	return err
}
