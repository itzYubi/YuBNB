package render

import (
	"net/http"
	"os"
	"testing"

	"github.com/itzYubi/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error("Failed to get session")
	}
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Failed to add default data as flash value of 123 was not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var w myWriter

	err = Template(&w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}

	err = Template(&w, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Error writing template to browser, should be an error")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/url", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(&testApp)

}

func TestCreateTemplateCache(t *testing.T) {
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error("Error creating Template Cache")
	}
	files, _ := os.ReadDir("./../../templates")

	if len(tc) != len(files) && tc == nil {
		t.Errorf("Did not render correctly got %d, expected %d", len(tc), len(files))
	}
}
