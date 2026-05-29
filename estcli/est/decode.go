package est

import (
    "bytes"
    "compress/gzip"
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "strings"
)

func DecodeESTBody(body []byte, headers http.Header) ([]byte, error) {
    var err error

    // Step 1: base64-decode if specified
    if strings.EqualFold(headers.Get("Content-Transfer-Encoding"), "base64") {
        bodyDecoded := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
        n, decodeErr := base64.StdEncoding.Decode(bodyDecoded, body)
        if decodeErr != nil {
            return nil, fmt.Errorf("base64 decode failed: %v", decodeErr)
        }
        body = bodyDecoded[:n]
    }

    // Step 2: decompress if gzip
    if strings.EqualFold(headers.Get("Content-Encoding"), "gzip") {
        gz, err := gzip.NewReader(bytes.NewReader(body))
        if err != nil {
            return nil, fmt.Errorf("gzip decompress failed: %v", err)
        }
        defer gz.Close()
        body, err = io.ReadAll(gz)
        if err != nil {
            return nil, fmt.Errorf("gzip read failed: %v", err)
        }
    }

    return body, err
}
