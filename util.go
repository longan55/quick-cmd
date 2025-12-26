package main

import (
	"fmt"
	"log"
)

func SQL_Slice_To_In_Args(slice any) (string, error) {
	log.Printf("SQL_Slice_To_In_Args: %+v", slice)
	if slice == nil {
		return "()", nil
	}
	args := "("
	switch sli := slice.(type) {
	case []string:
		for _, v := range sli {
			if args != "(" {
				args += ","
			}
			args += fmt.Sprintf("'%s'", v)
		}
		args += ")"
		return args, nil
	case []uint64:
		for _, v := range sli {
			if args != "(" {
				args += ","
			}
			args += fmt.Sprintf("%d", v)
		}
		args += ")"
		return args, nil
	case []int:
		for _, v := range sli {
			args += fmt.Sprintf("%d,", v)
		}
		args += ")"
		return args, nil
	default:
		return "()", fmt.Errorf("slice type %T not supported", slice)
	}
}
