package est

import (
    "net/http"
    "strings"
)

func ExtractFilename(header http.Header) string {
    cd := header.Get("Content-Disposition")
    if cd == "" {
        return ""
    }
    parts := strings.Split(cd, "filename=")
    if len(parts) < 2 {
        return ""
    }
    filename := strings.Trim(parts[1], "\"")
    return filename
}
