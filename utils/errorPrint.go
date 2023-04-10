package utils

import (
	"fmt"
	"log"
	"runtime"
)

func Errorf(any interface{}, a ...interface{}) error {
	if any != nil {
		err := (error)(nil)

		switch v := any.(type) {
		case string:
			err = fmt.Errorf(v, a...)
		case error:
			err = fmt.Errorf(v.Error(), a...)
		default:
			err = fmt.Errorf("%v", err)
		}

		_, fn, line, _ := runtime.Caller(1)
		log.Printf("Error: [%s:%d] %v \n", fn, line, err)

		return err
	}

	return nil
}
