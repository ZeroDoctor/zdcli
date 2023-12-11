package generate

import (
	"bytes"
	"html/template"
	"os"

	"github.com/zerodoctor/zdcli/util"
)

func TOTP(barcode string) error {
	tpl, err := template.New("qr_totp_template.html").
		ParseFiles(util.EXEC_PATH + "/assets/qr_totp_template.html")
	if err != nil {
		return err
	}

	data := struct {
		Barcode string
	}{
		Barcode: barcode,
	}

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, data)
	if err != nil {
		return err
	}

	return os.WriteFile("totp_qr.html", buffer.Bytes(), 0777)
}
