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

func TypeLog() *typeLog {
	return &typeLog{
		Error:         "[Error]",
		Debug:         "[Debug]",
		Info:          "[Info]",
		DatabaseError: "[Database Error]",
		DatabaseInfo:  "[Database Info]",
		Svc_Auth_Info: "[Auth Svc Info]",
	}
}

type typeLog struct {
	Error         string
	Debug         string
	Info          string
	DatabaseError string
	DatabaseInfo  string
	Svc_Auth_Info string
}
