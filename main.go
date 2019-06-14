package main

import "context"
import "fmt"
import "io/ioutil"
import "log"
import "cloud.google.com/go/storage"
import "google.golang.org/api/option"
import "golang.org/x/oauth2/google"
import "google.golang.org/api/iterator"

func list_bucket(client *storage.Client, bucketName string) {
	// List the objects in the bucket

	fmt.Println("Listing bucket ", bucketName)
	fmt.Println("--------------------------------------------------")

	ctx := context.Background()

	it := client.Bucket(bucketName).Objects(ctx, nil)

	for {
		attrs, err := it.Next()

		if err == iterator.Done {
			break;
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(attrs.Name)
	}
}

func main() {
	ctx := context.Background()

	first_sa := "first-service-account.json"
	bucketName := "replace-with-your-bucket-name"
	objectName := "second-service-account.json"

	// **********************************************************************
	// Phase 1
	// In this phase we will use the local service account JSON file
	// "first-service-account.json" to create a Cloud Storage client
	// This method loads credentials from a file
	// **********************************************************************

	fmt.Println("Phase 1")

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(first_sa))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// **********************************************************************
	// Phase 2
	// Try to list the objects in the bucket.
	// This should fail
	// **********************************************************************

	fmt.Println("Phase 2")

	list_bucket(client, bucketName)

	// **********************************************************************
	// Phase 3
	// Read the second service account stored in the bucket
	// **********************************************************************

	fmt.Println("Phase 3")

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)

	if err != nil {
		log.Fatalf("Failed to read object: %v", err)
	}

	defer rc.Close()

	data, err := ioutil.ReadAll(rc)

	if err != nil {
		log.Fatalf("Failed to read object: %v", err)
	}

	// fmt.Println("Data:", string(data))

	// **********************************************************************
	// Phase 4
	// Create credentials from second-service-account.json (in-memory data)
	// This method loads credentials from memory
	// **********************************************************************

	fmt.Println("Phase 4")

	creds, err := google.CredentialsFromJSON(ctx, data, storage.ScopeFullControl)

	// **********************************************************************
	// Phase 5
	// Create a new client from the second service account
	// **********************************************************************

	fmt.Println("Phase 5")

	client2, err := storage.NewClient(ctx, option.WithCredentials(creds))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// **********************************************************************
	// Phase 6
	// Try to list the objects in the bucket.
	// This should succeed
	// **********************************************************************

	fmt.Println("Phase 6")

	list_bucket(client2, bucketName)
}
