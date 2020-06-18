package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdenticonHandler(t *testing.T) {
	//This function uses the methods and structs in the
	//net/http/httptest package to build the
	//http.ResponseWriter and *http.Request that
	//will invoke the handler function,
	//and examine what was written to the response.

	// Add more test cases for other Headers and Status Code

	cases := []struct {
		name                string
		query               string
		expectedStatusCode  int
		expectedContentType string
	}{
		{
			"Valid Name Param",
			"name=test",
			http.StatusOK,
			contentTypePNG,
		},
		{
			"No name given",
			"name=",
			400,
			contentTypeError,
		},
		{
			"No name given COMPLETELY BLANK",
			"",
			400,
			contentTypeError,
		},
		{
			"Wrong URL Given",
			"/test",
			400,
			contentTypeError,
		},
	}

	for _, c := range cases {
		URL := fmt.Sprintf("/identicon?%s", c.query)
		req := httptest.NewRequest("GET", URL, nil)
		respRec := httptest.NewRecorder()
		IdenticonHandler(respRec, req)

		resp := respRec.Result()
		//check the response status code
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		} else if resp.Header.Get(headerContentType) != c.expectedContentType { //check the response header content type
			t.Errorf("case %s: incorrect header content type: expected %s but got %s",
				c.name, c.expectedContentType, resp.Header.Get(headerContentType))
		}

	}

}
