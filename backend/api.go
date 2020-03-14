package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/sergivb01/acmecopy/api"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *Server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.healthy.Load() {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// TODO: implement channels for errors
func (s *Server) readFiles(r *http.Request) []*api.File {
	var files []*api.File
	var wg sync.WaitGroup
	var m sync.Mutex
	for _, h := range r.MultipartForm.File["files"] {
		go func(h *multipart.FileHeader) {
			var buff bytes.Buffer
			file, err := h.Open()
			if err != nil {
				s.log.Error("error copying to file", zap.String("fileName", h.Filename), zap.Error(err))
			}

			if _, err := io.Copy(&buff, file); err != nil {
				s.log.Error("error copying to file", zap.String("fileName", h.Filename), zap.Error(err))
			}

			m.Lock()
			files = append(files, &api.File{
				FileName: h.Filename,
				Content:  buff.Bytes(),
			})
			m.Unlock()
		}(h)
		wg.Add(1)
	}
	wg.Wait()

	return files
}

func (s *Server) handleSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			s.log.Error("parsing multiform data", zap.Error(err))
			return
		}

		res, err := s.cli.CompileFiles(r.Context(), &api.CompileRequest{
			Files: s.readFiles(r),
			Input: []string{
				"hola test123",
				"#",
			},
			ExpectedOutput: []string{
				"ENTRA TEXT ACABAT EN #:",
				"TEXT REVES:",
				"test123 hola ",
			},
		})

		if err != nil {
			_, _ = fmt.Fprintf(w, "error compiling files %s\n", err.Error())
			s.log.Error("error compiling files", zap.Error(err))
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			_, _ = fmt.Fprintf(w, "error encoding to json %v\n", err)
			s.log.Error("error encoding to json", zap.Error(err))
		}
	}
}
