package authtoken

import (
	"os"
)


func TokenPath () string {
	home, _ := os.UserHomeDir()

	return home + "/Library/Application Support/ZeroTier/One/authtoken.secret"

}
