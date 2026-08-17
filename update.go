package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// appVersion is the version of this build. Override at compile time with:
//
//	go build -ldflags "-X main.appVersion=1.2.4"
var appVersion = "1.2.4"

const (
	githubOwner = "RykivSale"
	githubRepo  = "faceit-voicechat"
)

var (
	githubLatestURL = "https://github.com/" + githubOwner + "/" + githubRepo + "/releases/latest"
	githubAPIURL    = "https://api.github.com/repos/" + githubOwner + "/" + githubRepo + "/releases/latest"

	updateHTTPTimeout = 6 * time.Second
	checkUpdatesFn    = doUpdateCheck
)

type updateCheck struct {
	Current string
	Latest  string
	URL     string
	Newer   bool
	Err     error
}

type updateView struct {
	Current string `json:"current"`
	Latest  string `json:"latest"`
	URL     string `json:"url"`
	Newer   bool   `json:"newer"`
	Error   string `json:"error,omitempty"`
}

func toUpdateView(check updateCheck) updateView {
	v := updateView{
		Current: check.Current,
		Latest:  check.Latest,
		URL:     check.URL,
		Newer:   check.Newer,
	}
	if check.Err != nil {
		v.Error = check.Err.Error()
		v.Newer = false
	}
	return v
}

func startUpdateCheck() <-chan updateCheck {
	ch := make(chan updateCheck, 1)
	go func() {
		ch <- doUpdateCheck()
	}()
	return ch
}

func doUpdateCheck() updateCheck {
	res := updateCheck{
		Current: appVersion,
		URL:     githubLatestURL,
	}
	latest, err := fetchLatestVersion()
	if err != nil {
		res.Err = err
		return res
	}
	res.Latest = latest
	res.Newer = versionGreater(latest, appVersion)
	return res
}

func fetchLatestVersion() (string, error) {
	if v, err := fetchLatestFromRedirect(); err == nil && v != "" {
		return v, nil
	}
	return fetchLatestFromAPI()
}

func fetchLatestFromRedirect() (string, error) {
	client := &http.Client{
		Timeout: updateHTTPTimeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, githubLatestURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "faceit-voicechat/"+appVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, _ = io.CopyN(io.Discard, resp.Body, 1024)

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("no redirect from releases/latest")
	}
	return tagFromReleaseURL(loc)
}

func fetchLatestFromAPI() (string, error) {
	client := &http.Client{Timeout: updateHTTPTimeout}

	req, err := http.NewRequest(http.MethodGet, githubAPIURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "faceit-voicechat/"+appVersion)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api: %s", resp.Status)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.TagName == "" {
		return "", fmt.Errorf("github api: empty tag_name")
	}
	return normalizeVersion(payload.TagName), nil
}

func tagFromReleaseURL(u string) (string, error) {
	const marker = "/releases/tag/"
	i := strings.LastIndex(u, marker)
	if i < 0 {
		return "", fmt.Errorf("unexpected release URL: %s", u)
	}
	tag := strings.TrimSpace(u[i+len(marker):])
	if cut := strings.IndexAny(tag, "?#"); cut >= 0 {
		tag = tag[:cut]
	}
	tag = strings.Trim(tag, "/")
	if tag == "" || strings.EqualFold(tag, "latest") {
		return "", fmt.Errorf("unexpected release tag in URL: %s", u)
	}
	return normalizeVersion(tag), nil
}

func normalizeVersion(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimPrefix(s, "V")
	return s
}

func versionGreater(remote, local string) bool {
	rMaj, rMin, rPat, rok := parseVersion(remote)
	lMaj, lMin, lPat, lok := parseVersion(local)
	if !rok || !lok {
		return rok && normalizeVersion(remote) != normalizeVersion(local)
	}
	if rMaj != lMaj {
		return rMaj > lMaj
	}
	if rMin != lMin {
		return rMin > lMin
	}
	return rPat > lPat
}

func parseVersion(s string) (major, minor, patch int, ok bool) {
	s = normalizeVersion(s)
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) == 0 || len(parts) > 3 {
		return 0, 0, 0, false
	}
	nums := [3]int{}
	for i, p := range parts {
		if p == "" {
			return 0, 0, 0, false
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return 0, 0, 0, false
		}
		nums[i] = n
	}
	return nums[0], nums[1], nums[2], true
}

func printUpdateBanner(check updateCheck) {
	if check.Err != nil || !check.Newer {
		return
	}
	fmt.Println("--- Update available ---")
	fmt.Printf("A newer version is out: v%s  (you have v%s)\n", check.Latest, check.Current)
	fmt.Printf("Download: %s\n", check.URL)
	fmt.Println("Press U to open the download page.")
	fmt.Println()
}

func runManualUpdateCheck(reader *bufio.Reader) updateCheck {
	fmt.Println()
	fmt.Println("Checking for updates...")
	check := doUpdateCheck()
	if check.Err != nil {
		fmt.Println("Could not check for updates:", check.Err)
		fmt.Println("This does not affect the rest of the program.")
		return check
	}
	if !check.Newer {
		fmt.Printf("You already have the latest version: v%s\n", check.Current)
		return check
	}

	fmt.Println()
	fmt.Println("--- Update available ---")
	fmt.Printf("A newer version is out: v%s  (you have v%s)\n", check.Latest, check.Current)
	fmt.Printf("Download: %s\n", check.URL)
	askOpenDownload(reader, check.URL)
	return check
}

func askOpenDownload(reader *bufio.Reader, rawURL string) {
	fmt.Print("Open the download page? (Y/N) [Y]: ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line != "" && !strings.EqualFold(line, "y") && !strings.EqualFold(line, "yes") {
		return
	}
	if err := openURL(rawURL); err != nil {
		fmt.Println("Could not open the browser. Open this URL manually:")
		fmt.Println(rawURL)
		return
	}
	fmt.Println("Opened the download page in your browser.")
}

func openURL(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}

func onOff(v bool) string {
	if v {
		return "on"
	}
	return "off"
}
