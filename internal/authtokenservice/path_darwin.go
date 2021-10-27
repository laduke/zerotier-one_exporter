package authtoken

import (
	"os"
)


func tokenPath () string {
	home, _ := os.UserHomeDir()

	return home + "/Library/Application Support/ZeroTier/One/authtoken.secret"

}
