package license_generator_test

import (
	"testing"

	"github.com/vompressor/license_generator"
)

func TestGetLicenseKeys(t *testing.T) {
	ret, err := license_generator.GetLicenseKeys()

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("%#v", ret)

}

func TestGetLicenseInfo(t *testing.T) {
	ret, err := license_generator.GetLicenseInfo("apache-2.0")

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("%#v", ret)
}
