// slides-content extracts text and image URLs from a Google Slides presentation.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/goauth/google"
	"github.com/grokify/gogoogle/slidesutil/v1"
	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	var (
		credentialsFile       string
		goauthCredentialsFile string
		goauthCredentialsKey  string
		presentationID        string
		includeNotes          bool
		prettyPrint           bool
	)

	flag.StringVar(&credentialsFile, "credentials", os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		"Path to Google service account credentials JSON file")
	flag.StringVar(&goauthCredentialsFile, "goauth-credentials-file", os.Getenv("GOAUTH_CREDENTIALS_FILE"),
		"Path to goauth CredentialsSet JSON file")
	flag.StringVar(&goauthCredentialsKey, "goauth-credentials-account", os.Getenv("GOAUTH_CREDENTIALS_ACCOUNT"),
		"Account key within goauth CredentialsSet file")
	flag.StringVar(&presentationID, "presentation", "", "Google Slides presentation ID (required)")
	flag.BoolVar(&includeNotes, "notes", false, "Include speaker notes")
	flag.BoolVar(&prettyPrint, "pretty", true, "Pretty print JSON output")
	flag.Parse()

	if presentationID == "" {
		fmt.Fprintln(os.Stderr, "Error: -presentation flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate credentials
	hasGoogleCreds := credentialsFile != ""
	hasGoauthCreds := goauthCredentialsFile != "" && goauthCredentialsKey != ""

	if !hasGoogleCreds && !hasGoauthCreds {
		fmt.Fprintln(os.Stderr, "Error: credentials required")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Option 1: Google service account credentials")
		fmt.Fprintln(os.Stderr, "  -credentials /path/to/service-account.json")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Option 2: goauth CredentialsSet")
		fmt.Fprintln(os.Stderr, "  -goauth-credentials-file /path/to/credentials.json -goauth-credentials-account myaccount")
		os.Exit(1)
	}

	if hasGoogleCreds && hasGoauthCreds {
		fmt.Fprintln(os.Stderr, "Error: cannot use both -credentials and -goauth-credentials-file")
		os.Exit(1)
	}

	ctx := context.Background()

	// Create authenticated HTTP client
	var httpClient *http.Client
	var err error

	if hasGoauthCreds {
		httpClient, err = goauth.NewClient(ctx, goauthCredentialsFile, goauthCredentialsKey)
	} else {
		scopes := []string{
			slides.PresentationsReadonlyScope,
			slides.DriveReadonlyScope,
		}
		httpClient, err = google.NewClientSvcAccountFromFile(ctx, credentialsFile, scopes...)
	}
	if err != nil {
		log.Fatalf("Failed to create authenticated client: %v", err)
	}

	// Create Slides service
	svc, err := slides.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("Failed to create Slides service: %v", err)
	}

	// Get presentation
	pres, err := svc.Presentations.Get(presentationID).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to get presentation: %v", err)
	}

	// Extract content
	content := slidesutil.ExtractPresentationContent(pres, includeNotes)

	// Output JSON
	var output []byte
	if prettyPrint {
		output, err = json.MarshalIndent(content, "", "  ")
	} else {
		output, err = json.Marshal(content)
	}
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	fmt.Println(string(output))
}
