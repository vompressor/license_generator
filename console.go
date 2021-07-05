package license_generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PrintLicenseList() {
	keys, err := GetLicenseKeys()

	if err != nil {
		println("GitHub License api not work...")
		os.Exit(1)
	}

	println("key : id : name")
	for _, key := range keys {
		fmt.Printf("%s : %s : %s\n", key.Key, key.SpdxID, key.Name)
	}
}

func PrintLicenseBody(key string) {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		println("Invalid input...")
		os.Exit(1)
	}

	println(l.Body)
}

func PrintLicenseInfo(key string) {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		println("Invalid input...")
		os.Exit(1)
	}

	fmt.Printf("key: %s\n", l.Key)
	fmt.Printf("id:  %s\n", l.SpdxID)
	fmt.Printf("url: %s\n", l.LicenseURL)

	fmt.Printf("\ndescription:\n%s\n", l.Description)

	fmt.Printf("\nimplementation:\n%s\n", l.Implementation)

	print("\npermissions:")
	for _, c := range l.Permissions {
		print(" " + c)
	}

	print("\nconditions:")
	for _, c := range l.Conditions {
		print(" " + c)
	}

	print("\n limitations:")
	for _, c := range l.Limitations {
		print(" " + c)
	}
	print("\n")

}

func WriteLicenseBody(key string, w io.Writer) error {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, l.Body)
	if err != nil {
		return err
	}

	return nil
}

func WriteLicenseBodyToPath(key string, path string) error {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		return err
	}

	abs_path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	err = os.WriteFile(abs_path, []byte(l.Body), 0600)
	if err != nil {
		return err
	}

	return nil
}