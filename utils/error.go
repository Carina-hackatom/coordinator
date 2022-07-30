package utils

import (
	"fmt"
	"github.com/Carina-hackatom/coordinator/utils/types"
	"log"
	"os"
)

func CheckErr(err error, moreMsg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %v", moreMsg, err))

		}
	case types.KEEP:
		if err != nil {
			log.Printf("%s: \n %v\n", moreMsg, err)
		}
	}
}

func LogErrWithFd(fd *os.File, err error, msg string, action types.Code) {
	switch action {
	case types.EXIT:
		if err != nil {
			panic(fmt.Sprintf("%s: \n %v", msg, err))
		}
	case types.KEEP:
		if err != nil {
			l := log.New(fd, "ERROR (check) : ", log.Llongfile|log.LstdFlags)
			l.Printf("%s: \n %v\n", msg, err)
		}
	}
}
