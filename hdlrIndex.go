package main

import (
	"github.com/nvcnvn/DataPlaying/data"
	"log"
	"net/http"
	"path"
)

var DATASET *data.DataSet

func HandleIndex(c *context) {
	viewdata := c.ViewData("title")

	if c.Post("submit", false) == "Send" {
		file, fHeader, err := c.Request.FormFile("file")
		if err != nil {
			http.Error(c, "Internal Error",
				http.StatusInternalServerError)
			return
		}

		var set *data.DataSet

		if path.Ext(fHeader.Filename) == ".libsvm" {
			set, err = data.LibSVMConvert(file)
		} else {
			set, err = data.CSVConvert(file)
		}

		if err == nil {
			viewdata["dataset"] = set
			DATASET = set
		} else {
			http.Error(c, "Invalid data format",
				http.StatusBadRequest)
			log.Println("HandleIndex", err)
		}
	}

	c.View("index_detail.tmpl", viewdata)
}
