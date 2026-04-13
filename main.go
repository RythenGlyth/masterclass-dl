package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"go.nhat.io/cookiejar"
)

// sanitizeFilename removes or replaces characters that are illegal in filenames.
// On Windows: < > : " / \ | ? *
// Also handles characters that are problematic across platforms.
func sanitizeFilename(name string) string {
	// Characters to replace (illegal on Windows, problematic elsewhere)
	replacer := strings.NewReplacer(
		"?", "#",
		":", "#",
		"<", "",
		">", "",
		"\"", "'",
		"/", "-",
		"\\", "-",
		"|", "-",
		"*", "",
	)
	result := replacer.Replace(name)

	// On Windows, also handle trailing dots and spaces which are problematic
	if runtime.GOOS == "windows" {
		result = strings.TrimRight(result, ". ")
	}

	return result
}

// Common constants for requests
const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36"
	ja3       = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513,29-23-24,0"
)

func main() {
	var datDir string
	var debug bool
	var rootCmd = &cobra.Command{
		Use:   "masterclass-dl",
		Short: "A downloader for classes from masterclass.com",
	}
	rootCmd.PersistentFlags().StringVarP(&datDir, "datDir", "d", "", "Path to the directory where cookies and other data will be stored (default: $HOME/.masterclass/)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if datDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			datDir = path.Join(home, ".masterclass")
		}
		if _, err := os.Stat(datDir); os.IsNotExist(err) {
			if err := os.MkdirAll(datDir, 0755); err != nil {
				return err
			}
		}
		return nil
	}

	var outputDir string
	var downloadPdfs bool
	var downloadPosters bool
	var ytdlExec string
	var limit int
	var nameAsSeries bool
	var writeNfo bool
	var metadataOnly bool
	var forceDownload bool
	var concurrency int
	var subsOnly bool
	var downloadCmd = &cobra.Command{
		Use:     "download [class/chapter/category...]",
		Aliases: []string{"dl"},
		Short:   "Download a class, chapter, or category from masterclass.com",
		Long: `Download a class, chapter, or category from masterclass.com.
You can either specify a url or just the id. You can specify multiple URLs to download multiple at once.

Supported URL formats:
  - Class:    https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking
  - Chapter:  https://www.masterclass.com/classes/gordon-ramsay-teaches-cooking/chapters/introduction
  - Category: https://www.masterclass.com/homepage/science-and-tech`,
		Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			// Log enabled options
			if metadataOnly {
				fmt.Println("--metadata-only specified, downloading poster, fanart, and NFO only (no videos)")
			}
			if writeNfo && !metadataOnly {
				fmt.Println("--write-nfo specified, will write tvshow.nfo file")
			}
			if nameAsSeries {
				fmt.Println("--name-files-as-series specified, will use s01e01 naming format")
			}
			if !downloadPdfs && !metadataOnly {
				fmt.Println("--pdfs=false specified, skipping PDF downloads")
			}
			if !downloadPosters && !metadataOnly {
				fmt.Println("--posters=false specified, skipping poster/fanart downloads")
			}
			if limit != 10 {
				fmt.Printf("--limit=%d specified for category downloads\n", limit)
			}
			fmt.Println()

			for _, arg := range args {
				// Check if this is a category/homepage URL
				if strings.Contains(arg, "/homepage/") {
					err := downloadCategory(getClient(datDir), datDir, outputDir, downloadPdfs, downloadPosters, ytdlExec, limit, nameAsSeries, writeNfo, metadataOnly, forceDownload, concurrency, subsOnly, arg)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					err := download(getClient(datDir), datDir, outputDir, downloadPdfs, downloadPosters, ytdlExec, nameAsSeries, writeNfo, metadataOnly, forceDownload, concurrency, subsOnly, arg)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		},
	}
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	downloadCmd.Flags().BoolVarP(&downloadPdfs, "pdfs", "p", true, "Download PDFs")
	downloadCmd.Flags().BoolVar(&downloadPosters, "posters", true, "Download poster and fanart images")
	downloadCmd.Flags().StringVarP(&ytdlExec, "ytdl-exec", "y", "yt-dlp", "Path to the youtube-dl or yt-dlp executable")
	downloadCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of classes to download from a category (0 for unlimited)")
	downloadCmd.Flags().BoolVar(&nameAsSeries, "name-files-as-series", false, "Name files in TV series format (s01e01-Title.mp4)")
	downloadCmd.Flags().BoolVar(&writeNfo, "write-nfo", false, "Write tvshow.nfo metadata file for Plex/Jellyfin")
	downloadCmd.Flags().BoolVar(&metadataOnly, "metadata-only", false, "Download only metadata (poster, fanart, NFO) - no videos or PDFs")
	downloadCmd.Flags().BoolVar(&forceDownload, "force", false, "Re-download files even if they already exist")
	downloadCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent fragment downloads for yt-dlp")
	downloadCmd.Flags().BoolVarP(&subsOnly, "subs-only", "s", false, "Download only subtitles (no video)")
	downloadCmd.MarkFlagRequired("output")

	var loginCmd = &cobra.Command{
		Use:   "login [email] [password]",
		Short: "Login to masterclass.com",
		Long: `Login to masterclass.com with your email and password.

Password is resolved in this order:
  1. Command-line argument (not recommended — visible in shell history and process list)
  2. MASTERCLASS_PASSWORD environment variable
  3. Interactive secure prompt (recommended)`,
		Args: cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			email := args[0]
			var password string
			if len(args) >= 2 {
				password = args[1]
			} else if envPass := os.Getenv("MASTERCLASS_PASSWORD"); envPass != "" {
				password = envPass
			} else {
				prompt := promptui.Prompt{
					Label: "Password",
					Mask:  '*',
				}
				var err error
				password, err = prompt.Run()
				if err != nil {
					fmt.Printf("Error reading password: %v\n", err)
					return
				}
			}
			err := login(getClient(datDir), datDir, email, password, debug)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Login successful")
		},
	}

	var loginStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check login status",
		Run: func(cmd *cobra.Command, args []string) {
			err := loginStatus(getClient(datDir), datDir)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	var jsonOutput bool
	var metadataCmd = &cobra.Command{
		Use:   "metadata [url]",
		Short: "Show metadata for a class from masterclass.com",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := showMetadata(getClient(datDir), datDir, jsonOutput, args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	metadataCmd.Flags().BoolVar(&jsonOutput, "json", true, "Output as JSON")

	var safariLoginCmd = &cobra.Command{
		Use:   "safari-login",
		Short: "Login using your existing MasterClass session from Safari (macOS only)",
		Long: `Reads your MasterClass session cookies directly from Safari and writes them
to the masterclass-dl cookie store. Use this instead of 'login' if you
authenticate via SSO (Google, Apple, company SSO) and cannot set a password.

Requirements:
  - macOS only
  - Must be logged in to masterclass.com in Safari
  - Terminal may need Full Disk Access in System Settings → Privacy & Security`,
		Run: func(cmd *cobra.Command, args []string) {
			if runtime.GOOS != "darwin" {
				fmt.Println("safari-login is only supported on macOS")
				os.Exit(1)
			}
			err := safariLogin(getClient(datDir), datDir)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Login successful")
		},
	}

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(loginStatusCmd)
	rootCmd.AddCommand(metadataCmd)
	rootCmd.AddCommand(safariLoginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getClient(datDir string) *http.Client {
	jar := cookiejar.NewPersistentJar(
		cookiejar.WithFilePath(path.Join(datDir, "cookies.json")),
		cookiejar.WithFilePerm(0600),
		cookiejar.WithAutoSync(true),
	)

	return &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
}

// doWithRetry executes an HTTP request with retry and exponential backoff for transient errors.
func doWithRetry(client *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return nil, err
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			resp.Body.Close()
			if attempt < maxRetries {
				wait := time.Duration(1<<uint(attempt)) * time.Second
				fmt.Printf("  Server returned %d, retrying in %v...\n", resp.StatusCode, wait)
				time.Sleep(wait)
				continue
			}
		}
		return resp, nil
	}
	return resp, err
}

func login(client *http.Client, datDir string, email string, password string, debug bool) error {
	// Initialize CycleTLS client to bypass Cloudflare
	cycleclient := cycletls.Init()
	// Note: CycleTLS Close() has issues, so we'll use recover to handle panics
	defer func() {
		if r := recover(); r != nil {
			// Ignore panic from Close(), it's a known issue with CycleTLS
		}
	}()
	defer func() {
		cycleclient.Close()
	}()

	if debug {
		fmt.Printf("Attempting login with email: %s\n", email)
		fmt.Printf("Password length: %d characters\n", len(password))
	}

	// First, visit the home page to establish session (required for login to work)
	if debug {
		fmt.Println("Visiting home page...")
	}
	homePageResp, err := cycleclient.Do("https://www.masterclass.com/", cycletls.Options{
		Body:      "",
		Ja3:       ja3,
		UserAgent: userAgent,
	}, "GET")
	if err != nil {
		return fmt.Errorf("failed to visit home page: %v", err)
	}

	// Build cookie string from home page response
	var cookieStr string
	for _, cookie := range homePageResp.Cookies {
		if cookieStr != "" {
			cookieStr += "; "
		}
		cookieStr += cookie.Name + "=" + cookie.Value
	}

	// Now visit the login page with cookies from home page
	if debug {
		fmt.Println("Visiting login page...")
	}
	loginPageResp, err := cycleclient.Do("https://www.masterclass.com/auth/login", cycletls.Options{
		Body:      "",
		Ja3:       ja3,
		UserAgent: userAgent,
		Headers: map[string]string{
			"Referer": "https://www.masterclass.com/",
			"Cookie":  cookieStr,
		},
	}, "GET")
	if err != nil {
		return fmt.Errorf("failed to visit login page: %v", err)
	}

	if debug && strings.Contains(loginPageResp.Body, "hidden") {
		fmt.Println("Login page contains hidden form fields - might need to extract them")
	}

	// Update cookies from login page response
	for _, cookie := range loginPageResp.Cookies {
		// Check if cookie already exists, update it, otherwise append
		found := false
		for _, existing := range homePageResp.Cookies {
			if existing.Name == cookie.Name {
				existing.Value = cookie.Value
				found = true
				break
			}
		}
		if !found {
			homePageResp.Cookies = append(homePageResp.Cookies, cookie)
		}
	}
	// Get CSRF token
	if debug {
		fmt.Println("Getting CSRF token...")
	}
	csrfResp, err := cycleclient.Do("https://www.masterclass.com/api/v2/csrf-token", cycletls.Options{
		Body: "",
		Ja3:  ja3,
		Headers: map[string]string{
			"Referer": "https://www.masterclass.com/auth/login",
			"Cookie":  cookieStr,
		},
	}, "GET")
	if err != nil {
		return fmt.Errorf("failed to get CSRF token: %v", err)
	}
	if csrfResp.Status != 200 {
		return fmt.Errorf("failed to get CSRF token: status=%d, body=%s", csrfResp.Status, csrfResp.Body)
	}

	var csrfResponse CSRFResponse
	err = json.Unmarshal([]byte(csrfResp.Body), &csrfResponse)
	if err != nil {
		return fmt.Errorf("failed to parse CSRF response: %v", err)
	}
	if csrfResponse.Param == "" || csrfResponse.Token == "" || csrfResponse.Param != "authenticity_token" {
		return fmt.Errorf("invalid CSRF token response: param=%s, token=%s", csrfResponse.Param, csrfResponse.Token)
	}

	// Update cookies from CSRF response
	for _, cookie := range csrfResp.Cookies {
		// Check if cookie already exists, update it, otherwise append
		found := false
		for _, existing := range homePageResp.Cookies {
			if existing.Name == cookie.Name {
				existing.Value = cookie.Value
				found = true
				break
			}
		}
		if !found {
			homePageResp.Cookies = append(homePageResp.Cookies, cookie)
		}
	}

	// Rebuild cookie string
	cookieStr = ""
	for _, cookie := range homePageResp.Cookies {
		if cookieStr != "" {
			cookieStr += "; "
		}
		cookieStr += cookie.Name + "=" + cookie.Value
	}

	// Prepare login data - NO authenticity_token in the body! Only in X-Csrf-Token header
	data := url.Values{}
	data.Set("next_page", "")
	data.Set("auth_key", email)
	data.Set("password", password)
	data.Set("provider", "identity")

	if debug {
		fmt.Println("Logging in...")
		fmt.Printf("Form data: %s\n", data.Encode())
		fmt.Printf("CSRF token (header only): %s\n", csrfResponse.Token)
	}

	// Perform login - headers must match browser (cors, not navigate!)
	loginResp, err := cycleclient.Do("https://www.masterclass.com/auth/identity/callback", cycletls.Options{
		Body:      data.Encode(),
		Ja3:       ja3,
		UserAgent: userAgent,
		Headers: map[string]string{
			"Accept":             "*/*",
			"Accept-Language":    "en-GB,en-US;q=0.9,en;q=0.8",
			"Content-Type":       "application/x-www-form-urlencoded",
			"X-Csrf-Token":       csrfResponse.Token,
			"Cookie":             cookieStr,
			"Referer":            "https://www.masterclass.com/auth/login",
			"Origin":             "https://www.masterclass.com",
			"Sec-Ch-Ua":          "\"Chromium\";v=\"140\", \"Not=A?Brand\";v=\"24\", \"Google Chrome\";v=\"140\"",
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "\"macOS\"",
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-origin",
			"Priority":           "u=1, i",
		},
	}, "POST")
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	// Handle specific error statuses
	if loginResp.Status == 429 {
		return fmt.Errorf("rate limited by Masterclass. Please wait 15-60 minutes before trying again, or use a different network/VPN")
	}

	// Accept 200 or 302 (redirect) as success
	if loginResp.Status != 200 && loginResp.Status != 302 {
		return fmt.Errorf("failed to login: status=%d, body=%s", loginResp.Status, loginResp.Body)
	}

	if debug {
		fmt.Printf("Login response status: %d\n", loginResp.Status)
		fmt.Printf("Login response body length: %d bytes\n", len(loginResp.Body))
	}

	// Check if the HTML contains error indicators
	if strings.Contains(loginResp.Body, "Invalid email") ||
		strings.Contains(loginResp.Body, "Invalid password") ||
		strings.Contains(loginResp.Body, "incorrect email or password") {
		return fmt.Errorf("login failed: invalid credentials")
	}

	// Check all cookies in login response
	hasSessionCookie := false
	if debug {
		fmt.Printf("Number of cookies in login response: %d\n", len(loginResp.Cookies))
	}
	for _, cookie := range loginResp.Cookies {
		if debug {
			valuePreview := cookie.Value
			if len(valuePreview) > 50 {
				valuePreview = valuePreview[:50] + "..."
			}
			fmt.Printf("  Cookie: %s = %s (len: %d)\n", cookie.Name, valuePreview, len(cookie.Value))
		}
		if cookie.Name == "_mc_session" && len(cookie.Value) > 100 {
			hasSessionCookie = true
			if debug {
				fmt.Printf("  ✓ This is a valid authenticated session cookie\n")
			}
		}
	}

	if !hasSessionCookie {
		// Print more of the HTML body to see what page we actually got
		preview := loginResp.Body
		if len(preview) > 3000 {
			preview = preview[:3000]
		}
		if strings.Contains(preview, "<title>") {
			titleStart := strings.Index(preview, "<title>") + 7
			titleEnd := strings.Index(preview, "</title>")
			if titleEnd > titleStart {
				fmt.Printf("Page title: %s\n", preview[titleStart:titleEnd])
			}
		}
		return fmt.Errorf("login failed - no valid session cookie received")
	}

	// Extract cookies and save to cookiejar
	// Use the base domain (without www) so cookies work across all subdomains
	masterclassURL, _ := url.Parse("https://masterclass.com")
	edgeURL, _ := url.Parse("https://edge.masterclass.com")

	// Convert CycleTLS cookies to http.Cookie
	// Create NEW cookie objects with proper domain to ensure they work across subdomains
	var cookies []*http.Cookie

	// Collect all cookies from the session
	allCookies := make(map[string]*http.Cookie)

	// Add from home page
	for _, cookie := range homePageResp.Cookies {
		// Create a new cookie with proper domain
		domain := cookie.Domain
		if domain == "" || domain == "masterclass.com" {
			domain = ".masterclass.com" // Leading dot for subdomain sharing
		}
		if !strings.HasPrefix(domain, ".") && strings.Contains(domain, ".") {
			domain = "." + domain
		}
		newCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   domain,
			Expires:  cookie.Expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: cookie.SameSite,
		}
		if newCookie.Path == "" {
			newCookie.Path = "/"
		}
		allCookies[cookie.Name] = newCookie
	}

	// Update/add from login response (overwrites duplicates)
	for _, cookie := range loginResp.Cookies {
		// Create a new cookie with proper domain
		domain := cookie.Domain
		if domain == "" || domain == "masterclass.com" {
			domain = ".masterclass.com" // Leading dot for subdomain sharing
		}
		if !strings.HasPrefix(domain, ".") && strings.Contains(domain, ".") {
			domain = "." + domain
		}
		newCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   domain,
			Expires:  cookie.Expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
			SameSite: cookie.SameSite,
		}
		if newCookie.Path == "" {
			newCookie.Path = "/"
		}
		allCookies[cookie.Name] = newCookie
	}

	// Convert map to slice
	for _, cookie := range allCookies {
		cookies = append(cookies, cookie)
	}

	if debug {
		fmt.Printf("Saving %d cookies to jar\n", len(cookies))
		for _, c := range cookies {
			fmt.Printf("  - %s (domain: %s, value length: %d)\n", c.Name, c.Domain, len(c.Value))
		}
	}
	// Set cookies on both URLs to ensure they're available for all subdomains
	client.Jar.SetCookies(masterclassURL, cookies)
	client.Jar.SetCookies(edgeURL, cookies)

	// Build clean cookie string from our collected cookies (no duplicates)
	var cleanCookieStr string
	seenCookies := make(map[string]bool)
	for _, cookie := range cookies {
		if !seenCookies[cookie.Name] {
			if cleanCookieStr != "" {
				cleanCookieStr += "; "
			}
			cleanCookieStr += cookie.Name + "=" + cookie.Value
			seenCookies[cookie.Name] = true
		}
	}

	if debug {
		fmt.Printf("Using %d unique cookies\n", len(seenCookies))
	}

	// No need to visit profiles page - we already have the session cookie

	// Now fetch profiles API
	if debug {
		fmt.Println("Fetching profiles API...")
	}
	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/profiles?deep=true", nil)
	if err != nil {
		return fmt.Errorf("failed to create profiles request: %v", err)
	}
	req.Header.Set("Cookie", cleanCookieStr)
	req.Header.Set("Referer", "https://www.masterclass.com/profiles")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get profiles: %v", err)
	}
	defer resp.Body.Close()

	if debug {
		fmt.Printf("Profiles response status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get profiles: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var profiles []ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profiles)
	if err != nil {
		return fmt.Errorf("failed to parse profiles: %v", err)
	}

	if len(profiles) == 0 {
		return fmt.Errorf("no profiles found")
	}

	prompt := promptui.Select{
		Label: "Select Profile",
		Items: profiles,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .DisplayName }}",
			Active:   "\U0001F449 {{ .DisplayName }}",
			Inactive: "  {{ .DisplayName }}",
			Selected: "\U0001F64C {{ .DisplayName }}",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Selected profile: %s\n", profiles[i].DisplayName)

	// Write selected profile to datDir + "/profile.json"
	profileFile, err := os.Create(path.Join(datDir, "profile.json"))
	if err != nil {
		return err
	}
	defer profileFile.Close()
	err = json.NewEncoder(profileFile).Encode(profiles[i])
	if err != nil {
		return err
	}

	return nil
}

func getProfile(client *http.Client, datDir string) (*ProfileResponse, error) {
	profileFile, err := os.Open(path.Join(datDir, "profile.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile not found. Please login first")
		}
		return nil, err
	}
	defer profileFile.Close()
	var profile ProfileResponse
	err = json.NewDecoder(profileFile).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func fetchChapter(client *http.Client, profileUUID string, chapterID int) (*Chapter, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.masterclass.com/jsonapi/v1/chapters/%d?deep=true", chapterID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://www.masterclass.com/")
	req.Header.Set("Mc-Profile-Id", profileUUID)

	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch chapter: status %d", resp.StatusCode)
	}

	var chapter Chapter
	err = json.NewDecoder(resp.Body).Decode(&chapter)
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func loginStatus(client *http.Client, datDir string) error {
	if (client.Jar.Cookies(&url.URL{Scheme: "https", Host: "www.masterclass.com"}) == nil) {
		return fmt.Errorf("cookies not found. Please login first")
	}

	profile, err := getProfile(client, datDir)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/subscriptions/current?include=purchase_plan%2Cpurchase_plan.product%2Crenewal_purchase_plan%2Crenewal_purchase_plan.product", nil)
	req.Header.Set("Mc-Profile-Id", profile.UUID)
	req.Header.Set("Referer", "https://www.masterclass.com/homepage")
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get subscription status")
	}
	var subscription SubscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&subscription)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/user/cart-data?deep=true", nil)
	req.Header.Set("Mc-Profile-Id", profile.UUID)
	req.Header.Set("Referer", "https://www.masterclass.com/homepage")
	if err != nil {
		return err
	}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get login status")
	}
	var cartData CartDataResponse
	err = json.NewDecoder(resp.Body).Decode(&cartData)
	if err != nil {
		return err
	}
	fmt.Printf("Email: %s\n", cartData.Email)
	fmt.Printf("Subscription Status: %s\n", subscription.Status)
	fmt.Printf("Subscription Expires At: %s\n", subscription.ExpiresAt)
	fmt.Printf("Subscription Remaining Days: %d\n", subscription.RemainingDays)
	return nil
}
func showMetadata(client *http.Client, datDir string, jsonOutput bool, arg string) error {
	if client.Jar.Cookies(&url.URL{Scheme: "https", Host: "www.masterclass.com"}) == nil {
		return fmt.Errorf("cookies not found. Please login first")
	}

	profile, err := getProfile(client, datDir)
	if err != nil {
		return err
	}

	// Check if this is a category/homepage URL
	if strings.Contains(arg, "/homepage/") {
		return showCategoryMetadata(client, profile.UUID, jsonOutput, arg)
	}

	// Parse class slug from URL
	classSlug := arg
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/sessions/classes/")
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/classes/")
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/series/")
	classSlug = strings.TrimSuffix(classSlug, "/")
	// Remove any chapter suffix
	if strings.Contains(classSlug, "/chapters/") {
		classSlug = strings.Split(classSlug, "/chapters/")[0]
	}

	if classSlug == "" {
		return fmt.Errorf("invalid class URL")
	}

	return showCourseMetadata(client, profile.UUID, jsonOutput, classSlug)
}

func showCourseMetadata(client *http.Client, profileUUID string, jsonOutput bool, classSlug string) error {
	// Fetch course data
	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/courses/"+classSlug+"?deep=true&include=instructors,chapters,categories", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.masterclass.com/classes/"+classSlug)
	req.Header.Set("Mc-Profile-Id", profileUUID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return fmt.Errorf("failed to get class info: status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Read raw response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if jsonOutput {
		// Pretty print JSON
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			// If can't indent, just print raw
			fmt.Println(string(body))
		} else {
			fmt.Println(prettyJSON.String())
		}
	} else {
		// Parse and show key fields
		var course CourseResponse
		err = json.Unmarshal(body, &course)
		if err != nil {
			return err
		}
		fmt.Printf("Title: %s\n", course.Title)
		fmt.Printf("Skill: %s\n", course.Skill)
		fmt.Printf("Headline: %s\n", course.Headline)
		fmt.Printf("VanityName: %s\n", course.VanityName)
		fmt.Printf("InstructorName: %s\n", course.InstructorName)
		fmt.Printf("Slug: %s\n", course.Slug)
	}

	return nil
}

func showCategoryMetadata(client *http.Client, profileUUID string, jsonOutput bool, arg string) error {
	// Parse the category URL to extract bundle name
	categorySlug := arg
	categorySlug = strings.TrimPrefix(categorySlug, "https://www.masterclass.com/")
	categorySlug = strings.TrimPrefix(categorySlug, "http://www.masterclass.com/")
	categorySlug = strings.TrimSuffix(categorySlug, "/")

	// Convert path to bundle format: homepage/business -> homepage-business
	bundle := strings.ReplaceAll(categorySlug, "/", "-")

	fmt.Printf("Fetching category: %s (bundle: %s)\n", categorySlug, bundle)

	// Call the content-rows API
	apiURL := fmt.Sprintf("https://www.masterclass.com/jsonapi/v3/content-rows?filter[platform]=web&filter[bundle]=%s&include_items=true", bundle)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.masterclass.com/"+categorySlug)
	req.Header.Set("Mc-Profile-Id", profileUUID)

	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get category info: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var contentRows ContentRowsResponse
	err = json.NewDecoder(resp.Body).Decode(&contentRows)
	if err != nil {
		return fmt.Errorf("failed to parse content rows: %v", err)
	}

	// Extract unique course slugs
	courseMap := make(map[string]bool)
	var courseSlugs []string
	for _, row := range contentRows {
		for _, item := range row.Items {
			resource := item.Default.Resource
			if resource.EntitySlug != "" && resource.EntityType == "course" {
				if !courseMap[resource.EntitySlug] {
					courseMap[resource.EntitySlug] = true
					courseSlugs = append(courseSlugs, resource.EntitySlug)
				}
			}
		}
	}

	fmt.Printf("\nFound %d courses in category:\n", len(courseSlugs))
	fmt.Println(strings.Repeat("-", 80))

	if jsonOutput {
		// Output as JSON array of course metadata
		fmt.Println("[")
		for i, slug := range courseSlugs {
			err := showCourseMetadata(client, profileUUID, true, slug)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to get metadata for %s: %v\n", slug, err)
				continue
			}
			if i < len(courseSlugs)-1 {
				fmt.Println(",")
			}
		}
		fmt.Println("]")
	} else {
		// Output as table - fetch all courses first
		type courseInfo struct {
			Slug              string
			Title             string
			Skill             string
			Headline          string
			VanityName        string
			InstructorName    string
			InstructorTagline string
			Type              string
			NumChapters       int
			TotalSeconds      int
		}
		var courses []courseInfo

		for _, slug := range courseSlugs {
			req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/courses/"+slug+"?deep=true&include=instructors,chapters,categories", nil)
			if err != nil {
				continue
			}
			req.Header.Set("Referer", "https://www.masterclass.com/classes/"+slug)
			req.Header.Set("Mc-Profile-Id", profileUUID)

			resp, err := doWithRetry(client, req, 3)
			if err != nil {
				continue
			}

			var course CourseResponse
			err = json.NewDecoder(resp.Body).Decode(&course)
			resp.Body.Close()
			if err != nil {
				continue
			}

			courses = append(courses, courseInfo{
				Slug:              course.Slug,
				Title:             course.Title,
				Skill:             course.Skill,
				Headline:          course.Headline,
				VanityName:        course.VanityName,
				InstructorName:    course.InstructorName,
				InstructorTagline: course.InstructorTagline,
				Type:              course.Type,
				NumChapters:       course.NumChapters,
				TotalSeconds:      course.TotalSeconds,
			})
		}

		// Helper to truncate strings
		trunc := func(s string, max int) string {
			if len(s) > max {
				return s[:max-3] + "..."
			}
			return s
		}

		// Print table header
		fmt.Printf("\n%-3s | %-7s | %-45s | %-45s | %-30s | %-8s | %-30s | %-4s | %-6s | %-40s\n",
			"#", "Type", "Title", "Skill", "Headline", "Vanity", "Instructor", "Chap", "Mins", "Slug")
		fmt.Println(strings.Repeat("-", 240))

		// Print rows
		for i, c := range courses {
			fmt.Printf("%-3d | %-7s | %-45s | %-45s | %-30s | %-8s | %-30s | %-4d | %-6d | %-40s\n",
				i+1,
				c.Type,
				trunc(c.Title, 45),
				trunc(c.Skill, 45),
				trunc(c.Headline, 30),
				trunc(c.VanityName, 8),
				trunc(c.InstructorName, 30),
				c.NumChapters,
				c.TotalSeconds/60,
				trunc(c.Slug, 40))
		}
	}

	return nil
}

// downloadCamp handles MasterClass "Sessions" content which uses the camps API
func downloadCamp(client *http.Client, profileUUID string, outputDir string, downloadPdfs bool, downloadPosters bool, ytdlExec string, nameAsSeries bool, writeNfo bool, metadataOnly bool, forceDownload bool, concurrency int, subsOnly bool, campSlug string, taskSlug string) error {
	fmt.Printf("Detected Sessions content, fetching camp: %s\n", campSlug)

	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/camps/"+campSlug+"?include=camp_modules,camp_modules.camp_tasks,instructors", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.masterclass.com/sessions/classes/"+campSlug)
	req.Header.Set("Mc-Profile-Id", profileUUID)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get session info: status %d", resp.StatusCode)
	}

	var camp CampResponse
	err = json.NewDecoder(resp.Body).Decode(&camp)
	if err != nil {
		return fmt.Errorf("failed to parse session info: %v", err)
	}

	outputDir = path.Join(outputDir, sanitizeFilename(camp.Title))
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	fmt.Printf("Output directory: %s\n", outputDir)

	// Download artwork
	if downloadPosters {
		if camp.Primary2x3 != "" {
			fmt.Println("Downloading poster image")
			err = downloadImage(client, camp.Primary2x3, path.Join(outputDir, "poster.jpg"))
			if err != nil {
				fmt.Printf("Warning: failed to download poster: %v\n", err)
			}
		}
		if camp.Primary16x9 != "" {
			fmt.Println("Downloading fanart image")
			err = downloadImage(client, camp.Primary16x9, path.Join(outputDir, "fanart.jpg"))
			if err != nil {
				fmt.Printf("Warning: failed to download fanart: %v\n", err)
			}
		}
	}

	// Build a flat ordered list of video tasks across all modules
	type videoTask struct {
		moduleTitle string
		modulePos   int
		task        struct {
			ID           int
			Title        string
			Slug         string
			TaskType     string
			DurationSecs int
			Position     int
			ThumbURL     string
		}
		globalIndex int
	}

	var videoTasks []videoTask
	globalIdx := 1
	for _, mod := range camp.CampModules {
		for _, t := range mod.CampTasks {
			if t.TaskType == "video" || t.TaskType == "follow_along_video" {
				vt := videoTask{
					moduleTitle: mod.Title,
					modulePos:   mod.Position,
					globalIndex: globalIdx,
				}
				vt.task.ID = t.ID
				vt.task.Title = t.Title
				vt.task.Slug = t.Slug
				vt.task.TaskType = t.TaskType
				vt.task.DurationSecs = t.DurationSecs
				vt.task.Position = t.Position
				vt.task.ThumbURL = t.ThumbURL
				videoTasks = append(videoTasks, vt)
				globalIdx++
			}
		}
	}

	fmt.Printf("Found %d video tasks across %d modules\n", len(videoTasks), len(camp.CampModules))

	if metadataOnly {
		fmt.Println("Metadata only mode — skipping video download")
		return nil
	}

	apiKey := "b9517f7d8d1f48c2de88100f2c13e77a9d8e524aed204651acca65202ff5c6cb9244c045795b1fafda617ac5eb0a6c50"

	downloadedCount := 0
	for _, vt := range videoTasks {
		if taskSlug != "" && vt.task.Slug != taskSlug {
			continue
		}
		fmt.Printf("Downloading task %d: %s\n", vt.globalIndex, vt.task.Title)
		var downloaded bool
		var err error
		if subsOnly {
			downloaded, err = downloadCampTaskSubsOnly(client, profileUUID, outputDir, ytdlExec, campSlug, vt.globalIndex, vt.task.Slug, vt.task.Title, apiKey)
		} else {
			downloaded, err = downloadCampTask(client, profileUUID, outputDir, ytdlExec, campSlug, vt.globalIndex, len(videoTasks), vt.task.Slug, vt.task.Title, apiKey, nameAsSeries, forceDownload, concurrency)
		}
		if err != nil {
			fmt.Printf("Warning: task %s failed: %v\n", vt.task.Slug, err)
			continue
		}
		if downloaded {
			downloadedCount++
		}
	}

	fmt.Printf("Downloaded %d/%d videos\n", downloadedCount, len(videoTasks))
	return nil
}

func getCampTaskMediaUUID(client *http.Client, profileUUID string, campSlug string, taskSlug string) (string, error) {
	req, err := http.NewRequest("GET",
		"https://www.masterclass.com/jsonapi/v1/camp-tasks/"+taskSlug+"?include=video,video.video_segments",
		nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://www.masterclass.com/sessions/classes/"+campSlug+"/tasks/"+taskSlug)
	req.Header.Set("Mc-Profile-Id", profileUUID)

	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return "", err
	}

	var taskDetail CampTaskDetailResponse
	err = json.NewDecoder(resp.Body).Decode(&taskDetail)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	if taskDetail.Video == nil || taskDetail.Video.MediaUUID == "" {
		return "", nil
	}
	return taskDetail.Video.MediaUUID, nil
}

func downloadCampTaskSubsOnly(client *http.Client, profileUUID string, outputDir string, ytdlExec string, campSlug string, taskNum int, taskSlug string, taskTitle string, apiKey string) (bool, error) {
	mediaUUID, err := getCampTaskMediaUUID(client, profileUUID, campSlug, taskSlug)
	if err != nil {
		return false, err
	}
	if mediaUUID == "" {
		fmt.Printf("Skipping task %s: no media UUID found\n", taskSlug)
		return false, nil
	}

	safeTitle := sanitizeFilename(taskTitle)
	baseFilename := path.Join(outputDir, fmt.Sprintf("%03d-%s", taskNum, safeTitle))

	streamURL, textTracks, err := getChapterStreamInfo(client, profileUUID, mediaUUID, apiKey)
	if err != nil {
		return false, err
	}

	for _, track := range textTracks {
		if track.Src == "" {
			continue
		}
		subFilename := fmt.Sprintf("%s.%s.vtt", baseFilename, track.SrcLang)
		resp, err := client.Get(track.Src)
		if err != nil {
			fmt.Printf("  Warning: failed to download %s subtitle: %v\n", track.Label, err)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			fmt.Printf("  Warning: failed to download %s subtitle: status %d\n", track.Label, resp.StatusCode)
			continue
		}
		subFile, err := os.Create(subFilename)
		if err != nil {
			resp.Body.Close()
			fmt.Printf("  Warning: failed to create subtitle file for %s: %v\n", track.Label, err)
			continue
		}
		_, err = io.Copy(subFile, resp.Body)
		subFile.Close()
		resp.Body.Close()
		if err != nil {
			fmt.Printf("  Warning: failed to write subtitle file for %s: %v\n", track.Label, err)
			continue
		}
		fmt.Printf("  Downloaded subtitle: %s (%s)\n", track.Label, track.SrcLang)
	}

	if streamURL != "" && len(textTracks) == 0 {
		cmd := exec.Command(ytdlExec, "--skip-download", "--write-subs", "--all-subs", "-o", baseFilename+".%(ext)s", streamURL)
		if err := cmd.Run(); err != nil {
			fmt.Printf("  Warning: yt-dlp subtitle extraction failed: %v\n", err)
		}
	}

	return true, nil
}

func downloadCampTask(client *http.Client, profileUUID string, outputDir string, ytdlExec string, campSlug string, taskNum int, totalTasks int, taskSlug string, taskTitle string, apiKey string, nameAsSeries bool, forceDownload bool, concurrency int) (bool, error) {
	mediaUUID, err := getCampTaskMediaUUID(client, profileUUID, campSlug, taskSlug)
	if err != nil {
		return false, err
	}
	if mediaUUID == "" {
		fmt.Printf("Skipping task %s: no media UUID found\n", taskSlug)
		return false, nil
	}

	var baseFileName string
	if nameAsSeries {
		baseFileName = fmt.Sprintf("s01e%02d-%s", taskNum, sanitizeFilename(taskTitle))
	} else {
		baseFileName = fmt.Sprintf("%03d-%s", taskNum, sanitizeFilename(taskTitle))
	}
	outputFile := path.Join(outputDir, baseFileName+".mp4")

	if !forceDownload {
		if _, err := os.Stat(outputFile); err == nil {
			fmt.Printf("Skipping (already exists): %s\n", baseFileName)
			return false, nil
		}
	}

	streamURL, textTracks, err := getChapterStreamInfo(client, profileUUID, mediaUUID, apiKey)
	if err != nil {
		return false, err
	}
	if streamURL == "" {
		fmt.Printf("Skipping task %s: no video sources\n", taskSlug)
		return false, nil
	}

	fmt.Printf("Downloading task %d/%d: %s\n", taskNum, totalTasks, taskTitle)

	ytdlArgs := []string{
		streamURL,
		"-o", outputFile,
		"--no-warnings",
		"--embed-subs",
		"--all-subs",
		"--merge-output-format", "mp4",
		"-f", "bestvideo+bestaudio/best",
		"--concurrent-fragments", fmt.Sprintf("%d", concurrency),
		"--add-metadata",
	}
	if !forceDownload {
		ytdlArgs = append(ytdlArgs, "--no-overwrites")
	}
	for _, track := range textTracks {
		if track.Src != "" {
			ytdlArgs = append(ytdlArgs, "--sub-lang", track.SrcLang)
		}
	}

	ytdlCmd := exec.Command(ytdlExec, ytdlArgs...)
	ytdlCmd.Stdout = os.Stdout
	ytdlCmd.Stderr = os.Stderr
	if err := ytdlCmd.Run(); err != nil {
		return false, fmt.Errorf("yt-dlp failed: %v", err)
	}
	return true, nil
}

func download(client *http.Client, datDir string, outputDir string, downloadPdfs bool, downloadPosters bool, ytdlExec string, nameAsSeries bool, writeNfo bool, metadataOnly bool, forceDownload bool, concurrency int, subsOnly bool, arg string) error {
	if (client.Jar.Cookies(&url.URL{Scheme: "https", Host: "www.masterclass.com"}) == nil) {
		return fmt.Errorf("cookies not found. Please login first")
	}

	profile, err := getProfile(client, datDir)
	if err != nil {
		return err
	}

	classSlug := ""
	chapterSlug := ""
	// Handle both /chapters/ (classes) and /episodes/ (series) URL patterns
	if strings.Contains(arg, "/chapters/") {
		classSlug = strings.Split(arg, "/chapters/")[0]
		chapterSlug = strings.Split(arg, "/chapters/")[1]
	} else if strings.Contains(arg, "/episodes/") {
		classSlug = strings.Split(arg, "/episodes/")[0]
		chapterSlug = strings.Split(arg, "/episodes/")[1]
	} else {
		classSlug = arg
	}

	// Strip URL prefixes for both classes, sessions, and series
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/sessions/classes/")
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/classes/")
	classSlug = strings.TrimPrefix(classSlug, "https://www.masterclass.com/series/")
	classSlug = strings.TrimPrefix(classSlug, "sessions/classes/")
	classSlug = strings.TrimPrefix(classSlug, "classes/")
	classSlug = strings.TrimPrefix(classSlug, "series/")
	classSlug = strings.TrimSuffix(classSlug, "/")
	chapterSlug = strings.TrimPrefix(chapterSlug, "https://www.masterclass.com/classes/")
	chapterSlug = strings.TrimPrefix(chapterSlug, "https://www.masterclass.com/series/")
	chapterSlug = strings.TrimPrefix(chapterSlug, "classes/")
	chapterSlug = strings.TrimPrefix(chapterSlug, "series/")
	chapterSlug = strings.TrimSuffix(chapterSlug, "/")
	if classSlug == "" {
		return fmt.Errorf("invalid class/series slug")
	}

	//get class info (include=instructors,chapters,categories ensures we get full data for series)
	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/courses/"+classSlug+"?deep=true&include=instructors,chapters,categories", nil)
	req.Header.Set("Referer", "https://www.masterclass.com/classes/"+classSlug)
	req.Header.Set("Mc-Profile-Id", profile.UUID)
	if err != nil {
		return err
	}
	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		// Not a course — try the Sessions (camp) API
		return downloadCamp(client, profile.UUID, outputDir, downloadPdfs, downloadPosters, ytdlExec, nameAsSeries, writeNfo, metadataOnly, forceDownload, concurrency, subsOnly, classSlug, chapterSlug)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var class CourseResponse
	err = json.Unmarshal(bodyBytes, &class)
	if err != nil {
		return err
	}

	// Handle shallow chapters (API returns only {"id": ...} for some courses)
	// In this case, we need to fetch full chapter data individually
	if len(class.Chapters) > 0 && class.Chapters[0].Slug == "" {
		fmt.Println("Detected shallow chapter data, fetching full chapter details...")
		for i, ch := range class.Chapters {
			if ch.ID == 0 {
				continue
			}
			fullChapter, err := fetchChapter(client, profile.UUID, ch.ID)
			if err != nil {
				fmt.Printf("  Warning: failed to fetch chapter %d: %v\n", ch.ID, err)
				continue
			}
			class.Chapters[i] = *fullChapter
		}
	}
	outputDir = path.Join(outputDir, sanitizeFilename(class.Title))
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	fmt.Printf("Output directory: %s\n", outputDir)

	// Download show artwork (Plex naming convention)
	if downloadPosters {
		if class.Primary2x3 != "" {
			fmt.Println("Downloading poster image")
			err = downloadImage(client, class.Primary2x3, path.Join(outputDir, "poster.jpg"))
			if err != nil {
				fmt.Printf("Warning: failed to download poster: %v\n", err)
			}
		}
		if class.Primary16x9 != "" {
			fmt.Println("Downloading fanart image")
			err = downloadImage(client, class.Primary16x9, path.Join(outputDir, "fanart.jpg"))
			if err != nil {
				fmt.Printf("Warning: failed to download fanart: %v\n", err)
			}
		}
	}

	if downloadPdfs && !metadataOnly {
		fmt.Println("Downloading PDFs")
		for _, pdf := range class.AllPDFs {
			pdfURL := pdf.URL
			pdfTitle := pdf.Title

			// If PDF only has ID (common for series), fetch full details
			if pdfURL == "" && pdf.ID > 0 {
				fmt.Printf("Fetching PDF details for ID %d...\n", pdf.ID)
				pdfReq, err := http.NewRequest("GET", fmt.Sprintf("https://www.masterclass.com/jsonapi/v1/pdfs/%d", pdf.ID), nil)
				if err != nil {
					fmt.Printf("Warning: failed to create PDF request: %v\n", err)
					continue
				}
				pdfReq.Header.Set("Referer", "https://www.masterclass.com/classes/"+classSlug)
				pdfReq.Header.Set("Mc-Profile-Id", profile.UUID)
				pdfResp, err := doWithRetry(client, pdfReq, 3)
				if err != nil {
					fmt.Printf("Warning: failed to fetch PDF %d: %v\n", pdf.ID, err)
					continue
				}
				if pdfResp.StatusCode == 200 {
					var pdfDetails struct {
						Title string `json:"title"`
						URL   string `json:"url"`
					}
					err = json.NewDecoder(pdfResp.Body).Decode(&pdfDetails)
					pdfResp.Body.Close()
					if err == nil && pdfDetails.URL != "" {
						pdfURL = pdfDetails.URL
						pdfTitle = pdfDetails.Title
					}
				} else {
					pdfResp.Body.Close()
					fmt.Printf("Warning: PDF %d returned status %d\n", pdf.ID, pdfResp.StatusCode)
					continue
				}
			}

			if pdfURL == "" {
				fmt.Printf("Skipping PDF with empty URL: %s\n", pdfTitle)
				continue
			}
			req, err := http.NewRequest("GET", pdfURL, nil)
			if err != nil {
				return err
			}
			req.Header.Set("Referer", "https://www.masterclass.com/classes/"+classSlug)
			req.Header.Set("Mc-Profile-Id", profile.UUID)
			resp, err := doWithRetry(client, req, 3)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("failed to download PDF")
			}
			fmt.Printf("Downloading: %s\n", pdfTitle)
			safePdfTitle := sanitizeFilename(pdfTitle)
			pdfFile, err := os.Create(path.Join(outputDir, safePdfTitle+".pdf"))
			if err != nil {
				return err
			}
			defer pdfFile.Close()
			_, err = io.Copy(pdfFile, resp.Body)
			if err != nil {
				return err
			}
		}
	}

	downloadedCount := 0
	// Download videos or subtitles (skip if metadataOnly)
	if !metadataOnly {
		// Masterclass uses a fixed API key for media metadata requests
		apiKey := "b9517f7d8d1f48c2de88100f2c13e77a9d8e524aed204651acca65202ff5c6cb9244c045795b1fafda617ac5eb0a6c50"

		if chapterSlug != "" {
			fmt.Printf("Looking for chapter slug: %s\n", chapterSlug)
		}
		fmt.Printf("Found %d chapters:\n", len(class.Chapters))
		for _, ch := range class.Chapters {
			fmt.Printf("  Chapter %d: %s (slug: %s)\n", ch.Number, ch.Title, ch.Slug)
		}

		if subsOnly {
			for _, chapter := range class.Chapters {
				if chapterSlug != "" && chapter.Slug != chapterSlug {
					continue
				}
				fmt.Printf("Downloading chapter %d: %s\n", chapter.Number, chapter.Title)
				downloaded, err := downloadChapterSubsOnly(client, profile.UUID, outputDir, ytdlExec, chapter, apiKey)
				if err != nil {
					return err
				}
				if downloaded {
					downloadedCount++
				}
			}
		} else {
			// Create CycleTLS client once for all chapters (avoid memory leak)
			cycleclient := cycletls.Init()
			defer func() {
				defer func() { recover() }() // Catch panic from Close()
				cycleclient.Close()
			}()

			for _, chapter := range class.Chapters {
				if chapterSlug != "" && chapter.Slug != chapterSlug {
					continue
				}
				fmt.Printf("Downloading chapter %d: %s\n", chapter.Number, chapter.Title)
				downloaded, err := downloadChapterVideo(cycleclient, client, profile.UUID, outputDir, ytdlExec, chapter, class, apiKey, nameAsSeries, writeNfo, forceDownload, concurrency)
				if err != nil {
					return err
				}
				if downloaded {
					downloadedCount++
				}
			}
		}
	}

	// Write NFO metadata file (always write if metadataOnly, otherwise respect writeNfo flag)
	if writeNfo || metadataOnly {
		fmt.Println("Writing tvshow.nfo")
		err = writeNFO(class, outputDir)
		if err != nil {
			fmt.Printf("Warning: failed to write NFO: %v\n", err)
		}

		// Write episode NFOs when in metadata-only mode
		if metadataOnly {
			fmt.Println("Writing episode NFO files")
			for _, chapter := range class.Chapters {
				if chapterSlug != "" && chapter.Slug != chapterSlug {
					continue
				}
				// Generate filename matching video naming convention
				var baseFileName string
				if nameAsSeries {
					if chapter.IsExampleLesson {
						baseFileName = fmt.Sprintf("s01e%02d-%s-Extra_trailer", chapter.Number, sanitizeFilename(chapter.Title))
					} else {
						baseFileName = fmt.Sprintf("s01e%02d-%s", chapter.Number, sanitizeFilename(chapter.Title))
					}
				} else {
					baseFileName = fmt.Sprintf("%03d-%s", chapter.Number, sanitizeFilename(chapter.Title))
				}
				nfoFilename := baseFileName + ".nfo"
				err = writeEpisodeNFO(chapter, class, outputDir, nfoFilename)
				if err != nil {
					fmt.Printf("Warning: failed to write episode NFO for %s: %v\n", chapter.Title, err)
				}
			}
		}
	}

	if !metadataOnly {
		if subsOnly {
			fmt.Printf("Done - %d subtitle(s) downloaded successfully\n", downloadedCount)
		} else {
			fmt.Printf("Done - %d chapter(s) downloaded successfully\n", downloadedCount)
		}
	}

	return nil
}

func downloadCategory(client *http.Client, datDir string, outputDir string, downloadPdfs bool, downloadPosters bool, ytdlExec string, limit int, nameAsSeries bool, writeNfo bool, metadataOnly bool, forceDownload bool, concurrency int, subsOnly bool, arg string) error {
	if (client.Jar.Cookies(&url.URL{Scheme: "https", Host: "www.masterclass.com"}) == nil) {
		return fmt.Errorf("cookies not found. Please login first")
	}

	profile, err := getProfile(client, datDir)
	if err != nil {
		return err
	}

	// Parse the category URL to extract bundle name
	// Input: https://www.masterclass.com/homepage/science-and-tech
	// Bundle: homepage-science-and-tech
	categorySlug := arg
	categorySlug = strings.TrimPrefix(categorySlug, "https://www.masterclass.com/")
	categorySlug = strings.TrimPrefix(categorySlug, "http://www.masterclass.com/")
	categorySlug = strings.TrimSuffix(categorySlug, "/")

	// Convert path to bundle format: homepage/science-and-tech -> homepage-science-and-tech
	bundle := strings.ReplaceAll(categorySlug, "/", "-")

	fmt.Printf("Fetching category: %s (bundle: %s)\n", categorySlug, bundle)

	// Call the content-rows API
	apiURL := fmt.Sprintf("https://www.masterclass.com/jsonapi/v3/content-rows?filter[platform]=web&filter[bundle]=%s&include_items=true", bundle)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.masterclass.com/"+categorySlug)
	req.Header.Set("Mc-Profile-Id", profile.UUID)

	resp, err := doWithRetry(client, req, 3)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get category info: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var contentRows ContentRowsResponse
	err = json.NewDecoder(resp.Body).Decode(&contentRows)
	if err != nil {
		return fmt.Errorf("failed to parse content rows: %v", err)
	}

	// Extract all unique courses from all rows
	type CourseInfo struct {
		Slug     string
		Title    string
		Subtitle string
		Duration string
	}
	courseMap := make(map[string]CourseInfo)
	for _, row := range contentRows {
		for _, item := range row.Items {
			// Only include courses (not series or other content types)
			resource := item.Default.Resource
			if resource.EntitySlug != "" && resource.EntityType == "course" {
				courseMap[resource.EntitySlug] = CourseInfo{
					Slug:     resource.EntitySlug,
					Title:    item.Default.Title,
					Subtitle: item.Default.Subtitle,
					Duration: item.Default.Duration,
				}
			}
		}
	}

	// Convert to slice
	var courses []CourseInfo
	for _, info := range courseMap {
		courses = append(courses, info)
	}

	fmt.Printf("\nFound %d courses in category '%s':\n", len(courses), categorySlug)
	fmt.Println(strings.Repeat("-", 60))

	// Print all courses
	for i, course := range courses {
		duration := ""
		if course.Duration != "" {
			duration = fmt.Sprintf(" (%s)", course.Duration)
		}
		subtitle := ""
		if course.Subtitle != "" {
			subtitle = fmt.Sprintf(" - %s", course.Subtitle)
		}
		fmt.Printf("%3d. %s%s%s\n", i+1, course.Title, subtitle, duration)
	}
	fmt.Println(strings.Repeat("-", 60))

	// Apply limit
	downloadCount := len(courses)
	if limit > 0 && limit < downloadCount {
		downloadCount = limit
		fmt.Printf("\nDownloading first %d of %d courses (use --limit 0 for all):\n\n", downloadCount, len(courses))
	} else {
		fmt.Printf("\nDownloading all %d courses:\n\n", downloadCount)
	}

	// Download each course
	for i := 0; i < downloadCount; i++ {
		course := courses[i]
		fmt.Printf("\n[%d/%d] Downloading: %s\n", i+1, downloadCount, course.Title)
		fmt.Println(strings.Repeat("=", 60))

		err := download(client, datDir, outputDir, downloadPdfs, downloadPosters, ytdlExec, nameAsSeries, writeNfo, metadataOnly, forceDownload, concurrency, subsOnly, course.Slug)
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", course.Slug, err)
			// Continue with next course instead of stopping
			continue
		}
	}

	fmt.Printf("\n\nCategory download complete! Downloaded %d courses.\n", downloadCount)
	return nil
}

// getChapterStreamInfo fetches the stream URL and text tracks for a chapter
func getChapterStreamInfo(client *http.Client, profileUUID string, mediaUUID string, apiKey string) (string, []TextTrack, error) {
	// Use CycleTLS for the media metadata API request to bypass any Cloudflare protection
	cycleclient := cycletls.Init()
	defer func() {
		defer func() { recover() }() // Catch panic from Close()
		cycleclient.Close()
	}()

	// Build cookie string from jar
	wwwURL, _ := url.Parse("https://www.masterclass.com")
	edgeURL, _ := url.Parse("https://edge.masterclass.com")

	wwwCookies := client.Jar.Cookies(wwwURL)
	edgeCookies := client.Jar.Cookies(edgeURL)

	cookieMap := make(map[string]string)
	for _, c := range edgeCookies {
		cookieMap[c.Name] = c.Value
	}
	for _, c := range wwwCookies {
		cookieMap[c.Name] = c.Value
	}

	var cookieStr string
	first := true
	for name, value := range cookieMap {
		if !first {
			cookieStr += "; "
		}
		cookieStr += name + "=" + value
		first = false
	}

	metadataResp, err := cycleclient.Do("https://edge.masterclass.com/api/v1/media/metadata/"+mediaUUID, cycletls.Options{
		Body:      "",
		Ja3:       ja3,
		UserAgent: userAgent,
		Headers: map[string]string{
			"Accept":             "application/json",
			"Accept-Language":    "en-US,en;q=0.9",
			"Content-Type":       "application/json",
			"Origin":             "https://www.masterclass.com",
			"Referer":            "https://www.masterclass.com/",
			"Mc-Profile-Id":      profileUUID,
			"X-Api-Key":          apiKey,
			"Cookie":             cookieStr,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-site",
			"Sec-Ch-Ua":          `"Chromium";v="141", "Not?A_Brand";v="8"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"macOS"`,
		},
	}, "GET")

	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}

	if metadataResp.Status != 200 {
		return "", nil, fmt.Errorf("failed to get chapter metadata: status=%d", metadataResp.Status)
	}

	var chapterMetadata ChapterMetadataResponse
	err = json.Unmarshal([]byte(metadataResp.Body), &chapterMetadata)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse metadata: %v", err)
	}

	var streamURL string
	if len(chapterMetadata.Sources) > 0 {
		streamURL = chapterMetadata.Sources[0].Src
	}

	return streamURL, chapterMetadata.TextTracks, nil
}

func downloadChapterSubsOnly(client *http.Client, profileUUID string, outputDir string, ytdlExec string, chapter Chapter, apiKey string) (bool, error) {
	// Skip chapters without video content (e.g., PDF-only chapters)
	if chapter.MediaUUID == "" {
		fmt.Printf("Skipping chapter %d: %s (no video content)\n", chapter.Number, chapter.Title)
		return false, nil
	}

	safeTitle := sanitizeFilename(chapter.Title)
	baseFilename := path.Join(outputDir, fmt.Sprintf("%03d-%s", chapter.Number, safeTitle))

	// Get stream info from metadata API
	streamURL, textTracks, err := getChapterStreamInfo(client, profileUUID, chapter.MediaUUID, apiKey)
	if err != nil {
		fmt.Printf("  Warning: failed to get stream info: %v\n", err)
	}

	// Prefer TextTracks from metadata API, fall back to chapter data
	tracks := textTracks
	if len(tracks) == 0 {
		tracks = chapter.TextTracks
	}

	// First, try direct TextTracks download (faster than yt-dlp)
	downloadedSubs := 0
	if len(tracks) > 0 {
		for _, track := range tracks {
			if track.Src == "" {
				continue
			}
			// Create filename with language code
			subFilename := fmt.Sprintf("%s.%s.vtt", baseFilename, track.SrcLang)

			// Download the VTT file directly
			resp, err := client.Get(track.Src)
			if err != nil {
				fmt.Printf("  Warning: failed to download %s subtitle: %v\n", track.Label, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Printf("  Warning: failed to download %s subtitle: status %d\n", track.Label, resp.StatusCode)
				continue
			}

			subFile, err := os.Create(subFilename)
			if err != nil {
				fmt.Printf("  Warning: failed to create subtitle file for %s: %v\n", track.Label, err)
				continue
			}

			_, err = io.Copy(subFile, resp.Body)
			subFile.Close()
			if err != nil {
				fmt.Printf("  Warning: failed to write subtitle file for %s: %v\n", track.Label, err)
				continue
			}

			fmt.Printf("  Downloaded subtitle: %s (%s)\n", track.Label, track.SrcLang)
			downloadedSubs++
		}
	}

	if downloadedSubs > 0 {
		return true, nil
	}

	// Fallback: try yt-dlp subtitle extraction (works if subs are in HLS manifest)
	if streamURL != "" {
		fmt.Printf("  TextTracks empty, trying yt-dlp fallback...\n")
		cmd := exec.Command(ytdlExec, "--skip-download", "--write-subs", "--all-subs", "-o", baseFilename+".%(ext)s", streamURL)
		// yt-dlp can crash on some URLs due to regex bugs, so treat errors as non-fatal
		if err := cmd.Run(); err != nil {
			fmt.Printf("  Warning: yt-dlp subtitle extraction failed: %v\n", err)
		}

		// Check if yt-dlp produced any subtitle files
		entries, err := os.ReadDir(outputDir)
		if err != nil {
			fmt.Printf("  Warning: failed to read output directory: %v\n", err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, fmt.Sprintf("%03d-%s", chapter.Number, safeTitle)) &&
				(strings.HasSuffix(name, ".vtt") || strings.HasSuffix(name, ".srt") || strings.HasSuffix(name, ".ass")) {
				fmt.Printf("  Downloaded subtitle via yt-dlp: %s\n", name)
				downloadedSubs++
			}
		}
	}

	if downloadedSubs == 0 {
		fmt.Printf("Skipping chapter %d: %s (no subtitles available from any source)\n", chapter.Number, chapter.Title)
		return false, nil
	}
	return true, nil
}

func downloadChapterVideo(cycleclient cycletls.CycleTLS, client *http.Client, profileUUID string, outputDir string, ytdlExec string, chapter Chapter, course CourseResponse, apiKey string, nameAsSeries bool, writeNfo bool, forceDownload bool, concurrency int) (bool, error) {
	// Skip chapters without video content (e.g., PDF-only chapters)
	if chapter.MediaUUID == "" {
		fmt.Printf("Skipping chapter %d: %s (no video content)\n", chapter.Number, chapter.Title)
		return false, nil
	}

	// Build cookie string from jar - try getting from www.masterclass.com
	wwwURL, _ := url.Parse("https://www.masterclass.com")
	edgeURL, _ := url.Parse("https://edge.masterclass.com")

	// Get cookies from both URLs and merge them
	wwwCookies := client.Jar.Cookies(wwwURL)
	edgeCookies := client.Jar.Cookies(edgeURL)

	// Build a map to collect unique cookies, preferring www cookies
	cookieMap := make(map[string]string)
	for _, c := range edgeCookies {
		cookieMap[c.Name] = c.Value
	}
	for _, c := range wwwCookies {
		cookieMap[c.Name] = c.Value // Overwrite with www value if exists
	}

	var cookieStr string
	first := true
	for name, value := range cookieMap {
		if !first {
			cookieStr += "; "
		}
		cookieStr += name + "=" + value
		first = false
	}

	metadataResp, err := cycleclient.Do("https://edge.masterclass.com/api/v1/media/metadata/"+chapter.MediaUUID, cycletls.Options{
		Body:      "",
		Ja3:       ja3,
		UserAgent: userAgent,
		Headers: map[string]string{
			"Accept":             "application/json",
			"Accept-Language":    "en-US,en;q=0.9",
			"Content-Type":       "application/json",
			"Origin":             "https://www.masterclass.com",
			"Referer":            "https://www.masterclass.com/",
			"Mc-Profile-Id":      profileUUID,
			"X-Api-Key":          apiKey,
			"Cookie":             cookieStr,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-site",
			"Sec-Ch-Ua":          `"Chromium";v="141", "Not?A_Brand";v="8"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"macOS"`,
		},
	}, "GET")

	if err != nil {
		return false, fmt.Errorf("failed to fetch metadata: %v", err)
	}

	if metadataResp.Status != 200 {
		fmt.Printf("Response status: %d\n", metadataResp.Status)
		previewLen := min(len(metadataResp.Body), 500)
		if previewLen > 0 {
			fmt.Printf("Response body: %s\n", metadataResp.Body[:previewLen])
		}
		return false, fmt.Errorf("failed to get chapter metadata: status=%d", metadataResp.Status)
	}

	var chapterMetadata ChapterMetadataResponse
	err = json.Unmarshal([]byte(metadataResp.Body), &chapterMetadata)
	if err != nil {
		return false, fmt.Errorf("failed to parse metadata: %v", err)
	}

	// Check if there are video sources available
	if len(chapterMetadata.Sources) == 0 {
		fmt.Printf("Skipping chapter %d: %s (no video sources)\n", chapter.Number, chapter.Title)
		return false, nil
	}

	// Generate filename based on naming mode
	var baseFileName string
	if nameAsSeries {
		// TV series format: s01e01-Title.mp4 or s01e01-Title-Extra_trailer.mp4
		if chapter.IsExampleLesson {
			baseFileName = fmt.Sprintf("s01e%02d-%s-Extra_trailer", chapter.Number, sanitizeFilename(chapter.Title))
		} else {
			baseFileName = fmt.Sprintf("s01e%02d-%s", chapter.Number, sanitizeFilename(chapter.Title))
		}
	} else {
		// Default format: 001-Title.mp4
		baseFileName = fmt.Sprintf("%03d-%s", chapter.Number, sanitizeFilename(chapter.Title))
	}
	outputFile := path.Join(outputDir, baseFileName+".mp4")

	// Check if output file already exists (skip unless --force)
	if !forceDownload {
		if _, err := os.Stat(outputFile); err == nil {
			fmt.Printf("Skipping %s (already exists)\n", baseFileName+".mp4")
			return false, nil
		}
	}

	// Build metadata arguments - always embed full metadata regardless of naming mode
	// Extract date (YYYY-MM-DD) from UpdatedAt
	dateStr := ""
	if chapter.UpdatedAt != "" && len(chapter.UpdatedAt) >= 10 {
		dateStr = chapter.UpdatedAt[:10] // "2024-03-20T..." -> "2024-03-20"
	}

	// Build genre/tags from all categories
	genre := "Education"
	if len(course.Categories) > 0 {
		var genres []string
		for _, cat := range course.Categories {
			genres = append(genres, cat.Name)
		}
		genre = strings.Join(genres, ", ")
	}

	// Generate episode_id
	episodeID := fmt.Sprintf("s01e%02d", chapter.Number)

	// Full metadata for all downloads
	// Include multiple description tags for compatibility (description, comment, long_description, synopsis)
	metadataArgs := fmt.Sprintf(
		"ffmpeg:-metadata title=%q -metadata show=%q -metadata artist=%q -metadata genre=%q -metadata date=%q -metadata description=%q -metadata comment=%q -metadata long_description=%q -metadata synopsis=%q -metadata season_number=1 -metadata episode_sort=%d -metadata episode_id=%q -metadata network=%q",
		chapter.Title,
		course.Title,
		course.InstructorName,
		genre,
		dateStr,
		chapter.Abstract,
		chapter.Abstract,
		chapter.Abstract,
		course.Overview,
		chapter.Number,
		episodeID,
		"MasterClass",
	)

	// Build yt-dlp command with metadata embedding
	args := []string{
		"--continue",
		"--embed-subs", "--all-subs",
		"--embed-metadata",
		"-f", "bestvideo+bestaudio",
		"-N", fmt.Sprintf("%d", concurrency),
		"--postprocessor-args", metadataArgs,
		chapterMetadata.Sources[0].Src,
		"-o", outputFile,
	}

	cmd := exec.Command(ytdlExec, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return false, err
	}

	// Write episode NFO if requested
	if writeNfo {
		nfoFilename := baseFileName + ".nfo"
		err = writeEpisodeNFO(chapter, course, outputDir, nfoFilename)
		if err != nil {
			fmt.Printf("Warning: failed to write episode NFO: %v\n", err)
		}
	}

	return true, nil
}



func downloadImage(client *http.Client, imageURL string, outputPath string) error {
	resp, err := client.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download image: status=%d", resp.StatusCode)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// NFO XML structures for Kodi/Plex/Jellyfin compatibility
type TVShowNFO struct {
	XMLName   xml.Name    `xml:"tvshow"`
	Title     string      `xml:"title"`
	Plot      string      `xml:"plot"`
	Outline   string      `xml:"outline,omitempty"`
	Tagline   string      `xml:"tagline,omitempty"`
	Genres    []string    `xml:"genre"`
	Tags      []string    `xml:"tag,omitempty"`
	Studio    string      `xml:"studio"`
	Premiered string      `xml:"premiered,omitempty"`
	Runtime   int         `xml:"runtime,omitempty"`
	Actors    []NFOActor  `xml:"actor"`
	Thumbs    []NFOThumb  `xml:"thumb"`
	UniqueID  NFOUniqueID `xml:"uniqueid"`
}

type NFOActor struct {
	Name  string `xml:"name"`
	Role  string `xml:"role"`
	Thumb string `xml:"thumb,omitempty"`
	Bio   string `xml:"biography,omitempty"`
}

type NFOThumb struct {
	Aspect string `xml:"aspect,attr"`
	Value  string `xml:",chardata"`
}

type NFOUniqueID struct {
	Type    string `xml:"type,attr"`
	Default bool   `xml:"default,attr"`
	Value   string `xml:",chardata"`
}

// EpisodeNFO represents the episodedetails.nfo format for Kodi/Plex/Jellyfin
type EpisodeNFO struct {
	XMLName   xml.Name    `xml:"episodedetails"`
	Title     string      `xml:"title"`
	ShowTitle string      `xml:"showtitle,omitempty"`
	Season    int         `xml:"season"`
	Episode   int         `xml:"episode"`
	Plot      string      `xml:"plot"`
	Aired     string      `xml:"aired,omitempty"`
	Runtime   int         `xml:"runtime,omitempty"`
	Actors    []NFOActor  `xml:"actor"`
	Studio    string      `xml:"studio,omitempty"`
	UniqueID  NFOUniqueID `xml:"uniqueid"`
}

// splitInstructorNames splits an instructor string into individual names.
// Handles patterns like:
//   - "Kim Kardashian" -> ["Kim Kardashian"]
//   - "Mike Cessario and Laura Modi" -> ["Mike Cessario", "Laura Modi"]
//   - "Jeff Goodby & Rich Silverstein" -> ["Jeff Goodby", "Rich Silverstein"]
//   - "Rich Paul, Bob Myers, and Draymond Green" -> ["Rich Paul", "Bob Myers", "Draymond Green"]
func splitInstructorNames(instructorStr string) []string {
	var names []string

	// First split by comma
	parts := strings.Split(instructorStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Remove leading "and " if present (from "A, B, and C" pattern)
		part = strings.TrimPrefix(part, "and ")
		part = strings.TrimSpace(part)

		// Check for " and " within the part
		if strings.Contains(part, " and ") {
			subParts := strings.Split(part, " and ")
			for _, sp := range subParts {
				sp = strings.TrimSpace(sp)
				if sp != "" {
					names = append(names, sp)
				}
			}
			continue
		}

		// Check for " & " within the part
		if strings.Contains(part, " & ") {
			subParts := strings.Split(part, " & ")
			for _, sp := range subParts {
				sp = strings.TrimSpace(sp)
				if sp != "" {
					names = append(names, sp)
				}
			}
			continue
		}

		// Single name
		if part != "" {
			names = append(names, part)
		}
	}

	return names
}

func writeNFO(course CourseResponse, outputDir string) error {
	nfoPath := path.Join(outputDir, "tvshow.nfo")

	// Extract premiered date (YYYY-MM-DD) from UpdatedAt
	premiered := ""
	if course.UpdatedAt != "" && len(course.UpdatedAt) >= 10 {
		premiered = course.UpdatedAt[:10]
	}

	// Build genres from categories (skip empty names - series only have IDs)
	var genres []string
	for _, cat := range course.Categories {
		if cat.Name != "" {
			genres = append(genres, cat.Name)
		}
	}

	// Note: <tag> in NFO becomes Collections in Plex (via XBMCnfoTVImporter)
	// Use primary category as the tag (e.g., "Business") for collection grouping
	var tags []string
	if course.PrimaryCategory.ID != 0 {
		// Find the primary category name from categories array
		for _, cat := range course.Categories {
			if cat.ID == course.PrimaryCategory.ID && cat.Name != "" {
				tags = append(tags, cat.Name)
				break
			}
		}
	}

	// Build actor list from instructors array (if names populated) or fall back to splitting instructor_name
	var actors []NFOActor
	hasValidInstructors := false
	if len(course.Instructors) > 0 && course.Instructors[0].Name != "" {
		hasValidInstructors = true
		numInstructors := len(course.Instructors)
		for i, inst := range course.Instructors {
			role := "Instructor"
			if numInstructors > 1 {
				role = fmt.Sprintf("Instructor %d", i+1)
			}
			actor := NFOActor{
				Name: inst.Name,
				Role: role,
			}
			// Add bio if available
			if inst.Bio != nil {
				if bio, ok := inst.Bio.(string); ok && bio != "" {
					actor.Bio = bio
				}
			}
			// Use headshot URL for actor thumb
			if inst.HeadshotURL != nil {
				if headshot, ok := inst.HeadshotURL.(string); ok && headshot != "" {
					actor.Thumb = headshot + "?width=500&height=500&fit=cover&dpr=2"
				}
			}
			// Fallback to poster if no headshot
			if actor.Thumb == "" && course.Primary2x3 != "" {
				actor.Thumb = course.Primary2x3 + "?width=500&height=500&fit=cover&dpr=2"
			}
			actors = append(actors, actor)
		}
	}
	if !hasValidInstructors {
		// Fallback: split instructor_name string
		instructorNames := splitInstructorNames(course.InstructorName)
		numInstructors := len(instructorNames)
		for i, name := range instructorNames {
			role := "Instructor"
			if numInstructors > 1 {
				role = fmt.Sprintf("Instructor %d", i+1)
			}
			actor := NFOActor{
				Name: name,
				Role: role,
			}
			// Use poster as fallback thumb URL
			if course.Primary2x3 != "" {
				actor.Thumb = course.Primary2x3 + "?width=500&height=500&fit=cover&dpr=2"
			}
			actors = append(actors, actor)
		}
	}

	// Build the NFO struct
	nfo := TVShowNFO{
		Title:     course.Title,
		Plot:      course.Overview,
		Outline:   course.ShortOverview,
		Tagline:   course.InstructorTagline,
		Genres:    genres,
		Tags:      tags,
		Studio:    "MasterClass",
		Premiered: premiered,
		Runtime:   course.TotalSeconds / 60,
		Actors:    actors,
		Thumbs: []NFOThumb{
			{Aspect: "poster", Value: "poster.jpg"},
			{Aspect: "fanart", Value: "fanart.jpg"},
		},
		UniqueID: NFOUniqueID{
			Type:    "masterclass",
			Default: true,
			Value:   course.Slug,
		},
	}

	// Marshal to XML with indentation
	output, err := xml.MarshalIndent(nfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal NFO: %v", err)
	}

	// Add XML declaration
	xmlContent := []byte(xml.Header + string(output) + "\n")

	return os.WriteFile(nfoPath, xmlContent, 0644)
}

// writeEpisodeNFO writes an episodedetails.nfo file for a single episode
func writeEpisodeNFO(chapter Chapter, course CourseResponse, outputDir string, nfoFilename string) error {
	nfoPath := path.Join(outputDir, nfoFilename)

	// Extract aired date from chapter's UpdatedAt (YYYY-MM-DD)
	aired := ""
	if chapter.UpdatedAt != "" && len(chapter.UpdatedAt) >= 10 {
		aired = chapter.UpdatedAt[:10]
	}

	// Calculate runtime in minutes from DurationSecs
	runtime := 0
	if chapter.DurationSecs > 0 {
		runtime = chapter.DurationSecs / 60
	}

	// Build actors list with portraits (same logic as tvshow.nfo)
	var actors []NFOActor
	hasValidInstructors := false
	if len(course.Instructors) > 0 && course.Instructors[0].Name != "" {
		hasValidInstructors = true
		numInstructors := len(course.Instructors)
		for i, inst := range course.Instructors {
			role := "Instructor"
			if numInstructors > 1 {
				role = fmt.Sprintf("Instructor %d", i+1)
			}
			actor := NFOActor{
				Name: inst.Name,
				Role: role,
			}
			// Add bio if available
			if inst.Bio != nil {
				if bio, ok := inst.Bio.(string); ok && bio != "" {
					actor.Bio = bio
				}
			}
			// Use headshot URL for actor thumb (Plex link mode)
			if inst.HeadshotURL != nil {
				if headshot, ok := inst.HeadshotURL.(string); ok && headshot != "" {
					actor.Thumb = headshot + "?width=500&height=500&fit=cover&dpr=2"
				}
			}
			// Fallback to poster if no headshot
			if actor.Thumb == "" && course.Primary2x3 != "" {
				actor.Thumb = course.Primary2x3 + "?width=500&height=500&fit=cover&dpr=2"
			}
			actors = append(actors, actor)
		}
	}
	if !hasValidInstructors {
		// Fallback: split instructor_name string
		instructorNames := splitInstructorNames(course.InstructorName)
		numInstructors := len(instructorNames)
		for i, name := range instructorNames {
			role := "Instructor"
			if numInstructors > 1 {
				role = fmt.Sprintf("Instructor %d", i+1)
			}
			actor := NFOActor{
				Name: name,
				Role: role,
			}
			// Use poster as fallback thumb URL
			if course.Primary2x3 != "" {
				actor.Thumb = course.Primary2x3 + "?width=500&height=500&fit=cover&dpr=2"
			}
			actors = append(actors, actor)
		}
	}

	nfo := EpisodeNFO{
		Title:     chapter.Title,
		ShowTitle: course.Title,
		Season:    1,
		Episode:   chapter.Number,
		Plot:      chapter.Abstract,
		Aired:     aired,
		Runtime:   runtime,
		Actors:    actors,
		Studio:    "MasterClass",
		UniqueID: NFOUniqueID{
			Type:    "masterclass",
			Default: true,
			Value:   chapter.Slug,
		},
	}

	// Marshal to XML with indentation
	output, err := xml.MarshalIndent(nfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal episode NFO: %v", err)
	}

	// Add XML declaration
	xmlContent := []byte(xml.Header + string(output) + "\n")

	return os.WriteFile(nfoPath, xmlContent, 0644)
}

// ---------------------------------------------------------------------------
// Safari SSO login (macOS only)
// ---------------------------------------------------------------------------

// macEpochOffset is the number of seconds between 2001-01-01 and 1970-01-01.
// Safari stores cookie expiry as Mac absolute time (seconds since 2001-01-01).
const macEpochOffset = 978307200

type safariCookieEntry struct {
	Domain   string
	Name     string
	Value    string
	Path     string
	Expires  time.Time
	Secure   bool
	HttpOnly bool
}

// safariLogin reads MasterClass cookies from Safari's binary cookie store,
// injects them into the masterclass-dl cookie jar, then fetches the profiles
// API and writes profile.json — equivalent to a normal email/password login.
func safariLogin(client *http.Client, datDir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %v", err)
	}

	cookieDB := path.Join(home,
		"Library/Containers/com.apple.Safari/Data/Library/Cookies/Cookies.binarycookies")

	if _, err := os.Stat(cookieDB); os.IsNotExist(err) {
		return fmt.Errorf("Safari cookie database not found at:\n  %s\n\nMake sure you are logged in to masterclass.com in Safari", cookieDB)
	}

	fmt.Println("Reading Safari cookies...")
	allCookies, err := parseSafariBinaryCookies(cookieDB)
	if err != nil {
		return fmt.Errorf("failed to read Safari cookies: %v\n\n"+
			"You may need to grant Terminal Full Disk Access in:\n"+
			"  System Settings → Privacy & Security → Full Disk Access", err)
	}

	// Filter to masterclass.com cookies only
	var mcCookies []*http.Cookie
	for _, c := range allCookies {
		if strings.HasSuffix(c.Domain, "masterclass.com") {
			mcCookies = append(mcCookies, &http.Cookie{
				Name:     c.Name,
				Value:    c.Value,
				Path:     c.Path,
				Domain:   c.Domain,
				Expires:  c.Expires,
				Secure:   c.Secure,
				HttpOnly: c.HttpOnly,
			})
		}
	}

	if len(mcCookies) == 0 {
		return fmt.Errorf("no MasterClass cookies found in Safari.\nPlease log in to masterclass.com in Safari first")
	}

	// Verify the session cookie is present
	hasSession := false
	for _, c := range mcCookies {
		if c.Name == "_mc_session" {
			hasSession = true
			break
		}
	}
	if !hasSession {
		return fmt.Errorf("_mc_session cookie not found in Safari.\nPlease make sure you are fully logged in to masterclass.com in Safari")
	}

	fmt.Printf("Found %d MasterClass cookies (including _mc_session)\n", len(mcCookies))

	// Inject cookies into the jar for all relevant URLs
	for _, rawURL := range []string{
		"https://www.masterclass.com",
		"https://masterclass.com",
		"https://edge.masterclass.com",
	} {
		u, _ := url.Parse(rawURL)
		client.Jar.SetCookies(u, mcCookies)
	}

	fmt.Println("Fetching profiles...")

	req, err := http.NewRequest("GET", "https://www.masterclass.com/jsonapi/v1/profiles?deep=true", nil)
	if err != nil {
		return fmt.Errorf("failed to create profiles request: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://www.masterclass.com/")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch profiles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("profiles API returned %d: %s\n\n"+
			"Your session may have expired — reload masterclass.com in Safari and try again",
			resp.StatusCode, string(body))
	}

	var profiles []ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return fmt.Errorf("failed to parse profiles response: %v", err)
	}
	if len(profiles) == 0 {
		return fmt.Errorf("no profiles found on your account")
	}

	// Let user pick profile if multiple, otherwise use the first
	var selected ProfileResponse
	if len(profiles) == 1 {
		selected = profiles[0]
	} else {
		prompt := promptui.Select{
			Label: "Select Profile",
			Items: profiles,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ .DisplayName }}",
				Active:   "\U0001F449 {{ .DisplayName }}",
				Inactive: "  {{ .DisplayName }}",
				Selected: "\U0001F64C {{ .DisplayName }}",
			},
		}
		i, _, err := prompt.Run()
		if err != nil {
			return err
		}
		selected = profiles[i]
	}

	fmt.Printf("Selected profile: %s\n", selected.DisplayName)

	profileFile, err := os.Create(path.Join(datDir, "profile.json"))
	if err != nil {
		return fmt.Errorf("cannot create profile.json: %v", err)
	}
	defer profileFile.Close()
	if err := json.NewEncoder(profileFile).Encode(selected); err != nil {
		return fmt.Errorf("cannot write profile.json: %v", err)
	}

	return nil
}

// parseSafariBinaryCookies parses the binary cookies file format used by Safari.
func parseSafariBinaryCookies(filepath string) ([]safariCookieEntry, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	if len(data) < 8 || string(data[:4]) != "cook" {
		return nil, fmt.Errorf("not a valid Safari binary cookies file")
	}

	numPages := int(binary.BigEndian.Uint32(data[4:8]))
	pageOffsets := make([]int, numPages)
	pageSizes := make([]int, numPages)
	headerOffset := 8
	for i := 0; i < numPages; i++ {
		if headerOffset+4 > len(data) {
			break
		}
		pageSizes[i] = int(binary.BigEndian.Uint32(data[headerOffset : headerOffset+4]))
		headerOffset += 4
	}

	cursor := headerOffset
	for i := 0; i < numPages; i++ {
		pageOffsets[i] = cursor
		cursor += pageSizes[i]
	}

	var cookies []safariCookieEntry
	for i, pageSize := range pageSizes {
		if pageOffsets[i]+pageSize > len(data) {
			break
		}
		page := data[pageOffsets[i] : pageOffsets[i]+pageSize]
		if len(page) < 8 {
			continue
		}
		numCookies := int(binary.LittleEndian.Uint32(page[4:8]))
		for j := 0; j < numCookies; j++ {
			start := 8 + j*4
			if start+4 > len(page) {
				break
			}
			co := int(binary.LittleEndian.Uint32(page[start : start+4]))
			c, err := parseSafariCookie(page, co)
			if err != nil {
				continue
			}
			cookies = append(cookies, c)
		}
	}
	return cookies, nil
}

func parseSafariCookie(page []byte, offset int) (safariCookieEntry, error) {
	if offset+56 > len(page) {
		return safariCookieEntry{}, fmt.Errorf("offset out of bounds")
	}
	flags := binary.LittleEndian.Uint32(page[offset+8 : offset+12])
	urlOff := int(binary.LittleEndian.Uint32(page[offset+16 : offset+20]))
	nameOff := int(binary.LittleEndian.Uint32(page[offset+20 : offset+24]))
	pathOff := int(binary.LittleEndian.Uint32(page[offset+24 : offset+28]))
	valueOff := int(binary.LittleEndian.Uint32(page[offset+28 : offset+32]))
	expiryBits := binary.LittleEndian.Uint64(page[offset+40 : offset+48])

	readStr := func(base int) (string, error) {
		abs := offset + base
		if abs >= len(page) {
			return "", fmt.Errorf("string offset out of bounds")
		}
		end := abs
		for end < len(page) && page[end] != 0 {
			end++
		}
		return string(page[abs:end]), nil
	}

	domain, err := readStr(urlOff)
	if err != nil {
		return safariCookieEntry{}, err
	}
	name, err := readStr(nameOff)
	if err != nil {
		return safariCookieEntry{}, err
	}
	cookiePath, err := readStr(pathOff)
	if err != nil {
		return safariCookieEntry{}, err
	}
	value, err := readStr(valueOff)
	if err != nil {
		return safariCookieEntry{}, err
	}

	var expires time.Time
	if expFloat := math.Float64frombits(expiryBits); expFloat > 0 {
		expires = time.Unix(int64(expFloat)+macEpochOffset, 0)
	}

	return safariCookieEntry{
		Domain:   domain,
		Name:     name,
		Value:    value,
		Path:     cookiePath,
		Expires:  expires,
		Secure:   flags&1 != 0,
		HttpOnly: flags&4 != 0,
	}, nil
}
