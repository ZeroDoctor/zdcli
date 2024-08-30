package generate

import (
	"fmt"
	"os"
	"testing"

	"github.com/zerodoctor/zdcli/util"
)

func TestGenerateTOTP(t *testing.T) {
	util.EXEC_PATH = ".."

	err := TOTP("test", "barcode")
	if err != nil {
		fmt.Printf("generate [error=%s]\n", err.Error())
	}

	err = os.Remove("./totp_qr.html")
	if err != nil {
		fmt.Printf("remove [error=%s]\n", err.Error())
	}
}
