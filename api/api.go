package main

import (
	"flag"
	"fmt"
    "strconv"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"

    "github.com/jctanner/galaxygo/pkg/database_queries"
    "github.com/jctanner/galaxygo/pkg/galaxy_database"
)


type Galaxy struct {}


func (g *Galaxy) Api(c *gin.Context) {
    c.JSON(200, gin.H{
        "available_versions": gin.H{
            "v3": "v3/",
        },
        "current_version": "v1",
        "description": "Galaxy GoLang",
    })
}


func (g *Galaxy) ApiV3(c *gin.Context) {
    c.JSON(200, gin.H{
        "available_versions": gin.H{
            "collections": "collections/",
            "collectionversions": "collectionversions/",
        },
    })
}


func (g *Galaxy) ApiV3CollectionsList(c *gin.Context) {

    limit := c.DefaultQuery("limit", "10")
    offset := c.DefaultQuery("offset", "0")
    order_by := c.DefaultQuery("order_by", "pulp_created")

    limit_int,err := strconv.Atoi(limit)
    if err != nil {
        fmt.Println(err)
    }
    offset_int,err := strconv.Atoi(offset)
    if err != nil {
        fmt.Println(err)
    }

    qs := database_queries.ListCollections + " ORDER BY " + order_by + " LIMIT " + limit + " OFFSET " + offset

    count_rows,err := galaxy_database.ExecuteQuery(database_queries.CountCollections)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(count_rows[0]["count"])
    count := count_rows[0]["count"]
    count_int := int(count.(int64))

    collection_rows,err := galaxy_database.ExecuteQuery(qs)
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

    if (offset_int > 0) {
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
            "first": first_url,
            "previous": previous_url,
            "next": next_url,
            "last": last_url,
        },
        "data": collection_rows,
    })
}


func (g *Galaxy) ApiV3CollectionVersionsList(c *gin.Context) {

    limit := c.DefaultQuery("limit", "10")
    offset := c.DefaultQuery("offset", "0")
    order_by := c.DefaultQuery("order_by", "pulp_created")

    limit_int,err := strconv.Atoi(limit)
    if err != nil {
        fmt.Println(err)
    }
    offset_int,err := strconv.Atoi(offset)
    if err != nil {
        fmt.Println(err)
    }

    qs := database_queries.ListCollectionVersions + " ORDER BY " + order_by + " LIMIT " + limit + " OFFSET " + offset

    count_rows,err := galaxy_database.ExecuteQuery(database_queries.CountCollectionVersions)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(count_rows[0]["count"])
    count := count_rows[0]["count"]
    count_int := int(count.(int64))

    collection_rows,err := galaxy_database.ExecuteQuery(qs)
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

    if (offset_int > 0) {
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
            "first": first_url,
            "previous": previous_url,
            "next": next_url,
            "last": last_url,
        },
        "data": collection_rows,
    })
}


func main() {
    var artifacts string
    var port string
    galaxy := Galaxy{}

	// https://pkg.go.dev/flag
	//	Package flag implements command-line flag parsing.
    flag.StringVar(&artifacts, "artifacts", "artifacts", "Location of the artifacts dir")
    flag.StringVar(&port, "port", "8080", "Port")
    flag.Parse()

    r := gin.Default()
    r.RedirectTrailingSlash = true
    r.Use(location.Default())

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
    r.GET("/api/v3/collections/", galaxy.ApiV3CollectionsList)
    r.GET("/api/v3/collectionversions/", galaxy.ApiV3CollectionVersionsList)

    //r.Static("/artifacts", amanda.Artifacts)
    r.Run("0.0.0.0:" + port)
}
