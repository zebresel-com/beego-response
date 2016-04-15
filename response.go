package beego_response

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego/context"
)

const (
	PARAM_SKIP  string = "skip"
	PARAM_LIMIT string = "limit"
)

type Response struct {
	data    map[string]map[string]interface{} `json:"data,omitempty"`
	context *context.Context                  `json:"-"`
}

// New initializes an empty response object
func New(context *context.Context) *Response {

	dataMap := make(map[string]map[string]interface{})

	dataMap["data"] = make(map[string]interface{})

	response := &Response{
		data:    dataMap,
		context: context,
	}

	return response
}

// AddContent creates a new map entry to wrap results. E.g. "users" : {result..}.
func (self *Response) AddContent(objectType string, data interface{}) {

	// Reserved entrys cant be overwritten
	if objectType == "error" || objectType == "pagination" {
		return
	}

	// Return if error already was set
	if _, ok := self.data["data"]["error"]; ok {
		return
	}

	self.data["data"][objectType] = data
}

// Error sets the specified status code, obtains a default message and directly writes the output (ServeJSON gets called).
// It is necessary that you call return or prevent any other output after this function!
// UserInfo can be specified additionally.
func (self *Response) Error(code int, userInfo ...interface{}) {

	if len(userInfo) > 0 {
		self.SetError(code, userInfo[0])
	} else {
		self.SetError(code)
	}

	self.ServeJSON()
}

// CustomError sets the specified status code with a custom message and directly writes the output (ServeJSON gets called).
// It is necessary that you call return or prevent any other output after this function!
// UserInfo can be specified additionally.
func (self *Response) CustomError(code int, customCode int, message string, userInfo ...interface{}) {

	if len(userInfo) > 0 {
		self.SetCustomError(code, customCode, message, userInfo[0])
	} else {
		self.SetCustomError(code, customCode, message)
	}

	self.ServeJSON()
}

// SetError clears all response data and only leaves the error object.
// Response object must be manually sent to the output.
// UserInfo can be specified additionally.
func (self *Response) SetError(code int, userInfo ...interface{}) {

	if len(userInfo) > 0 {
		self.SetCustomError(code, 0, http.StatusText(code), userInfo[0])
	} else {
		self.SetCustomError(code, 0, http.StatusText(code))
	}
}

// SetCustomError clears all response data and only leaves the error object.
// Additionally you can pass a custom message and code
// Response object must be manually sent to the output.
// UserInfo can be specified additionally.
func (self *Response) SetCustomError(code int, customCode int, message string, userInfo ...interface{}) {

	// Clear everything from payload except errors
	for key, _ := range self.data["data"] {

		if key != "error" {
			delete(self.data["data"], key)
		}
	}

	if customCode == 0 {
		customCode = code
	}

	// Handle user information
	userInfoAll := make([]*UserInfo, 0, 0)

	if len(userInfo) > 0 {

		// Check for single error
		if userError, ok := userInfo[0].(error); ok {

			info := &UserInfo{
				Message: userError.Error(),
			}

			userInfoAll = append(userInfoAll, info)

			// Check for string
		} else if userError, ok := userInfo[0].(string); ok {

			info := &UserInfo{
				Message: userError,
			}

			userInfoAll = append(userInfoAll, info)

			// Check for multiple errors
		} else if userErrors, ok := userInfo[0].([]error); ok {

			for _, userError := range userErrors {

				info := &UserInfo{
					Message: userError.Error(),
				}

				userInfoAll = append(userInfoAll, info)
			}

		} else {
			panic("userInfo for response must be error or []error")
		}
	}

	// Check if error object already exists
	if _, ok := self.data["data"]["error"]; !ok {
		self.data["data"]["error"] = &Error{
			Code:     customCode,
			Message:  message,
			UserInfo: userInfoAll,
		}
	} else {

		err := self.data["data"]["error"].(*Error)

		err.Code = customCode
		err.Message = message
		err.UserInfo = userInfoAll
	}

	self.context.Output.SetStatus(code)
}

// CreatePaging automatically generates a paging for the response instance.
// The skip param is the offset to start with and limit is the amount to show.
// Records is the document count for a query. currentRecords is the amount which is on the current page.
func (self *Response) CreatePaging(skip int, limit int, records int, currentRecords int) {

	if _, ok := self.data["data"]["error"]; ok {
		return
	}

	paging := &Paging{}

	var previousSkip, previousLimit, nextSkip, nextLimit int
	var hasPrevious, hasNext bool = true, true

	if skip == 0 {

		hasPrevious = false

	} else if skip > records {

		previousSkip = records - limit
		previousLimit = limit

		if previousSkip < 0 {
			previousSkip = 0
		}

	} else if limit > skip {

		previousSkip = 0
		previousLimit = skip

	} else {

		previousSkip = skip - limit
		previousLimit = limit

	}

	if skip+limit >= records {

		hasNext = false

	} else {

		nextSkip = skip + limit
		nextLimit = limit
	}

	requestURI := self.context.Request.RequestURI

	nextUrl := ""
	previousUrl := ""

	// Check if url has params to keep old ones
	if strings.Contains(requestURI, "?") {

		limitReg := regexp.MustCompile(PARAM_LIMIT + "=[0-9]*")
		skipReg := regexp.MustCompile(PARAM_SKIP + "=[0-9]*")

		// Limit exists
		if strings.Contains(requestURI, PARAM_LIMIT+"=") {

			nextUrl = limitReg.ReplaceAllString(requestURI, PARAM_LIMIT+"="+strconv.Itoa(nextLimit))
			previousUrl = limitReg.ReplaceAllString(requestURI, PARAM_LIMIT+"="+strconv.Itoa(previousLimit))

			paging.First = limitReg.ReplaceAllString(requestURI, PARAM_LIMIT+"="+strconv.Itoa(limit))
			paging.Last = limitReg.ReplaceAllString(requestURI, PARAM_LIMIT+"="+strconv.Itoa(limit))

		} else { // Add limit param manually

			nextUrl = fmt.Sprintf("%s&%s=%s", requestURI, PARAM_LIMIT, strconv.Itoa(nextLimit))
			previousUrl = fmt.Sprintf("%s&%s=%s", requestURI, PARAM_LIMIT, strconv.Itoa(previousLimit))

			paging.First = fmt.Sprintf("%s&%s=%s", requestURI, PARAM_LIMIT, strconv.Itoa(limit))
			paging.Last = fmt.Sprintf("%s&%s=%s", requestURI, PARAM_LIMIT, strconv.Itoa(limit))
		}

		// Skip exists
		if strings.Contains(requestURI, PARAM_SKIP+"=") {

			fmt.Println(previousSkip)

			nextUrl = skipReg.ReplaceAllString(nextUrl, PARAM_SKIP+"="+strconv.Itoa(nextSkip))
			previousUrl = skipReg.ReplaceAllString(previousUrl, PARAM_SKIP+"="+strconv.Itoa(previousSkip))

			paging.First = skipReg.ReplaceAllString(paging.First, PARAM_SKIP+"="+strconv.Itoa(0))
			paging.Last = skipReg.ReplaceAllString(paging.Last, PARAM_SKIP+"="+strconv.Itoa(previousSkip))

		} else { // Add skip param manually

			nextUrl = fmt.Sprintf("%s&%s=%s", nextUrl, PARAM_SKIP, strconv.Itoa(nextSkip))
			previousUrl = fmt.Sprintf("%s&%s=%s", previousUrl, PARAM_SKIP, strconv.Itoa(previousSkip))

			paging.First = fmt.Sprintf("%s&%s=%s", paging.First, PARAM_SKIP, strconv.Itoa(0))
			paging.Last = fmt.Sprintf("%s&%s=%s", paging.Last, PARAM_SKIP, strconv.Itoa(records-limit))
		}

	} else { // Otherwise build custom url without replacing

		previousUrl = fmt.Sprintf("%s?%s=%s&%s=%s", self.context.Request.RequestURI, PARAM_SKIP, strconv.Itoa(previousSkip), PARAM_LIMIT, strconv.Itoa(previousLimit))
		nextUrl = fmt.Sprintf("%s?%s=%s&%s=%s", self.context.Request.RequestURI, PARAM_SKIP, strconv.Itoa(nextSkip), PARAM_LIMIT, strconv.Itoa(nextLimit))

		paging.First = fmt.Sprintf("%s?%s=%s&%s=%s", self.context.Request.RequestURI, PARAM_SKIP, strconv.Itoa(0), PARAM_LIMIT, strconv.Itoa(limit))
		paging.Last = fmt.Sprintf("%s?%s=%s&%s=%s", self.context.Request.RequestURI, PARAM_SKIP, strconv.Itoa(records-limit), PARAM_LIMIT, strconv.Itoa(limit))
	}

	fmt.Println(skip)

	if hasPrevious {
		paging.Previous = previousUrl
	}

	if hasNext {
		paging.Next = nextUrl
	}

	paging.RecordsTotal = records
	paging.RecordsPage = currentRecords
	paging.Pages = int(math.Ceil(float64(records) / float64(limit)))

	if skip >= records {
		paging.CurrentPage = 0
	} else {
		paging.CurrentPage = paging.Pages - int(math.Ceil(float64(records-skip)/float64(limit))) + 1
	}

	if paging.Pages == paging.CurrentPage {
		paging.Last = self.context.Request.RequestURI
	}

	self.data["data"]["pagination"] = paging
}

// Set custom status code
func (self *Response) SetStatus(code int) {
	self.context.Output.SetStatus(code)
}

// Return only the payload as map (normally not needed)
func (self *Response) Data() *map[string]interface{} {

	dataMap := self.data["data"]

	return &dataMap
}

// Use this method to output the response data (instead of beego method)
func (self *Response) ServeJSON() {

	self.context.Output.Json(self.data["data"], true, false)
}
