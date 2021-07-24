package main

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type API struct {
	USERAGENT    string
	VERSION      string
	KeyVersion   string
	KEY          string
	CAPABILITIES string
}

type HttpResponse struct {
	Err                 error
	ResStatus           int
	Req                 *http.Request
	Res                 *http.Response
	Body                string
	Headers             http.Header
	Cookies             []*http.Cookie
	RequestSizeByBytes  float64
	ResponseSizeByBytes float64
}

type Secret struct {
	version string
	key     string
}

func randDevice(version string) string {

	dpi := []string{
		"480", "320", "640", "515", "120", "160", "240", "800",
	}
	manufacturers := []string{
		"HUAWEI", "Xiaomi", "samsung", "OnePlus", "LGE/lge", "ZTE", "HTC",
		"LENOVO", "MOTOROLA", "NOKIA", "OPPO", "SONY", "VIVO", "LAVA",
	}

	randResolution := randomdata.Number(2, 9) * 180
	lowerResolution := randResolution - 180

	androidVersion := randomdata.Number(18, 25)
	androidRelease := fmt.Sprintf("%d.%d", randomdata.Number(1, 7), randomdata.Number(0, 7))
	if randomdata.Boolean() {
		androidRelease = fmt.Sprintf("%d.%d.%d", randomdata.Number(1, 7), randomdata.Number(0, 7), randomdata.Number(1, 7))
	}
	_dpi := dpi[randomdata.Number(0, len(dpi))]
	resolution := fmt.Sprintf("%dx%d", lowerResolution, randResolution)
	manufacturer := manufacturers[randomdata.Number(0, len(manufacturers))]
	device := fmt.Sprintf("%s-%s", manufacturers[randomdata.Number(0, len(manufacturers))], randomdata.RandStringRunes(5))
	model := randomdata.RandStringRunes(4)
	cpu := fmt.Sprintf("%s%d", randomdata.RandStringRunes(2), randomdata.Number(1000, 9999))

	UserAgentBase := "Instagram %s Android (%d/%s; %s; %s; %s; %s; %s; %s; en_US)"
	return fmt.Sprintf(UserAgentBase, version, androidVersion, androidRelease, _dpi, resolution, manufacturer, device, model, cpu)

}

func newDeviceID() string {
	return "android-" + randomdata.RandStringRunes(16)
}

func GetAPI() API { // random choise

	var version string
	var key string

	version = fmt.Sprintf("%d.%d.%d", randomdata.Number(3, 138), randomdata.Number(5, 10), randomdata.Number(0, 10))
	if randomdata.Boolean() {
		version = fmt.Sprintf("%d.%d.%d", randomdata.Number(4, 138), randomdata.Number(0, 10), randomdata.Number(0, 10))
	}
	key = "SIGNATURE"

	USERAGENT := randDevice(version)
	IG_VERSION := version
	IG_SIG_KEY := key
	SIG_KEY_VERSION := "4"
	if randomdata.Boolean() {
		SIG_KEY_VERSION = "5"
	}
	X_IG_Capabilities := "3brTvw=="
	if randomdata.Boolean() {
		X_IG_Capabilities = fmt.Sprintf("%s==", randomdata.RandStringRunes(6))
	}

	_API := API{VERSION: IG_VERSION, KEY: IG_SIG_KEY, KeyVersion: SIG_KEY_VERSION, CAPABILITIES: X_IG_Capabilities, USERAGENT: USERAGENT}

	return _API
}

func MakeHttpResponse(Response *http.Response, Request *http.Request, Error error, RequestSizeByBytes float64, ResponseSizeByBytes float64) HttpResponse {

	var res = ""
	var StatusCode = 0
	var Headers http.Header = nil
	var cookies []*http.Cookie = nil
	var err error

	if Error != nil {
		err = Error
	}
	if Response != nil {

		if Response.StatusCode != 0 {
			StatusCode = Response.StatusCode
		}

		if Response.Header != nil {
			Headers = Response.Header
		}

		cookies = Response.Cookies()
		var reader io.ReadCloser
		switch Response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ = gzip.NewReader(Response.Body)
		default:
			reader = Response.Body
		}

		if reader != nil {

			//try start
			func() {
				defer func() {
					ex := recover() // except Exception as ex || Catch ex As Exception
					if ex != nil {
						//println(ex)
						//do something with ex
						//e1++
					}
				}()
				//try:

				body, err := ioutil.ReadAll(reader)
				if err == nil {
					defer reader.Close()
					res = string(body)
				}

			}()
			//try end

		}

	}

	return HttpResponse{ResStatus: StatusCode, Res: Response, ResponseSizeByBytes: ResponseSizeByBytes, Req: Request, RequestSizeByBytes: RequestSizeByBytes, Body: res, Headers: Headers, Cookies: cookies, Err: err}
}

func instRequest(iurl string, signedbody map[string]string, payload string,
	Headers map[string]string, api API, proxy string,
	cookie []*http.Cookie, usecookies bool, timeoutMilliseconds int) HttpResponse {

	_url := iurl

	if ((!strings.Contains(_url, "https")) || (!strings.Contains(_url, "http"))) && _url[0] != '/' {
		_url = "https://i.instagram.com/api/v1/" + _url
	} else if ((!strings.Contains(_url, "https")) || (!strings.Contains(_url, "http"))) && _url[0] == '/' {
		_url = "https://i.instagram.com/api/v1" + _url
	}

	_api := API{}
	if api == (API{}) {
		_api = GetAPI()
	} else {
		_api = api
	}

	_payload := ""
	if signedbody != nil {
		_data, _ := json.Marshal(signedbody)
		_json := string(_data)
		_signed := fmt.Sprintf("SIGNATURE.%s", _json)
		_payload = "ig_sig_key_version=" + _api.KeyVersion + "&signed_body=" + _signed
	} else if payload != "" {
		_payload = payload
	}

	var req *http.Request
	if _payload != "" {
		req, _ = http.NewRequest("POST", _url, bytes.NewBuffer([]byte(_payload)))
	} else {
		req, _ = http.NewRequest("GET", _url, nil)
	}

	req.Header.Set("User-Agent", "Instagram "+_api.VERSION+" Android (19/4.4.2; 480dpi; 1080x1920; samsung; SM-N900T; hltetmo; qcom; en_US)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("X-IG-Capabilities", _api.CAPABILITIES)
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "Keep-Alive")

	if Headers != nil {
		var keys []string
		for key := range Headers {
			keys = append(keys, key)
		}
		var values []string
		for _, value := range Headers {
			values = append(values, value)
		}

		for i := 0; i < len(keys); i++ {
			req.Header.Set(keys[i], values[i])
		}
	}

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(_url)
	jar.SetCookies(u, cookie)

	transport := http.Transport{}
	if proxy != "" {
		proxyUrl := &url.URL{Host: proxy}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	if timeoutMilliseconds != 0 {
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeoutMilliseconds)}
	}
	if usecookies {
		if timeoutMilliseconds != 0 {
			client = &http.Client{Timeout: time.Millisecond * time.Duration(timeoutMilliseconds), Jar: jar}
		} else {
			client = &http.Client{Jar: jar}
		}
	}

	client.Transport = &transport
	resp, err := client.Do(req)

	RawReq, _ := httputil.DumpRequest(req, true)
	ReqSize := float64(len(RawReq))
	if _payload != "" {
		ReqSize += float64(len([]byte(_payload)))
		ReqSize += 4
	}

	if resp == nil {
		if err != nil {
			return MakeHttpResponse(nil, req, err, ReqSize, 0)
		}
		return MakeHttpResponse(nil, req, nil, ReqSize, 0)
	}
	RawRes, _ := httputil.DumpResponse(resp, true)
	ResSize := float64(len(RawRes))

	if err != nil {
		return MakeHttpResponse(resp, req, err, ReqSize, ResSize)
	}
	defer resp.Body.Close()
	return MakeHttpResponse(resp, req, nil, ReqSize, ResSize)
}

type account struct {
	username, password, sessionID, uuid, uid, device_id string
	api                                                 API
	loggedIn                                            bool
}

func newAccount(username, password string, proxy string, timeoutMilliseconds int) (*account, bool) {
	u, _ := uuid.NewUUID()
	guid := u.String()
	device_id := newDeviceID()
	Account := &account{
		username:  username,
		password:  password,
		sessionID: "",
		uid:       "",
		uuid:      guid,
		device_id: device_id,
		api:       GetAPI(),
		loggedIn:  false,
	}
	flag := Account.login(proxy, timeoutMilliseconds)
	return Account, flag
}

func (a *account) followers(userID, maxID string, proxy string, timeoutMilliseconds int) ([]string, []string, string, bool) {

	if !a.loggedIn && a.sessionID == "" {
		return nil, nil, "", false
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    a.sessionID,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie)

	path := fmt.Sprintf("friendships/%s/followers/?search_surface=follow_list_page&order=default&query=&enable_groups=true&rank_token=missing", userID)
	if maxID != "" {
		path += fmt.Sprintf("&max_id=%s", maxID)
	}

	errors := -1
retry:
	errors++
	if errors >= 1 {
		return nil, nil, "", false
	}

	res := instRequest(path, nil, "", nil, a.api, proxy, cookies, true, timeoutMilliseconds)

	if res.ResStatus == 200 {

		var usernames []string
		var userIDs []string

		usernames = extractRegex2(`"username":\s*"(.*?)"`, res.Body, 0)
		userIDs = extractRegex2(`"pk":\s*([0-9]*\s*?)`, res.Body, 14)
		next_max_id := gjson.Get(res.Body, "next_max_id").String()

		if next_max_id != "" {
			return usernames, userIDs, next_max_id, true
		} else {
			return usernames, userIDs, next_max_id, false
		}

	}

	goto retry

}

func (a *account) CheckSessionID(proxy string, timeoutMilliseconds int) bool {

	if !a.loggedIn && a.sessionID == "" {
		return false
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    a.sessionID,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie)

	errors := -1
retry:
	errors++
	if errors >= 1 {
		return false
	}

	res := instRequest("accounts/current_user/?edit=true", nil, "", nil, a.api, proxy, cookies, true, timeoutMilliseconds)

	if res.ResStatus == 200 {
		return true
	} else {
		goto retry
	}

}

func (a *account) login(proxy string, timeoutMilliseconds int) bool {

	url := "https://i.instagram.com/api/v1/accounts/login/"

	var Cookies []*http.Cookie

	post := make(map[string]string)
	post["phone_id"] = a.uuid
	post["_csrftoken"] = "missing"
	post["username"] = a.username
	post["password"] = a.password
	post["device_id"] = a.uuid
	post["guid"] = a.uuid
	post["login_attempt_count"] = "0"

	errors := -1
retry:
	errors++
	if errors >= 1 {
		return false
	}

	res := instRequest(url, post, "", nil, a.api, proxy, Cookies, true, timeoutMilliseconds)
	for _, cookie := range res.Cookies {
		if strings.ToLower(cookie.Name) == "sessionid" || strings.ToLower(cookie.Name) == "session_id" {
			a.sessionID = cookie.Value
			a.loggedIn = true
		}
		if strings.ToLower(cookie.Name) == "ds_user_id" || strings.ToLower(cookie.Name) == "uid" {
			a.uid = cookie.Value
		}
	}

	if a.sessionID != "" && a.uid != "" && a.loggedIn {
		return true
	}

	goto retry

}

func (a *account) changePassword(new string, proxy string, timeoutMilliseconds int) bool {

	if !a.loggedIn && a.sessionID == "" {
		return false
	}

	post := make(map[string]string)
	post["_uuid"] = a.uuid
	post["_uid"] = a.uid
	post["_csrftoken"] = "missing"
	post["old_password"] = a.password
	post["new_password1"] = new
	post["new_password2"] = new

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    a.sessionID,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie)

	errors := -1
retry:
	errors++
	if errors >= 1 {
		return false
	}

	res := instRequest("accounts/change_password/", post, "", nil, a.api, proxy, cookies, true, timeoutMilliseconds)
	if res.ResStatus == 200 {
		return true
	} else {
		goto retry
	}

	return false
}

func (a *account) IDByUsername(username string, proxy string, timeoutMilliseconds int) (string, error) {

	if !a.loggedIn && a.sessionID == "" {
		return "", nil
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    a.sessionID,
		Path:     "/",
		Domain:   ".instagram.com",
		Secure:   true,
		HttpOnly: true,
	}
	cookies = append(cookies, cookie)

	errors := -1
retry:
	errors++
	if errors >= 1 {
		return "", fmt.Errorf("An error with request")
	}

	res := instRequest("users/"+username+"/usernameinfo/", nil, "", nil, a.api, proxy, cookies, true, timeoutMilliseconds)

	if res.Err != nil || res.Body == "" || !strings.Contains(res.Body, "pk") {
		goto retry
	}

	results := extractRegex2(`"pk":\s*([0-9]*\s*?)`, res.Body, 0)
	if len(results) > 0 {
		return results[len(results)-1], nil
	} else {
		return "", fmt.Errorf("Pk is empty")
	}

}
