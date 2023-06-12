package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"net/http"

	"database/sql"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/noirbizarre/gonja"
	"golang.org/x/crypto/pbkdf2"

	"github.com/go-redis/redis"

	"github.com/google/uuid"

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

	db := c.MustGet("db").(*sql.DB)

	// find the staging repository id from the default distro base path ...
	tpl2, err2 := gonja.FromString(database_queries.GetRepositoryIdByName)
	if err2 != nil {
		fmt.Println(err2)
	}
	logger.Debug("---------------------------------------------------")
	get_qs2, err := tpl2.Execute(gonja.Context{"repository_name": "staging"})
	logger.Debug(get_qs2)
	logger.Debug("---------------------------------------------------")

	repo_rows2, err := galaxy_database.ExecuteQueryWithDatabase(get_qs2, db)
	if err != nil {
		fmt.Println(err)
	}
	staging_repository := fmt.Sprintf("%v", repo_rows2[0]["pulp_id"])
	logger.Debug(fmt.Sprintf("%v", staging_repository))
	//staging_repository_id := fmt.Sprintf("%v", staging_repository["pulp_id"])

	// need the sha256sum
	sha256Value := c.PostForm("sha256")
	logger.Debug(fmt.Sprintf("incoming sha256 %v", sha256Value))

	file, err := c.FormFile("file")
	if err != nil {
		logger.Debug(fmt.Sprintf("Error retrieving the file: %s", err.Error()))
		c.String(400, "Bad Request")
		return
	}
	logger.Debug(fmt.Sprintf("file %v", file))
	fileName := file.Filename
	namespace := utils.TarFilenameToNamespace(fileName)

	// Save the uploaded file to disk
	//tmp_filename := sha256Value
	tmp_filepath := "/tmp/" + fileName
	logger.Debug(tmp_filepath)
	/*
		err = c.SaveUploadedFile(file, tmp_filepath)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error saving the file: %s", err.Error()))
			c.String(500, "Internal Server Error")
			return
		}
	*/
	utils.SaveUploadedFile(c, tmp_filepath)

	fileInfo, err := os.Stat(tmp_filepath)
	fileSize := fileInfo.Size()

	//utils.PrintTarballFilenames(tmp_filepath)

	// save it to s3 or to /var/lib/pulp
	if !settings.Use_s3 {
		c.JSON(500, gin.H{"message": "on disk uploads not yet implemented"})
		return
	}
	sha256ValueFirstTwo := sha256Value[:2]
	sha256ValueLastPart := sha256Value[2:]
	s3Destination := settings.Aws_s3_bucket_name + "/artifact/" + sha256ValueFirstTwo + "/" + sha256ValueLastPart
	logger.Debug(s3Destination)
	galaxy_aws.PutS3ObjectByFilepath(tmp_filepath, s3Destination)

	/**************************************************
	    Schema requirements
	**************************************************/

	// need the default domain ID ...
	domain_id := galaxy_database.RunQueryAndReturnColumnByName(
		db,
		database_queries.GetDefaultDomainID,
		"pulp_id",
	)
	logger.Debug(fmt.Sprintf("domain id [%v]", domain_id))

	// need a timestamp
	currentTime := time.Now().UTC()
	logger.Debug(fmt.Sprintf("timestamp %v", currentTime))

	// need a new pulp_id
	newUUID := fmt.Sprintf("%s", uuid.New())
	logger.Debug(newUUID)

	// need the logging cid
	logging_cid := strings.ReplaceAll(newUUID, "-", "")

	// need the username
	username, _, _ := c.Request.BasicAuth()

	// need all the checksums
	//fmd5 := utils.GetFileMd5Sum(tmp_filepath)
	//fsha1 := utils.GetFileSha1Sum(tmp_filepath)
	fsha224 := utils.GetFileSha224Sum(tmp_filepath)
	fsha256 := utils.GetFileSha256Sum(tmp_filepath)
	fsha384 := utils.GetFileSha384Sum(tmp_filepath)
	fsha512 := utils.GetFileSha512Sum(tmp_filepath)

	/**************************************************
	    Make a pulp artifact
	**************************************************/
	stmt, err := db.Prepare(`
        INSERT INTO core_artifact (
            pulp_id,
            pulp_domain_id,
            pulp_created,
			timestamp_of_interest,
            file,
            size,
			sha224,
            sha256,
			sha384,
			sha512
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
    `)
	_, err = stmt.Exec(
		newUUID,
		domain_id,
		currentTime,
		currentTime,
		"artifact/"+sha256ValueFirstTwo+"/"+sha256ValueLastPart,
		fileSize,
		fsha224,
		fsha256,
		fsha384,
		fsha512,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	logger.Info("artifact insert successful")

	/**************************************************
	    Make an upload task
	**************************************************/

	// need args
	//args := ""

	// need kwargs
	kwargs := fmt.Sprintf(`{
        "repository_pk": "%s",
        "data": {
            "sha256": "%s",
            "artifact": "%s",
            "repository": "%s"
        },
        "context": {
            "filename": "%s"
        },
        "username": "%s",
        "filename_ns": "%s",
        "general_args": [
            "ansible",
            "CollectionVersionUploadSerializer"
        ]
    }`,
		staging_repository,
		sha256Value,
		"/api/pulp/api/v3/artifacts/"+newUUID+"/",
		"/api/pulp/api/v3/repositories/ansible/ansible/"+staging_repository+"/",
		fileName,
		username,
		namespace,
	)
	logger.Debug(kwargs)

	/*
		// create a task row ...
		tpl3, err := gonja.FromString(database_queries.NewArtifactUploadTask)
		if err != nil {
			fmt.Println(err)
		}
		task_qs, err := tpl3.Execute(gonja.Context{
			"pulp_created":   currentTime,
			"pulp_id":        newUUID,
			"logging_cid":    logging_cid,
			"pulp_domain_id": domain_id,
			"args":           args,
			"kwargs":         kwargs,
		})
		logger.Debug(task_qs)
	*/

	task_stmt, err := db.Prepare(`
        INSERT INTO core_task (
			pulp_id,
			pulp_domain_id,
			pulp_created,
			state,
			name,
			logging_cid,
			kwargs
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )
    `)
	_, err = task_stmt.Exec(
		newUUID,
		domain_id,
		currentTime,
		"waiting",
		"galaxy_ng.app.tasks.publishing.import_and_auto_approve",
		logging_cid,
		kwargs,
	)

	// cursor.execute("NOTIFY pulp_worker_wakeup")
	db.Exec("NOTIFY pulp_worker_wakeup")

	c.JSON(200, gin.H{"task": newUUID})
}


func (g *Galaxy) ApiV3CollectionImportTask(c *gin.Context) {
}


func (g *Galaxy) ApiUiV1NamespaceCreate(c *gin.Context) {
	/*
	   {
	       "pulp_href":"/api/pulp/api/v3/pulp_ansible/namespaces/39/",
	       "id":39,
	       "name":"foo",
	       "company":"",
	       "EMail":"",
	       "avatar_url":"",
	       "description":"",
	       "links":[],
	       "groups":[],
	       "resources":"",
	       "related_fields":{},
	       "metadata_sha256":"afc89e99d9af497eccccd2aba782ff446d195472117a49752226371537162ba8",
	       "avatar_sha256":null
	   }
	*/
	/*
	   {
	       "errors":[
	           {
	               "status":"409",
	               "code":"conflict",
	               "title":"Data conflicts with existing entity.",
	               "detail":"A namespace named foo already exists.",
	               "source":{"parameter":"name"}
	           }
	       ]
	   }
	*/

	db := c.MustGet("db").(*sql.DB)

	var data database_queries.NamespacePostData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	namespaceName := data.Name

	// does the namespace already exist?
	tpl, err := gonja.FromString(database_queries.CheckIfNamespaceExists)
	if err != nil {
		fmt.Println(err)
	}
	count_qs, err := tpl.Execute(gonja.Context{"name": namespaceName})
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug(count_qs)
	count_rows, err := galaxy_database.ExecuteQueryWithDatabase(count_qs, db)
	if err != nil {
		fmt.Println(err)
	}
	count := count_rows[0]["count"]
	count_int := int(count.(int64))
	logger.Debug(fmt.Sprintf("rowcount %v", count_int))

	if count_int >= 1 {
		error_data := []map[string]interface{}{
			{
				"status": "400",
				"code":   "conflict",
				"title":  "Data conflicts with existing entity.",
				"detal":  fmt.Sprintf("A namespace named %v already exists", namespaceName),
			},
		}
		c.JSON(400, error_data)
		return
	}

	// make a new one
	tpl, err = gonja.FromString(database_queries.CreateNamespace)
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug("---------------------------------------------------")
	create_qs, err := tpl.Execute(gonja.Context{"namespace_name": namespaceName})
	logger.Debug(create_qs)
	logger.Debug("---------------------------------------------------")
	create_rows, err := galaxy_database.ExecuteQueryWithDatabase(create_qs, db)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		c.JSON(500, gin.H{"error": fmt.Sprintf("%v", err)})
	}
	logger.Debug(fmt.Sprintf("%v", create_rows))

	// return the id, name and pulp_href
	tpl, err = gonja.FromString(database_queries.GetNamespaceByName)
	if err != nil {
		fmt.Println(err)
	}
	logger.Debug("---------------------------------------------------")
	get_qs, err := tpl.Execute(gonja.Context{"namespace_name": namespaceName})
	logger.Debug(get_qs)
	logger.Debug("---------------------------------------------------")

	ns_rows, err := galaxy_database.ExecuteQueryWithDatabase(get_qs, db)
	if err != nil {
		fmt.Println(err)
	}
	namespace := ns_rows[0]
	logger.Debug(fmt.Sprintf("%v", namespace))

	c.JSON(200, gin.H{
		"name":       namespace["name"],
		"id":         namespace["id"],
		"company":    namespace["company"],
		"email":      namespace["email"],
		"avatar_url": namespace["avatar_url"],
	})
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

		// create templater
		tpl, err := gonja.FromString(database_queries.GetUsernameAndPassword)
		if err != nil {
			fmt.Println(err)
		}

		// render the query
		qs, err := tpl.Execute(gonja.Context{"username": username})
		if err != nil {
			fmt.Println(err)
		}
		logger.Debug(qs)

		// run query
		rows, err := galaxy_database.ExecuteQueryWithDatabase(qs, db)
		if err != nil {
			fmt.Println(err)
		}

		// extract hash
		db_hash, ok := rows[0]["password"].(string)
		parts := strings.Split(db_hash, "$")
		fmt.Println(parts)

		// extract salt
		db_salt := parts[2]
		fmt.Println(db_salt)

		// extract iterations
		db_iterations_string := parts[1]
		db_iterations, err := strconv.Atoi(db_iterations_string)
		fmt.Println(db_iterations)

		// compile a new hash based on input password
		hashed := pbkdf2.Key([]byte(password), []byte(db_salt), db_iterations, 32, sha256.New)
		encoded_hash := base64.StdEncoding.EncodeToString(hashed)
		logger.Info(fmt.Sprintf("hashed pw %v", encoded_hash))
		salted_hash := "pbkdf2_sha256" + "$" + strconv.Itoa(db_iterations) + "$" + db_salt + "$" + encoded_hash
		fmt.Println(salted_hash)

		if salted_hash == db_hash {
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

	r.POST("/api/_ui/v1/namespaces/", authMiddleware(), galaxy.ApiUiV1NamespaceCreate)

	// ... | POST "/api/v3/artifacts/collections/"
	r.POST("/api/v3/artifacts/collections/", authMiddleware(), galaxy.ApiV3ArtifactPublish)
	r.GET("/api/v3/imports/collections/:taskid/", galaxy.ApiV3CollectionImportTask)

	//r.Static("/artifacts", amanda.Artifacts)
	r.Run("0.0.0.0:" + port)
}
