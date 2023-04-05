package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCPBucket represents the GCP bucket configuration.
type GCPBucket struct {
	Name           string `json:"name"`
	CredentialFile string `json:"credentialFile"`
}

type FileInfo struct {
	Name string `json:"name"`
}

type FileListResponse struct {
	Limit    int         `json:"limit"`
	NextPage interface{} `json:"next_page"`
	Page     int         `json:"page"`
	PrevPage interface{} `json:"prev_page"`
	Results  []FileInfo  `json:"results"`
	Total    int         `json:"total"`
}

func main() {
	http.HandleFunc("/list", listFilesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var pageTokenCache sync.Map

// listFilesHandler handles the HTTP request to list all the files in the GCS bucket.
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	bucketInfo := GCPBucket{
		Name:           "bucket-name",
		CredentialFile: "gcp_credential.json",
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	pageToken := r.URL.Query().Get("pageToken")
	page := 1
	if r.URL.Query().Get("page") != "" {
		page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	}

	files, nextPageToken, err := listFilesInBucket(bucketInfo, limit, pageToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the next page token in the cache with the current page as the key
	pageTokenCache.Store(page, nextPageToken)

	// Retrieve the previous page token from the cache
	var prevPageToken interface{}
	if page > 2 {
		prevPageToken, _ = pageTokenCache.Load(page - 2)
	}

	// Generate the next and previous page URLs
	var nextPage, prevPage interface{}
	if nextPageToken != "" {
		nextPage = fmt.Sprintf("http://localhost:8080%s?page=%d&limit=%d&pageToken=%s", r.URL.Path, page+1, limit, nextPageToken)
	}
	if page == 2 {
		prevPage = fmt.Sprintf("http://localhost:8080%s?page=%d&limit=%d", r.URL.Path, page-1, limit)
	} else if page > 2 {
		prevPage = fmt.Sprintf("http://localhost:8080%s?page=%d&limit=%d&pageToken=%s", r.URL.Path, page-1, limit, prevPageToken)
	}

	response := FileListResponse{
		Limit:    limit,
		NextPage: nextPage,
		Page:     page,
		PrevPage: prevPage,
		Results:  files,
		Total:    len(files),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// listFilesInBucket lists all the files in the specified GCS bucket with pagination support.
func listFilesInBucket(bucketInfo GCPBucket, limit int, pageToken string) ([]FileInfo, string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(bucketInfo.CredentialFile))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketInfo.Name)

	it := bucket.Objects(ctx, &storage.Query{
		Prefix: "",
	})

	p := iterator.NewPager(it, limit, pageToken)

	var files []*storage.ObjectAttrs
	nextPageToken, err := p.NextPage(&files)
	if err != nil {
		return nil, "", fmt.Errorf("failed to iterate over objects: %v", err)
	}
	var fileInfos []FileInfo
	for _, file := range files {
		fileInfo := FileInfo{
			Name: file.Name,
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nextPageToken, nil
}
