package exocomp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"

	//"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	//"github.com/docker/docker/builder/dockerfile/parser"
)

func ImageDigests(repoPath string) error {

	files, err := findFiles(repoPath, "Dockerfile")
	if err != nil {
		return fmt.Errorf("failed to search Dockerfiles: %s", err)
	}

	for _, file := range files {
		fmt.Println(file)
		ParseFile(file)
	}

	return nil
}

func findFiles(root, fileName string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.Name() == fileName {
			files = append(files, path)
		}
		return err
	})

	return files, err
}

// Parse a Dockerfile from a filename.  An IOError or ParseError may occur.
func ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return ParseReader(file)
}

// Parse a Dockerfile from a reader.  A ParseError may occur.
func ParseReader(file io.Reader) error {

	res, err := parser.Parse(file)
	if err != nil {
		return err
	}

	for _, child := range res.AST.Children {

		if child.Value == "from" {
			fmt.Println(child.StartLine)
			fmt.Println(child.Original)
			img := strings.Split(child.Original, " ")[1]
			err := digest(img)
			if err != nil {
				return err
			}
		}

		// Only happens for ONBUILD
		if child.Next != nil && len(child.Next.Children) > 0 {
			child = child.Next.Children[0]
		}

	}
	return nil
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
