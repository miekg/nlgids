package nlgids

import (
	"bufio"
	"net/http"
	"path"

	"github.com/miekg/nlgids/webinvoice"
)

const templateDir = "/opt/tmpl/nlgids"

func Test(w http.ResponseWriter, r *http.Request) {
	testInvoice := &webinvoice.Invoice{
		Tour:     "Van Koninklijke Huize",
		Persons:  2,
		Time:     "11.00",
		Duration: 2.0,
		Cost:     50.0,
		Date:     "2015/12/10",
		Name:     "Miek",
		FullName: "Miek Gieben",
		Email:    "miek@miek.nl",
		Where:    "Green Park metro station",
		How:      "Ik sta buiten de de fontein om",
	}

	tmpl := path.Join(templateDir, webinvoice.DefaultTemplate)

	pdf, err := testInvoice.Create(templateDir, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(pdf) == 0 {
		http.Error(w, "no pdf produced", http.StatusInternalServerError)
		return
	}
	rd := bufio.NewReader(pdf)
	Download(rd, i.FileName, w, r)
	return http.StatusOK, nil
}
