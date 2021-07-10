package license_generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PrintLicenseList is Prints the license list received from the github api.
func PrintLicenseList() error {
	keys, err := GetLicenseKeys()

	if err != nil {
		return err
	}

	println("key : id : name")
	for _, key := range keys {
		fmt.Printf("%s : %s : %s\n", key.Key, key.SpdxID, key.Name)
	}
	return nil
}

// PrintLicenseBody is Prints the license body received from the github api.
func PrintLicenseBody(key string) error {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		return err
	}

	println(l.Body)

	return nil
}

// PrintLicenseInfo is Prints general information about the license.
func PrintLicenseInfo(key string) error {
	l, err := GetLicenseInfo(strings.ToLower(key))

	// TODO:: err type
	if err != nil {
		return err
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

	print("\nlimitations:")
	for _, c := range l.Limitations {
		print(" " + c)
	}
	print("\n")
	return nil
}

// WriteLicenseBody is Write license body to input io.Writer
func WriteLicenseBody(key string, w io.Writer, year string, owner string) error {
	l, err := GetLicenseInfo(strings.ToLower(key))
	ret := l.Body

	// TODO:: err type
	if err != nil {
		return err
	}
	if year != "" {
		ret = strings.Replace(ret, "[year]", year, 1)
	}
	if owner != "" {
		ret = strings.Replace(ret, "[owner]", owner, 1)
		ret = strings.Replace(ret, "[fullname]", owner, 1)
	}

	_, err = io.WriteString(w, ret)
	if err != nil {
		return err
	}

	return nil
}

// WriteLicenseBody is Write license body to input path
func WriteLicenseBodyToPath(key string, path string, year string, owner string) error {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	f, err := os.Create(abs_path)
	if err != nil {
		return err
	}
	defer f.Close()

	WriteLicenseBody(key, f, year, owner)

	return nil
}
