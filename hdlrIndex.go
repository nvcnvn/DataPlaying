package main

import (
	"github.com/nvcnvn/DataPlaying/data"
)

func HandleIndex(c *context) {
	viewdata := c.ViewData("title")
	if c.Post("submit", false) == "Send" {
		file, _, _ := c.Request.FormFile("file")
		data.CSVConvert(file)
	}

	c.View("index_detail.tmpl", viewdata)
}
