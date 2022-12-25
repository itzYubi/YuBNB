package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"natsu home", "/natsu-home", "GET", []postData{}, http.StatusOK},
	{"yuki home", "/yuki-home", "GET", []postData{}, http.StatusOK},
	{"Search Availablity page", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make reservation page", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"contact page", "/contact", "GET", []postData{}, http.StatusOK},
	{"post-search-availability page", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2021-01-03"},
	}, http.StatusOK},
	{"post-search-availability-json page", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2021-01-03"},
	}, http.StatusOK},
	{"post-make-reservation page", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Yubi"},
		{key: "last_name", value: "Uchicha"},
		{key: "email", value: "yubi.uchiha@gmail.com"},
		{key: "phone", value: "2422123242"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			response, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			response, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
