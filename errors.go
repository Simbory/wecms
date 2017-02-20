package wecms

import "fmt"

func assertNotEmpty(data string, paramName string) {
	if len(data) == 0 {
		panic(errParamEmpty(paramName))
	}
}

func assertNotNil(data interface{}, paramName string) {
	if data == nil {
		panic(errParamNil(paramName))
	}
}

func errParamEmpty(paramName string) error {
	return fmt.Errorf("The parameter '%s' cannot be empty", paramName)
}

func errParamNil(paramName string) error {
	return fmt.Errorf("The parameter '%s' cannot be nil", paramName)
}

func errSessionNil(dbName string) error {
	return fmt.Errorf("the data session of this repository is nil. Database: %s", dbName)
}