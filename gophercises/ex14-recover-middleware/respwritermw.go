package main

import "net/http"

type respWriterMw struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *respWriterMw) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *respWriterMw) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *respWriterMw) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}
