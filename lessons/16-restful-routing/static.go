package restfulrouting

import (
	"net/http"
	"strings"
	"time"
)

// FileServerHandler shows the standard-library way to serve a directory of
// static files. Three moves:
//
//  1. http.Dir wraps a filesystem path as an http.FileSystem.
//  2. http.FileServer produces an http.Handler that serves that filesystem.
//  3. http.StripPrefix removes the URL prefix before the FileServer looks
//     up the path — without it, /assets/foo.txt would try to find
//     "./public/assets/foo.txt" instead of "./public/foo.txt".
func FileServerHandler(root, urlPrefix string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	// Ensure prefix has both leading and trailing slash for StripPrefix
	// symmetry with the mux pattern below.
	prefix := "/" + strings.Trim(urlPrefix, "/") + "/"
	return http.StripPrefix(prefix, fs)
}

// ServeReadmeContent shows the underlying primitive http.FileServer uses:
// http.ServeContent. It honors If-Modified-Since, ETag/If-None-Match, and
// Range headers — so a client asking for bytes 100–200 gets a 206 Partial
// Content response automatically. Use ServeContent when you already have
// the bytes in memory or as an io.ReadSeeker but don't want to reinvent
// range and cache handling.
func ServeReadmeContent(w http.ResponseWriter, r *http.Request, name string, modTime time.Time, content string) {
	// The last arg must be an io.ReadSeeker — strings.Reader satisfies that.
	http.ServeContent(w, r, name, modTime, strings.NewReader(content))
}
