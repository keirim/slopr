package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultAPIURL = "https://api.slop.sh"
)

var (
	success = color.New(color.FgGreen, color.Bold)
	info    = color.New(color.FgCyan)
	link    = color.New(color.FgMagenta, color.Underline)
	warn    = color.New(color.FgYellow)
)

type UploadResponse struct {
	ID      string    `json:"id"`
	URL     string    `json:"url"`
	Expires time.Time `json:"expires"`
}

func formatDuration(t time.Time) string {
	duration := t.Sub(time.Now())
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

func uploadFile(filePath string, copyURL bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Show upload starting message
	info.Printf("Uploading %s (%.2f MB)...\n", 
		filepath.Base(filePath), 
		float64(fileInfo.Size())/1024/1024,
	)

	// Prepare multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file to form
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	writer.Close()

	// Send request
	apiURL := viper.GetString("api_url")
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	req, err := http.NewRequest("POST", apiURL+"/upload", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Print the response in a clean format
	fmt.Println()
	success.Println("Upload successful! 🚀")
	link.Printf("URL: %s\n", uploadResp.URL)
	info.Printf("Expires: %s ", uploadResp.Expires.Format("2006-01-02"))
	fmt.Printf("(in %s)\n", formatDuration(uploadResp.Expires))

	// Try to copy URL to clipboard if requested
	if copyURL {
		if err := clipboard.WriteAll(uploadResp.URL); err != nil {
			warn.Println("\nNote: Could not copy to clipboard. URL is displayed above.")
		} else {
			success.Println("\nURL copied to clipboard! 📋")
		}
	}

	return nil
}

func main() {
	var noCopy bool

	var rootCmd = &cobra.Command{
		Use:   "slop [file]",
		Short: "Upload files to the temporary file server",
		Long: `A CLI tool for uploading files to the temporary file server.
Files are automatically deleted after 7 days.

Example:
  slop image.png     Upload a file and copy URL to clipboard
  slop --no-copy file.txt   Upload without copying to clipboard`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uploadFile(args[0], !noCopy)
		},
	}

	rootCmd.Flags().BoolVar(&noCopy, "no-copy", false, "Don't copy URL to clipboard")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/slop")
	viper.SetEnvPrefix("SLOP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
