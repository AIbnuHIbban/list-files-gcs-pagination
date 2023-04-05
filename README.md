
# Google Cloud Storage File List with Pagination
This Golang application demonstrates how to list files stored in a Google Cloud Storage (GCS) bucket with pagination support. The application uses Google Cloud Storage API to retrieve file information and returns the results in a JSON format with next and previous page URLs.

## Prerequisites

1. Go 1.14 or higher installed.

2. A Google Cloud Platform (GCP) project with a GCS bucket containing files.

3. A GCP service account key in JSON format (e.g., `gcp_credential.json`).

### GCPBucket struct

The `GCPBucket` struct represents the GCS bucket configuration, including the bucket name and the JSON credential file.

### FileInfo struct

The `FileInfo` struct represents the metadata of a file in the GCS bucket, such as the file name, creation time, content type, and media URL.

### FileListResponse struct

The `FileListResponse` struct represents the JSON response format, including the limit, next and previous page URLs, current page number, file information, and total number of files returned.

## Usage
1. Clone this repository and navigate to the project directory:

```
git clone https://github.com/AIbnuHibban/gcs-file-list-pagination.git
cd gcs-file-list-pagination
```
2. Replace the `bucketInfo` in `main.go` with your GCS bucket information and the `gcp_credential.json` file with your GCP service account key.

3. Run the application:

```
go run main.go
```

4. Access the `/list` endpoint with your browser or a tool like `curl` or Postman, using the appropriate query parameters, for example:
```
http://localhost:8080/list?page=1&limit=10
```
  
The response will be in JSON format, with the file information and pagination URLs.