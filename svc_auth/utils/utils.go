package utilsApi

import (
	"log"
)

type error interface {
	Error() string
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DefaultEnvField(env string, defaul string) string {
	res := env
	if env == "" {
		res = defaul
	}
	return res
}
