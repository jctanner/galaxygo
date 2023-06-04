package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"
	"time"

	"net/http"

	"database/sql"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"golang.org/x/crypto/pbkdf2"

	"github.com/go-redis/redis"

	"github.com/jctanner/galaxygo/pkg/database_queries"
	"github.com/jctanner/galaxygo/pkg/galaxy_aws"
	"github.com/jctanner/galaxygo/pkg/galaxy_database"
	"github.com/jctanner/galaxygo/pkg/galaxy_logger"
	"github.com/jctanner/galaxygo/pkg/galaxy_settings"
	"github.com/jctanner/galaxygo/pkg/utils"
)

type Galaxy struct{}

var settings = galaxy_settings.NewGalaxySettings()
var logger = galaxy_logger.Logger{}
var redisClient = redis.NewClient(&redis.Options{
	//Addr:     "redis:6379",
	Addr:     settings.Redis_address,
	Password: "", // Provide password if required
	DB:       0,  // Use default database
})

func (g *Galaxy) Api(c *gin.Context) {
	c.JSON(200, gin.H{
		"available_versions": gin.H{
			"v3": "v3/",
		},
		"current_version": "v1",
		"description":     "Galaxy GoLang",
	})
}

func (g *Galaxy) ApiV3(c *gin.Context) {
	c.JSON(200, gin.H{
		"available_versions": gin.H{
			"collections":        "collections/",
			"collectionversions": "collectionversions/",
		},
	})
}

func (g *Galaxy) ApiV3CollectionsList(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	order_by := c.DefaultQuery("order_by", "pulp_created")

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println(err)
	}
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println(err)
	}

	qs := database_queries.ListCollections + " ORDER BY " + order_by + " LIMIT " + limit + " OFFSET " + offset
	logger.Debug(qs)

	count_rows, err := galaxy_database.ExecuteQueryWithDatabase(database_queries.CountCollections, db)
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug(fmt.Sprintf("rowcount %v", count_rows[0]["count"]))
	count := count_rows[0]["count"]
	count_int := int(count.(int64))

	collection_rows, err := galaxy_database.ExecuteQueryWithDatabase(qs, db)
	if err != nil {
		fmt.Println(err)
	}

	baseurl := "/api/v3/collections/"
	first_url := ""
	previous_url := ""
	previous_offset := offset_int - limit_int
	next_url := ""
	next_offset := limit_int + offset_int
	last_url := ""
	last_offset := (count_int / limit_int)

	if offset_int > 0 {
		previous_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(previous_offset)
	} else {
		previous_url = ""
	}
	first_url = baseurl + "?" + "limit=" + limit + "&offset=0"
	next_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(next_offset)
	last_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(last_offset)

	c.JSON(200, gin.H{
		"meta": gin.H{
			"count": count,
		},
		"links": gin.H{
			"first":    first_url,
			"previous": previous_url,
			"next":     next_url,
			"last":     last_url,
		},
		"data": collection_rows,
	})
}

func (g *Galaxy) ApiV3CollectionSummary(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	namespace := c.Param("namespace")
	name := c.Param("name")

	tpl, err := gonja.FromString(database_queries.CollectionSummary)
	if err != nil {
		fmt.Println(err)
	}

	qs, err := tpl.Execute(gonja.Context{"namespace": namespace, "name": name})
	fmt.Println(qs)

	collection_rows, err := galaxy_database.ExecuteQueryWithDatabase(qs, db)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(collection_rows)

	collection := collection_rows[0]

	c.JSON(200, collection)
}

// r.GET("/api/v3/collections/:namespace/:name/versions/", galaxy.ApiV3CollectionVersionsSummary)
func (g *Galaxy) ApiV3CollectionVersionsSummary(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	namespace := c.Param("namespace")
	name := c.Param("name")

	tpl, err := gonja.FromString(database_queries.CollectionVersionsSummaryCount)
	if err != nil {
		fmt.Println(err)
	}

	count_qs, err := tpl.Execute(gonja.Context{"namespace": namespace, "name": name})
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug(count_qs)

	tpl2, err2 := gonja.FromString(database_queries.CollectionVersionsSummary)
	if err2 != nil {
		fmt.Println(err2)
	}

	versions_qs, err := tpl2.Execute(gonja.Context{"namespace": namespace, "name": name})
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug(versions_qs)

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	order_by := "cc.pulp_created"

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println(err)
	}
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println(err)
	}

	count_rows, err := galaxy_database.ExecuteQueryWithDatabase(count_qs, db)
	if err != nil {
		fmt.Println(err)
	}
	count := count_rows[0]["count"]
	count_int := int(count.(int64))
	logger.Debug(fmt.Sprintf("rowcount %v", count_int))

	qs := versions_qs + " ORDER BY " + order_by + " DESC " + " LIMIT " + limit + " OFFSET " + offset
	logger.Debug(qs)
	collection_rows, err := galaxy_database.ExecuteQueryWithDatabase(qs, db)
	if err != nil {
		fmt.Println(err)
	}

	baseurl := "/api/v3/collections/" + namespace + "/" + name + "/versions/"
	first_url := ""
	previous_url := ""
	previous_offset := offset_int - limit_int
	next_url := ""
	next_offset := limit_int + offset_int
	last_url := ""
	last_offset := (count_int / limit_int)

	if offset_int > 0 {
		previous_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(previous_offset)
	}
	first_url = baseurl + "?" + "limit=" + limit + "&offset=0"

	if next_offset <= count_int {
		next_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(next_offset)
	}

	last_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(last_offset)

	c.JSON(200, gin.H{
		"meta": gin.H{
			"count": count,
		},
		"links": gin.H{
			"first":    first_url,
			"previous": previous_url,
			"next":     next_url,
			"last":     last_url,
		},
		"data": collection_rows,
	})
}

func (g *Galaxy) ApiV3CollectionVersionsList(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	order_by := c.DefaultQuery("order_by", "pulp_created")

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		fmt.Println(err)
	}
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Println(err)
	}

	qs := database_queries.ListCollectionVersions + " ORDER BY " + order_by + " LIMIT " + limit + " OFFSET " + offset
	logger.Debug(qs)
	count_rows, err := galaxy_database.ExecuteQueryWithDatabase(database_queries.CountCollectionVersions, db)
	if err != nil {
		fmt.Println(err)
	}
	count := count_rows[0]["count"]
	count_int := int(count.(int64))
	logger.Debug(fmt.Sprintf("rowcount %v", count_int))

	collection_rows, err := galaxy_database.ExecuteQuery(qs)
	if err != nil {
		fmt.Println(err)
	}

	baseurl := "/api/v3/collectionversions/"
	first_url := ""
	previous_url := ""
	previous_offset := offset_int - limit_int
	next_url := ""
	next_offset := limit_int + offset_int
	last_url := ""
	last_offset := (count_int / limit_int)

	if offset_int > 0 {
		previous_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(previous_offset)
	} else {
		previous_url = ""
	}
	first_url = baseurl + "?" + "limit=" + limit + "&offset=0"
	next_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(next_offset)
	last_url = baseurl + "?" + "limit=" + limit + "&offset=" + strconv.Itoa(last_offset)

	c.JSON(200, gin.H{
		"meta": gin.H{
			"count": count,
		},
		"links": gin.H{
			"first":    first_url,
			"previous": previous_url,
			"next":     next_url,
			"last":     last_url,
		},
		"data": collection_rows,
	})
}

// r.GET("/api/v3/collections/:namespace/:name/versions/:version/", galaxy.ApiV3CollectionVersionDetail)
func (g *Galaxy) ApiV3CollectionVersionDetail(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	namespace := c.Param("namespace")
	name := c.Param("name")
	version := c.Param("version")

	// get the request host ...
	rhost, err := utils.GetRequestHostFromContext(c)

	// make the templater for the cv SQL
	tpl, err := gonja.FromString(database_queries.CollectionVersionDetail)
	if err != nil {
		fmt.Println(err)
	}

	// render the SQL
	qs, err := tpl.Execute(gonja.Context{"namespace": namespace, "name": name, "version": version})
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug(qs)

	// run query
	cv_rows, err := galaxy_database.ExecuteQueryWithDatabase(qs, db)
	if err != nil {
		fmt.Println(err)
	}
	cv := cv_rows[0]

	// cast to string and unmarshal dependencies
	var cv_deps map[string]interface{}
	err2 := json.Unmarshal([]byte(fmt.Sprintf("%v", cv["dependencies"])), &cv_deps)
	if err2 != nil {
		panic(err2)
	}

	// serialize the response
	ds := gin.H{
		"pulp_id": cv["pulp_id"],
		"href":    cv["href"],
		"namespace": gin.H{
			"name": cv["namespace"],
		},
		"name":         cv["name"],
		"version":      cv["version"],
		"created_at":   cv["created_at"],
		"updated_at":   cv["updated_at"],
		"download_url": rhost + fmt.Sprintf("%v", cv["download_url"]),
		"collection": gin.H{
			"id":   cv["collection_id"],
			"name": cv["name"],
			"href": cv["collection_href"],
		},
		"artifact": gin.H{
			"filename": cv["filename"],
			"sha256":   cv["sha256"],
			"size":     cv["size"],
		},
		"requires_ansible": cv["requires_ansible"],
		"metadata": gin.H{
			"authors":       nil,
			"contents":      nil,
			"dependencies":  cv_deps,
			"description":   cv["description"],
			"documentation": cv["documentation"],
			"homepage":      cv["homepage"],
			"issues":        cv["issues"],
			"license":       nil,
			"repository":    cv["repository"],
			"tags":          nil,
		},
		"git_url":        nil,
		"git_commit_sha": nil,
		"manifest": gin.H{
			"format":          nil,
			"collection_info": nil,
		},
		"file_manifest_file": gin.H{
			"name":          nil,
			"ftype":         nil,
			"format":        nil,
			"chksum_type":   nil,
			"chksum_sha256": nil,
		},
		"files": gin.H{
			"files": nil,
		},
	}

	c.JSON(200, ds)
}

func (g *Galaxy) ApiV3Artifact(c *gin.Context) {

	// ArtifactPathByNamespaceNameVersion
	filename := c.Param("filename")
	filepath := ""

	if settings.Use_redis {

		redisCacheKey := "artifact_path_" + filename
		filepath, err := redisClient.Get(redisCacheKey).Result()
		if err != nil {
			filepath = utils.GetFilepathFromDatabase(filename)
			err = redisClient.Set(redisCacheKey, filepath, 5*time.Minute).Err()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {
		filepath = utils.GetFilepathFromDatabase(filename)
	}

	if settings.Use_s3 {
		resp := galaxy_aws.GetS3ObjectByFilepath(filepath)

		// Set the appropriate Content-Type header
		c.Header("Content-Type", *resp.ContentType)

		// Stream the file contents to the client
		_, err := io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to stream file")
			return
		}

	} else {
		baseurl := settings.Content_hostname
		if baseurl == "" {
			baseurl = "http://localhost:5001/api/v3/plugin/ansible/content/community/collections/artifacts/"
		}
		redirect_url := baseurl + filename
		c.Redirect(http.StatusFound, redirect_url)
	}

}

func (g *Galaxy) ApiV3ArtifactPublish(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for GET requests
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		db := c.MustGet("db").(*sql.DB)

		username, password, ok := c.Request.BasicAuth()
		logger.Info(fmt.Sprintf("basic username %v", username))
		logger.Info(fmt.Sprintf("basic password %v", password))
		logger.Info(fmt.Sprintf("basic ok %v", ok))

		form_username := c.PostForm("username")
		form_password := c.PostForm("password")
		logger.Info(fmt.Sprintf("form username %v", form_username))
		logger.Info(fmt.Sprintf("form password %v", form_password))

		//hashedPassword := "pbkdf2_sha256$10000$SALT$HASH"

		var salt string
		iterations := 260000
		salt = "8rZWNpsMtx9HcuhxOVY8qW"
		hashed := pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)
		encoded_hash := base64.StdEncoding.EncodeToString(hashed)
		logger.Info(fmt.Sprintf("hashed pw %v", encoded_hash))
		salted_hash := "pbkdf2_sha256" + "$" + strconv.Itoa(iterations) + "$" + salt + "$" + encoded_hash
		logger.Info(salted_hash)

		// create templater
		tpl, err := gonja.FromString(database_queries.CheckUsernameAndPassword)
		if err != nil {
			fmt.Println(err)
		}

		// render the query
		count_qs, err := tpl.Execute(gonja.Context{"username": username, "password": salted_hash})
		if err != nil {
			fmt.Println(err)
		}
		logger.Debug(count_qs)

		// run query
		count_rows, err := galaxy_database.ExecuteQueryWithDatabase(count_qs, db)
		if err != nil {
			fmt.Println(err)
		}
		count := count_rows[0]["count"]
		count_int := int(count.(int64))
		logger.Debug(fmt.Sprintf("rowcount %v", count_int))

		if count_int >= 1 {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		c.Abort()
		return
	}

}

func main() {
	var artifacts string
	var port string
	galaxy := Galaxy{}

	// https://pkg.go.dev/flag
	//    "flag" implements command-line flag parsing.
	flag.StringVar(&artifacts, "artifacts", "artifacts", "Location of the artifacts dir")
	flag.StringVar(&port, "port", "8080", "Port")
	flag.Parse()

	db, err := galaxy_database.OpenDatabaseConnection()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(location.Default())

	// Middleware to inject the database connection
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// root
	r.GET("/api/", galaxy.Api)

	/*
	   // v1
	   r.GET("/api/v1/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/users/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/users/:userid/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/namespaces/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/namespaces/:namespaceid/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/namespaces/:namespaceid/content/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/namespaces/:namespaceid/owners/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/roles/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/roles/:roleid/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v1/roles/:roleid/versions/", galaxy_proxy.UpstreamHandler)

	   // v2
	   r.GET("/api/v2/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v2/collections/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v2/collections/:namespace/:name/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v2/collections/:namespace/:name/versions/", galaxy_proxy.UpstreamHandler)
	   r.GET("/api/v2/collections/:namespace/:name/versions/:version/", galaxy_proxy.UpstreamHandler)

	   // downloads
	   r.GET("/download/:artifact", galaxy_proxy.ArtifactHandler)
	*/

	// v3
	r.GET("/api/v3/", galaxy.ApiV3)
	r.GET("/api/v3/artifacts/:filename", galaxy.ApiV3Artifact)
	r.GET("/api/v3/artifacts/collections/community/:filename", galaxy.ApiV3Artifact)
	r.GET("/api/v3/collections/", galaxy.ApiV3CollectionsList)
	r.GET("/api/v3/collections/:namespace/:name/", galaxy.ApiV3CollectionSummary)
	r.GET("/api/v3/collections/:namespace/:name/versions/", galaxy.ApiV3CollectionVersionsSummary)
	r.GET("/api/v3/collections/:namespace/:name/versions/:version/", galaxy.ApiV3CollectionVersionDetail)
	r.GET("/api/v3/collectionversions/", galaxy.ApiV3CollectionVersionsList)

	// ... | POST "/api/v3/artifacts/collections/"
	r.POST("/api/v3/artifacts/collections/", authMiddleware(), galaxy.ApiV3ArtifactPublish)

	//r.Static("/artifacts", amanda.Artifacts)
	r.Run("0.0.0.0:" + port)
}
