package util

import "os"

func IsUat() bool {
	return os.Getenv("DEPLOY_ENV") == "uat"
}

func IsPre() bool {
	return os.Getenv("DEPLOY_ENV") == "pre"
}

func IsProd() bool {
	return os.Getenv("DEPLOY_ENV") == "prod"
}
