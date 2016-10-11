// Copyright 1999-2016. Parallels IP Holdings GmbH.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (self *BackupToAmazon) StartHttp() {
	http.HandleFunc("/", self.index)
	http.HandleFunc("/status/set", self.httpStatusSet)
	http.HandleFunc("/status/unset", self.httpStatusUnSet)

	s := &http.Server{
		Addr:           self.HttpHost + ":" + self.HttpPort,
		Handler:        nil,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		self.HttpStatusSetter = true
		self.Log.Printf("Failed to start http server: %s\n", err)
	}

}

func (self *BackupToAmazon) index(w http.ResponseWriter, r *http.Request) {
	jsEncoder := json.NewEncoder(w)
	err := jsEncoder.Encode(self.HttpStatus)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to show status: %s", err), 400)
		return
	}
}

func (self *BackupToAmazon) httpStatusSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	status := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		self.Log.Printf("Failed to set status: %s", err)
		http.Error(w, fmt.Sprintf("Failed to set status: %s", err), 400)
		return
	}

	for k, v := range status {
		self.HttpStatus[k] = v
	}

	_, err = fmt.Fprint(w, "OK")
	if err != nil {
		self.Log.Printf("Failed to write response: %s", err)
		return
	}
}

func (self *BackupToAmazon) httpStatusUnSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	status := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		self.Log.Printf("Failed to unset status: %s", err)
		http.Error(w, fmt.Sprintf("Failed to set status: %s", err), 400)
		return
	}

	for k, _ := range status {
		delete(self.HttpStatus, k)
	}

	_, err = fmt.Fprint(w, "OK")
	if err != nil {
		self.Log.Printf("Failed to write response: %s", err)
		return
	}
}

func (self *BackupToAmazon) setHttpStatus(key, val string) {
	status := make(map[string]string)
	status[key] = val

	for k, v := range status {
		self.HttpStatus[k] = v
	}

	if !self.HttpStatusSetter {
		self.Log.Printf("HttpStatusSetter: %t\n", self.HttpStatusSetter)
		return
	}

	var payload []byte
	buf := bytes.NewBuffer(payload)
	json.NewEncoder(buf).Encode(self.HttpStatus)
	_, err := http.Post("http://"+self.HttpHost+":"+self.HttpPort+"/status/set", "string", buf)
	if err != nil {
		self.Log.Printf("Failed to set HTTP status: %s\n", err)
	}
	self.Log.Printf("Set HTTP status: %s:%s\n", key, val)
}

func (self *BackupToAmazon) unSetHttpStatus(key, val string) {
	status := make(map[string]string)
	status[key] = val

	for k, v := range status {
		self.HttpStatus[k] = v
	}

	if !self.HttpStatusSetter {
		self.Log.Printf("HttpStatusSetter: %t\n", self.HttpStatusSetter)
		return
	}

	var payload []byte
	buf := bytes.NewBuffer(payload)
	json.NewEncoder(buf).Encode(self.HttpStatus)
	_, err := http.Post("http://"+self.HttpHost+":"+self.HttpPort+"/status/unset", "string", buf)
	if err != nil {
		self.Log.Printf("Failed to set HTTP status: %s\n", err)
	}
	self.Log.Printf("Unset HTTP status: %s:%s\n", key, val)
}

func (self *BackupToAmazon) getHttpStatus() map[string]string {
	r, err := http.Get("http://" + self.HttpHost + ":" + self.HttpPort)
	if err != nil {
		self.Log.Println("HTTP status: nothing in progress")
		return nil
	}
	defer self.IoClose(r.Body)

	status := make(map[string]string)
	err = json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		self.Log.Printf("Failed to get HTTP status: %s", err)
		return nil
	}

	return status
}
