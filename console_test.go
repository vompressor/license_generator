package license_generator_test

import (
	"testing"

	"github.com/vompressor/license_generator"
)

func TestPrintLicenseList(t *testing.T) {
	license_generator.PrintLicenseList()
}

func TestPrintLicenseInfo(t *testing.T) {
	license_generator.PrintLicenseInfo("Apache-2.0")
}

func TestPrintLicenseBody(t *testing.T) {
	license_generator.PrintLicenseBody("mit")
}

func TestWriteLicenseBodyToPath(t *testing.T) {
	err := license_generator.WriteLicenseBodyToPath("mit", "LICENSE", "2021", "owner")
	if err != nil {
		t.Fatal(err.Error())
	}
}
