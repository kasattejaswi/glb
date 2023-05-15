package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyToServerHandler(t *testing.T) {
	t.Run("Should proxy requests successfully to the origin server", func(t *testing.T) {
		// create a test server to use as the destination host
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("a", "b")
			w.WriteHeader(201)
			fmt.Fprint(w, "Hello, world!")
		}))
		defer testServer.Close()

		// create a request to pass to the ProxyToServerHandler function
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// create a response recorder to capture the response from the ProxyToServerHandler function
		rr := httptest.NewRecorder()

		// call the ProxyToServerHandler function with the test server as the destination host
		handler := ProxyToServerHandler(testServer.URL)
		handler.ServeHTTP(rr, req)

		// check that the response body is "Hello, world!"
		expectedBody := "Hello, world!"
		actualBody := rr.Body.String()
		if actualBody != expectedBody {
			t.Errorf("Unexpected response body. Expected: %s, Actual: %s", expectedBody, actualBody)
		}

		// check that the response status code is 201 OK
		expectedStatusCode := http.StatusCreated
		actualStatusCode := rr.Code
		if actualStatusCode != expectedStatusCode {
			t.Errorf("Unexpected response status code. Expected: %d, Actual: %d", expectedStatusCode, actualStatusCode)
		}

		// check if the response contains the headers returned by the server
		expectedHeaderValue := "b"
		actualHeaderValue := rr.Result().Header.Get("a")
		if actualHeaderValue != expectedHeaderValue {
			t.Errorf("Unexpeced header value. Expected %s, Actual %s", expectedHeaderValue, actualHeaderValue)
		}

	})

	t.Run("Should return 500 response if target server is down", func(t *testing.T) {

		// create a test server which is not running
		testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer testServer.Close()

		// create a request to pass to the ProxyToServerHandler function
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// create a response recorder to capture the response from the ProxyToServerHandler function
		rr := httptest.NewRecorder()

		// call the ProxyToServerHandler function with the test server as the destination host
		handler := ProxyToServerHandler(testServer.URL)
		handler.ServeHTTP(rr, req)
		// check that the response status code is 201 OK
		expectedStatusCode := http.StatusInternalServerError
		actualStatusCode := rr.Code
		if actualStatusCode != expectedStatusCode {
			t.Errorf("Unexpected response status code. Expected: %d, Actual: %d", expectedStatusCode, actualStatusCode)
		}
	})

}
