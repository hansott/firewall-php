package utils

import (
	"path/filepath"
	"strings"
)

const (
	NOT_FOUND          = 404
	METHOD_NOT_ALLOWED = 405
	MOVED_PERMANENTLY  = 301
	FOUND              = 302
	SEE_OTHER          = 303
	TEMPORARY_REDIRECT = 307
	PERMANENT_REDIRECT = 308
)

var ERROR_CODES = []int{NOT_FOUND, METHOD_NOT_ALLOWED}
var REDIRECT_CODES = []int{
	MOVED_PERMANENTLY,
	FOUND,
	SEE_OTHER,
	TEMPORARY_REDIRECT,
	PERMANENT_REDIRECT,
}
var EXCLUDED_METHODS = []string{"OPTIONS", "HEAD"}
var IGNORE_EXTENSIONS = []string{"properties", "asp", "aspx", "jsp", "config"}
var ALLOW_EXTENSIONS = []string{"html", "php"}
var IGNORE_STRINGS = []string{"cgi-bin"}

func ShouldDiscoverRoute(statusCode int, route, method string) bool {
	if containsStr(EXCLUDED_METHODS, method) {
		return false
	}

	if containsInt(ERROR_CODES, statusCode) {
		return false
	}

	if containsInt(REDIRECT_CODES, statusCode) {
		return false
	}

	segments := strings.Split(route, "/")

	// e.g. /path/to/.file or /.directory/file
	for _, segment := range segments {
		if isDotFile(segment) {
			return false
		}
	}

	for _, segment := range segments {
		if containsIgnoredString(segment) {
			return false
		}
	}

	for _, segment := range segments {
		if !isAllowedExtension(segment) {
			return false
		}
	}

	return true
}

func isAllowedExtension(segment string) bool {
	extension := filepath.Ext(segment)

	if extension != "" && strings.HasPrefix(extension, ".") {
		extension = extension[1:]

		if containsStr(ALLOW_EXTENSIONS, extension) {
			return true
		}

		if len(extension) >= 2 && len(extension) <= 5 {
			return false
		}

		if containsStr(IGNORE_EXTENSIONS, extension) {
			return false
		}
	}

	return true
}

func isDotFile(segment string) bool {
	if segment == ".well-known" {
		return false
	}

	return strings.HasPrefix(segment, ".") && len(segment) > 1
}

func containsIgnoredString(segment string) bool {
	for _, str := range IGNORE_STRINGS {
		if strings.Contains(segment, str) {
			return true
		}
	}
	return false
}

func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
