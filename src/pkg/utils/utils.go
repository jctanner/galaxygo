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
	logger.Debug(fmt.Sprintf("filepath %v", filepath))

	return filepath
}
