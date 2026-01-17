package burpsuite

// Burp Suite constants
const (
	// JarName is the name of the Burp Suite JAR file.
	JarName = "burpsuite.jar"

	// baseURL is the PortSwigger CDN base URL for downloads.
	baseURL = "https://portswigger-cdn.net/burp/releases/download"
)

// burpDownloadURL constructs the download URL for a Burp Suite release.
// Downloads latest version - version parameter not needed.
func burpDownloadURL(edition string) string {
	product := editionToProduct(edition)
	return baseURL + "?product=" + product + "&type=Jar"
}

// editionToProduct maps config edition to PortSwigger product identifier.
func editionToProduct(edition string) string {
	switch edition {
	case "community":
		return "community"
	case "professional":
		return "pro"
	default:
		return "pro"
	}
}
