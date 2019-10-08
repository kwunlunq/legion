package glob

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Errors []error

func (errs Errors) Add(newErrs ...error) Errors {
	for _, err := range newErrs {
		if err == nil {
			continue
		}
		errs = append(errs, err)
	}
	return errs
}

func (errs Errors) Error() string {
	var errors = []string{}
	for _, err := range errs {
		if err == nil {
			continue
		}
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, "; ")
}

func (errs Errors) IsNil() bool {
	for _, err := range errs {
		if err != nil {
			return false
		}
	}
	return true
}

func (errs Errors) MarshalJSON() (buf []byte, err error) {
	// if errs == nil {
	// 	return []byte(`null`), nil
	// }

	var strs []string
	for _, err := range errs {
		if err == nil {
			continue
		}
		strs = append(strs, err.Error())
	}

	return json.Marshal(strs)
}

func (errs *Errors) UnmarshalJSON(data []byte) (err error) {
	// if string(data) == "null" {
	// 	return nil
	// }

	var strPtrs []*string
	err = json.Unmarshal(data, &strPtrs)
	if err != nil {
		return
	}

	for _, strPtr := range strPtrs {
		if strPtr == nil {
			continue
		}
		*errs = append(*errs, fmt.Errorf("%s", *strPtr))
	}

	return
}
