package exocomp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"
)

type line struct {
	offset  int
	content string
}

const fileWithImages = `^(Dockerfile|.*\.ya?ml)$`
const linesFrom = `^FROM ([^\s]*)`

func ImageDigests(repoPath string) error {

	var filesFilter = regexp.MustCompile(fileWithImages)
	var fromFilter = regexp.MustCompile(linesFrom)

	files, err := findFiles(repoPath, filesFilter)
	if err != nil {
		return fmt.Errorf("failed to search Dockerfiles: %s", err)
	}

	for _, file := range files {
		fmt.Println(file)

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = findLines(f, fromFilter)
		if err != nil {
			return err
		}
	}

	return nil
}

func findFiles(root string, filter *regexp.Regexp) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filter.MatchString(info.Name()) {
			files = append(files, path)
		}
		return err
	})

	return files, err
}

func findLines(r io.Reader, reg *regexp.Regexp) ([]line, error) {

	ll := []line{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if !reg.MatchString(t) {
			continue
		}
		fmt.Println(t)
		m := reg.FindStringSubmatch(t)
		fmt.Println(m[1])
		digest(m[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ll, nil
}

func digest(img string) error {
	image, err := registry.ParseImage(img)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Create the registry client.
	r, err := createRegistryClient(ctx, image.Domain)
	if err != nil {
		return err
	}

	// Get the digest.
	digest, err := r.Digest(ctx, image)
	if err != nil {
		return err
	}

	fmt.Println(digest.String())
	return nil
}

var (
	insecure    bool
	forceNonSSL bool
	skipPing    bool

	timeout = time.Minute

	authURL  string
	username string
	password string

	debug bool
)

func createRegistryClient(ctx context.Context, domain string) (*registry.Registry, error) {
	// Use the auth-url domain if provided.
	authDomain := authURL
	if authDomain == "" {
		authDomain = domain
	}
	auth, err := repoutils.GetAuthConfig(username, password, authDomain)
	if err != nil {
		return nil, err
	}

	// Prevent non-ssl unless explicitly forced
	if !forceNonSSL && strings.HasPrefix(auth.ServerAddress, "http:") {
		return nil, fmt.Errorf("attempted to use insecure protocol! Use force-non-ssl option to force")
	}

	return registry.New(ctx, auth, registry.Opt{
		Domain:   domain,
		Insecure: insecure,
		Debug:    debug,
		SkipPing: skipPing,
		NonSSL:   forceNonSSL,
		Timeout:  timeout,
	})
}
