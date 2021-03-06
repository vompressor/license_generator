package license_generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"time"
)

const githubLicenseAPI = "https://api.github.com/licenses"

// TODO::
const vompressorLicenseInfoAPI = ""

var useCache = true

//

// LicenseKey is Struct of license key, name, and URL to get detail information.
type LicenseKey struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

// License is Struct of Contains detailed information about the license.
type License struct {
	Key            string
	Name           string
	SpdxID         string `json:"spdx_id"`
	URL            string
	LicenseURL     string `json:"html_url"`
	Description    string
	Implementation string
	Permissions    []string
	Conditions     []string
	Limitations    []string
	Body           string
	Featured       bool
}

type cacheData struct {
	Created int64
	Expire  int64
	TTL     int64
}

// cache file config
const cacheFileHead = "lfm cache data"
const cacheTTL = time.Hour * 24 * 3
const cacheDir = "lfm/"

func cacheThis(name string, item interface{}, ttl time.Duration) error {
	p, _ := os.UserCacheDir()
	cachePath := filepath.Join(p, cacheDir)

	var cd cacheData

	cd.TTL = int64(ttl)
	cd.Created = time.Now().Unix()
	cd.Expire = time.Now().Add(ttl).Unix()

	os.MkdirAll(cachePath, os.ModePerm)

	f, err := os.OpenFile(filepath.Join(cachePath, name), os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	writer.WriteString(cacheFileHead + "\n")

	jenc := json.NewEncoder(writer)
	jenc.Encode(cd)
	jenc.Encode(item)

	return nil
}

// TODO:: error impl

func readCache(name string, d interface{}) error {
	var cd cacheData
	p, _ := os.UserCacheDir()
	cachePath := filepath.Join(p, cacheDir, name)
	f, err := os.Open(cachePath)

	if err != nil {
		return err
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	data, ispfx, err := reader.ReadLine()

	if err != nil {
		return err
	}

	if ispfx {
		var e WrongCacheError
		return e
	}

	if string(data) != cacheFileHead {
		var e WrongCacheError
		return e
	}

	dec := json.NewDecoder(reader)

	err = dec.Decode(&cd)
	if err != nil {
		return err
	}

	if time.Now().Unix() > cd.Expire {
		var e TTLExpireError
		e.ExpireAt = time.Unix(cd.Expire, 0)
		e.TTL = time.Duration(cd.TTL)
		return e
	}

	err = dec.Decode(d)
	if err != nil {
		return err
	}

	return nil
}

func DelCache() error {
	p, _ := os.UserCacheDir()
	cachePath := filepath.Join(p, cacheDir)

	err := os.RemoveAll(cachePath)
	if err != nil {
		return err
	}
	return nil
}

func GetLicenseKeys() ([]LicenseKey, error) {

	var licenses []LicenseKey

	err := readCache("list", &licenses)

	if err == nil {
		return licenses, nil
	}

	req, err := http.NewRequest("GET", githubLicenseAPI, nil)
	if err != nil {
		return nil, ServerError{1, err}
	}

	req.Header.Add("accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ServerError{1, err}
	}

	if resp.StatusCode != 200 {
		return nil, HttpCodeError{HttpCode: resp.StatusCode}
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&licenses)
	if err != nil {
		return nil, err
	}

	cacheThis("list", licenses, cacheTTL)

	return licenses, nil
}

func GetLicenseInfo(license string) (*License, error) {

	var ret License

	err := readCache(license, &ret)

	if err == nil {
		return &ret, nil
	}

	req, err := http.NewRequest("GET", githubLicenseAPI+"/"+license, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ServerError{1, err}
	}

	if resp.StatusCode != 200 {
		return nil, HttpCodeError{HttpCode: resp.StatusCode}
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}

	cacheThis(license, ret, cacheTTL)

	return &ret, nil
}

type TTLExpireError struct {
	TTL      time.Duration
	ExpireAt time.Time
}

func (e TTLExpireError) Error() string {
	return fmt.Sprintf("%f over...", time.Now().Sub(e.ExpireAt).Seconds())
}

type WrongCacheError struct {
}

func (e WrongCacheError) Error() string {
	return "Invalid cache file"
}

type HttpCodeError struct {
	HttpCode int
}

func (e HttpCodeError) Error() string {
	return http.StatusText(e.HttpCode)
}

type ServerError struct {
	Code      int
	SourceErr error
}

func (e ServerError) Error() string {
	return ecode[e.Code]
}

var ecode = map[int]string{
	1: "server error",
}
