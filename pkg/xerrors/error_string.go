package xerrors

import (
	"errors"
	"fmt"
)

func GetNotExist(key string) error {
	return errors.New(fmt.Sprintf("%s not exist. Please check again", key))
}

func GetMin(key string, min int) error {
	return errors.New(fmt.Sprintf("%s min value is %v. Please check again", key, min))
}

func GetRegexp(key string, pattern string) error {
	return errors.New(fmt.Sprintf("%s must match pattern %s. Please check again", key, pattern))
}

func GetEmpty(key string) error {
	return errors.New(fmt.Sprintf("%s is required to not be empty. Please check again", key))
}

func GetRequired(key string) error {
	return errors.New(fmt.Sprintf("%s is required. Please check again", key))
}

func GetEnum(key string, enums []string) error {
	return errors.New(fmt.Sprintf("%s value must be in the enum %v. Please check again", key, enums))
}

func GetRange(key string, min, max int32) error {
	return errors.New(fmt.Sprintf("%s value must be in the range %d - %d. Please check again", key, min, max))
}

func GetLen(key string, min, max int32) error {
	return errors.New(fmt.Sprintf("%s length must be in the range %d - %d. Please check again", key, min, max))
}

func GetComparseTime(key string) error {
	return errors.New(fmt.Sprintf("%s start time must be before end time. Please check again", key))
}

func GetDate(key string) error {
	return errors.New(fmt.Sprintf("%s invalid datetime. Please check again", key))
}

func GetEmail(key string) error {
	return errors.New(fmt.Sprintf("%s invalid email. Please check again", key))
}

func GetExist(key string) error {
	return errors.New(fmt.Sprintf("%s already exists. Please check again", key))
}
