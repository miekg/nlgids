package date

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	gnuDate = "date"
	langNL  = "nl_NL.UTF-8"
)

// NL the GNU data program with LANG=nl_NL.UTF-8 so that
// it will return a date in the correct format/language for the form.
func NL(opts ...string) (string, error) {
	// Yes, I can also do time.Time something something and then
	// translate the names to Dutch with a map or something.
	os.Setenv("LANG", langNL)

	cmd := exec.Command(gnuDate, opts...)
	cmdOut, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	out, _ := ioutil.ReadAll(cmdOut)
	out = bytes.TrimSuffix(out, []byte("\n"))

	err := cmd.Wait()
	return string(out), err
}
