package authtoken

import (
	"io/ioutil"
)

func Guess () string {
	b, err := ioutil.ReadFile(tokenPath())
	if err != nil { panic(err) }

	token := string(b)
	return token
}
