package est

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ESTOptions struct {
	Server    string
	Label     string
	Operation string
	Username  string
	Password  string
	CSRPath   string
	Output    string
	CN        string
	AllowHTTP bool
}

func RunEST(opts ESTOptions) error {
	if opts.Server == "" || opts.Operation == "" {
		return fmt.Errorf("server and operation are required")
	}

	url := fmt.Sprintf("%s/.well-known/est/%s/%s", strings.TrimRight(opts.Server, "/"), opts.Label, opts.Operation)

	var req *http.Request
	var err error

	switch opts.Operation {
	case "cacerts":
		req, err = http.NewRequest("GET", url, nil)
	case "simpleenroll", "simplereenroll":
		csrBytes, readErr := os.ReadFile(opts.CSRPath)
		if readErr != nil {
			return fmt.Errorf("failed to read CSR: %v", readErr)
		}
		req, err = http.NewRequest("POST", url, bytes.NewReader(csrBytes))
		req.Header.Set("Content-Type", "application/pkcs10")
	default:
		return fmt.Errorf("unsupported operation: %s", opts.Operation)
	}

	if err != nil {
		return err
	}

	req.SetBasicAuth(opts.Username, opts.Password)

	client := &http.Client{}
	if strings.HasPrefix(opts.Server, "https://") && opts.AllowHTTP {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // WARNING: not safe for prod
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		rawBody, err := io.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("EST server error: %s - %s", resp.Status, rawBody)
		} else {
			return fmt.Errorf("EST server error: %s", resp.Status)
		}
	}

	var outFile string
	if opts.Output != "" {
		outFile = opts.Output
	} else {
		outFile = defaultFilenameFromHeaders(resp.Header, opts.Operation)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	body, err := DecodeESTBody(rawBody, resp.Header)
	if err != nil {
		return fmt.Errorf("failed to decode EST body: %v", err)
	}

	if err := os.WriteFile(outFile, body, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	if strings.HasSuffix(outFile, ".p7b") || strings.HasSuffix(outFile, ".der") {
		parsedPEM := strings.TrimSuffix(outFile, filepath.Ext(outFile)) + ".pem"
		if err := SaveCertAsPEM(body, parsedPEM); err == nil {
			fmt.Printf("📜 Certificate parsed and saved as: %s\n", parsedPEM)
		}
	}

	fmt.Printf("✅ Wrote %s (%d bytes)\n", outFile, len(body))
	return nil
}

func defaultFilenameFromHeaders(h http.Header, operation string) string {
	cd := h.Get("Content-Disposition")
	if cd != "" && strings.Contains(cd, "filename=") {
		if parts := strings.Split(cd, "filename="); len(parts) > 1 {
			name := strings.Trim(parts[1], "\"")
			return name
		}
	}
	switch operation {
	case "cacerts":
		return "cacerts.p7b"
	case "simpleenroll", "simplereenroll":
		return "cert.p7b"
	default:
		return "out.p7b"
	}
}
