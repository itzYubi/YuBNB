package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIsValid(t *testing.T) {
	req := httptest.NewRequest("POST", "/fake", nil)
	form := New(req.PostForm)
	if !form.IsValid() {
		t.Error("Should have said that form is valid, but says it is invalid")
	}
}

func TestFormRequired(t *testing.T) {
	req := httptest.NewRequest("POST", "/fake", nil)
	form := New(req.PostForm)
	form.Required("a", "b", "c")
	if form.IsValid() {
		t.Error("form shows valid when it should be invalid")
	}
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")
	r := httptest.NewRequest("POST", "/fake", nil)
	r.PostForm = postData

	new_form := New(r.PostForm)
	new_form.Required("a", "b", "c")
	if !new_form.IsValid() {
		t.Error("form shows invalid when it should be valid")
	}
}

func TestHas(t *testing.T) {
	req := httptest.NewRequest("POST", "/fake", nil)
	form := New(req.PostForm)
	if form.Has("a") {
		t.Error("form shows it has field: a, but it does not")
	}

	postData := url.Values{}
	postData.Add("a", "bb")
	r, _ := http.NewRequest("POST", "/fake", nil)
	r.PostForm = postData
	form2 := New(r.PostForm)
	if !form2.Has("a") {
		t.Error("form shows it does not have field: a, but it does")
	}
}

func TestMinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("a", "bb")
	r := httptest.NewRequest("POST", "/fake", nil)
	r.PostForm = postData
	form := New(r.PostForm)
	if form.MinLength("a", 3) {
		t.Error("form shows that it has correct min length but it does not have")
	}

	if form.Errors.Get("a") == "" {
		t.Error("Should have error but shows does not have")
	}

	postData = url.Values{}
	postData.Add("a", "bbbb")
	r = httptest.NewRequest("POST", "/fake", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	if !form.MinLength("a", 3) {
		t.Error("form shows that it does not have correct min length but it does")
	}

	if form.Errors.Get("a") != "" {
		t.Error("Should not have error but shows has error")
	}
}

func TestEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/fake", nil)
	form := New(r.PostForm)
	form.IsEmail("x")
	if form.IsValid() {
		t.Error("Form shows Email valid for a non existent field")
	}

	postData := url.Values{}
	postData.Add("email", "me@here.com")
	r = httptest.NewRequest("POST", "/fake", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.IsEmail("email")
	if !form.IsValid() {
		t.Error("Form shows Email invalid for a vaid email")
	}

	postData = url.Values{}
	postData.Add("email", "com")
	r = httptest.NewRequest("POST", "/fake", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.IsEmail("email")
	if form.IsValid() {
		t.Error("Form shows Email valid for an invaid email")
	}
}
