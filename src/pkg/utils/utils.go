package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	return fmt.Sprintf("Keys in MultipartForm: %v", keys)
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

func SaveUploadedFile(c *gin.Context, dstPath string) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}

	src, err := fileHeader.Open()
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
	defer dst.Close()

	// Decode the base64-encoded data
	decoder := base64.NewDecoder(base64.StdEncoding, src)

	// Copy the decoded data to the destination file
	_, err = io.Copy(dst, decoder)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}
}

func PrintTarballFilenames(filename string) {
	// Open the tarball file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the tarball
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tarball
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Print the filename
		fmt.Println(header.Name)
	}
}

func TarFilenameToNamespace(filename string) string {
	parts := strings.Split(filename, "-")
	return parts[0]
}

func GetFileMd5Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new MD5 hash instance
	hash := md5.New()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}

func GetFileSha1Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA1 hash instance
	hash := sha1.New()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}

func GetFileSha224Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA-224 hash instance
	hash := sha256.New224()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}

func GetFileSha256Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA-256 hash instance
	hash := sha256.New()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}

func GetFileSha384Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA-384 hash instance
	hash := sha512.New384()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}

func GetFileSha512Sum(filename string) string {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA-512 hash instance
	hash := sha512.New()

	// Copy the file content to the hash object
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash sum as a byte slice
	hashSum := hash.Sum(nil)

	// Convert the byte slice to a hex string
	hashString := fmt.Sprintf("%x", hashSum)

	return hashString
}
