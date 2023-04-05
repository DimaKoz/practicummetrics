package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodCheckerWrongMethod(t *testing.T) {

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t.Error("This handler must not be called.")

	})

	// create the handler to test, using our custom "next" handler
	handlerToTest := MethodChecker(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)

	// call the handler using a mock response recorder
	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)
	res := w.Result()
	// check the result
	want := http.StatusMethodNotAllowed
	if res.StatusCode != want {
		t.Errorf("StatusCode got: %v, want: %v", res.StatusCode, want)
	}

	res.Body.Close()
}

func TestMethodCheckerOkMethod(t *testing.T) {
	isNextUsed := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		isNextUsed = true

	})

	handlerToTest := MethodChecker(nextHandler)
	req := httptest.NewRequest("POST", "http://testing", nil)
	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)
	res := w.Result()
	want := true
	if isNextUsed != want {
		t.Errorf("isNextUsed got: %v, want: %v", isNextUsed, want)
	}

	res.Body.Close()
}
