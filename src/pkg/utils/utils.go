package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jctanner/galaxygo/pkg/database_queries"
	"github.com/jctanner/galaxygo/pkg/galaxy_database"
	"github.com/jctanner/galaxygo/pkg/galaxy_logger"
	"github.com/noirbizarre/gonja"
)

var logger = galaxy_logger.Logger{}

func GetRequestHostFromContext(c *gin.Context) (string, error) {

	// get the request host ...
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	rhost := scheme + "://" + c.Request.Host

	return rhost, nil
}

func GetFilepathFromDatabase(filename string) (filepath string) {

	// assemble the templater
	tpl, err := gonja.FromString(database_queries.ArtifactPathByFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// render the SQL
	qs, err := tpl.Execute(gonja.Context{"filename": filename})
	if err != nil {
		fmt.Println(err)
		return
	}

	// run query
	fp_rows, err := galaxy_database.ExecuteQuery(qs)
	if err != nil {
		fmt.Println(err)
	}

	filepath = fp_rows[0]["filepath"].(string)
	logger.Debug(fmt.Sprintf("pulp filepath %v", filepath))

	return filepath
}

func ShowKeysInMultipartForm(c *gin.Context) string {

	err := c.Request.ParseMultipartForm(32 << 20) // MaxMemory: 32MB
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	// Get the keys from the MultipartForm
	keys := make([]string, 0, len(c.Request.MultipartForm.Value))
	for key := range c.Request.MultipartForm.Value {
		keys = append(keys, key)
	}

	// Print the keys
	return fmt.Sprint("Keys in MultipartForm: %v", keys)
}

func ShowFormData(c *gin.Context) (output string) {

	output = ""

	form := c.Request.Form

	// Dump all the POST data
	for key, values := range form {
		output += fmt.Sprintf("Key: %s\t", key)
		for _, value := range values {
			output += fmt.Sprintf("  Value: %s\n", value)
		}
	}

	return output
}

func ShowFormFile(c *gin.Context) string {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return fmt.Sprintf("Error getting uploaded file: %s", err.Error())
	} else {
		return fmt.Sprintf("Uploaded File Name: %s", file.Filename)
	}
}
