package internal_serve

func isWriteableMime(mimeType string) bool {
	if mimeType == "text/html" ||
		mimeType == "text/javascript" ||
		mimeType == "text/json" ||
		mimeType == "application/javascript" ||
		mimeType == "application/json" ||
		mimeType == "text/css" {
		return true
	}
	return false
}
