// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Json struct {
	logger logger.Logger

	viewModelOrchestrator orchestrator.ViewModelOrchestrator
	fallbackHandler       Handler
}

func (handler *Json) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "json" or ".json" suffix from the path
		path = strings.TrimSuffix(path, "json")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if viewModel, found := handler.viewModelOrchestrator.GetViewModel(requestRoute); found {
			writeJson(w, viewModel)
			return
		}

		// fallback to the item handler
		fallbackHandler := handler.fallbackHandler.Func()
		fallbackHandler(w, r)
	}
}

func writeJson(writer io.Writer, viewModel viewmodel.Model) error {
	bytes, err := json.MarshalIndent(viewModel, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
