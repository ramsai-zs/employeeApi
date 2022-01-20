package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {

	db = connect()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Println(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/employee?ID=1", nil)
	w := httptest.NewRecorder()
	get(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
	result := `{"ID":1,"Name":"ram","Address":"rjy"}`
	if string(data) != result {
		t.Errorf("expected %v but %v", result, string(data))
	}
}

func TestPost(t *testing.T) {
	db = connect()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Println(err)
	}

	tables := []struct {
		//input
		method string
		body   []byte
		//output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"POST", []byte(`{"ID":32,"Name":"Mahesh","Address":"rhgd"}`), http.StatusOK, []byte("Success")},
	}

	for _, v := range tables {
		req := httptest.NewRequest(v.method, "/employee", bytes.NewReader(v.body))
		w := httptest.NewRecorder()
		post(w, req)
		if w.Code != v.expectedStatusCode {
			t.Errorf("expected %v but got %v", v.expectedStatusCode, w.Code)
		}

		result := bytes.NewBuffer(v.expectedResponse)
		if !reflect.DeepEqual(w.Body, result) {
			t.Errorf("Expected %v but got %v", result.String(), w.Body.String())
		}
	}
}

func TestPut(t *testing.T) {
	db = connect()
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Errorf("Excepted nil but got %v", err)
	}

	tables := []struct {
		//input
		method string
		body   []byte
		//output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"PUT", []byte(`{ID:40,"Name":"Sai"}`), http.StatusOK, []byte("Success")},
	}

	for _, v := range tables {
		res := httptest.NewRequest(v.method, "/employee", bytes.NewReader(v.body))
		w := httptest.NewRecorder()
		put(w, res)
		if w.Code != v.expectedStatusCode {
			t.Errorf("expected %v but got %v", v.expectedResponse, w.Code)
		}

		result := bytes.NewBuffer(v.expectedResponse)
		if !reflect.DeepEqual(w.Body, result) {
			t.Errorf("expected %v but got %v", result.String(), w.Body.String())
		}
	}
}

func TestDeleteId(t *testing.T) {
	db = connect()
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	tables := []struct {
		//input
		method string
		body   []byte
		id     int
		//output
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"DELETE", nil, 32, http.StatusOK, []byte("success")},
	}

	for _, v := range tables {
		query := "/employee?ID=" + strconv.Itoa(v.id)
		res := httptest.NewRequest(v.method, query, bytes.NewReader(v.body))
		w := httptest.NewRecorder()
		deleteId(w, res)
		if w.Code != v.expectedStatusCode {
			t.Errorf("expected %v but got %v", v.expectedStatusCode, w.Code)
		}
		result := bytes.NewBuffer(v.expectedResponse)
		if !reflect.DeepEqual(w.Body, result) {
			t.Errorf("expected %v but got %v", result.String(), w.Body.String())
		}
	}
}
