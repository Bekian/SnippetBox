package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
)

// 500 range errors
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// 400 range errors
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// write template to buffer to catch error prematurely
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// template Data initializer to populate the current year
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// s/n theres gotta be a better way to write this
	// parse the form
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// attempt to decode the form into the destination struct
	err = app.formDecorder.Decode(dst, r.PostForm)
	if err != nil {
		// here we use this variable to check if the dst is invalid
		var invalidDecoderErr *form.InvalidDecoderError
		// if it is, panic
		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}

		// else return
		return err
	}

	return nil
}
