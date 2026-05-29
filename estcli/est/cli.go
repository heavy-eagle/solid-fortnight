package est

import (
	"encoding/pem"
	"flag"
	"fmt"
	"os"
)

func ExecuteCLI(args []string) error {
	fs := flag.NewFlagSet("estcli", flag.ExitOnError)

	var (
		server    = fs.String("server", "", "EST server base URL (e.g. https://example.com)")
		label     = fs.String("label", "default", "EST label")
		operation = fs.String("operation", "", "Operation: cacerts, simpleenroll, simplereenroll")
		username  = fs.String("user", "", "Username for basic auth")
		password  = fs.String("pass", "", "Password for basic auth")
		csrPath   = fs.String("csr", "", "Path to CSR file (for enrollment)")
		output    = fs.String("out", "", "Output file path (optional)")
		allowHTTP = fs.Bool("allow-http", false, "Allow plain HTTP (for local/dev testing)")
		generate  = fs.Bool("gen", false, "Generate key + CSR automatically")
		cn        = fs.String("cn", "localhost", "Common Name for generated CSR")
		dns       = fs.String("dns", "", "DNS SAN for generated CSR")
	)

	fs.Parse(args[1:]) // skip program name

	opts := ESTOptions{
		Server:    *server,
		Label:     *label,
		Operation: *operation,
		Username:  *username,
		Password:  *password,
		CSRPath:   *csrPath,
		Output:    *output,
		AllowHTTP: *allowHTTP,
	}

	if *generate {
		csrDER, priv, err := GenerateCSRInMemory(*cn, *dns)
		if err != nil {
			return err
		}

		csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
		csrPath := *cn + ".csr"
		keyPath := *cn + ".key"

		// Save files
		if err := os.WriteFile(csrPath, csrPEM, 0600); err != nil {
			return fmt.Errorf("failed to write CSR: %v", err)
		}

		keyPEM, _ := encodePrivateKey(priv)
		if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
			return fmt.Errorf("failed to write key: %v", err)
		}

		fmt.Printf("🔐 Generated key: %s\n📄 CSR written to: %s\n", keyPath, csrPath)
		opts.CSRPath = csrPath
	}

	return RunEST(opts)
}
