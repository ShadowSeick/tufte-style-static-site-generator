package bunny

// In this package I have to get to use:
// - Access key for zone storage
// - Content type must be defined as the type we want to update. We need to set types of the content type. Maybe I should make a pkg for talking through the net or whatever. JUST make it work!
//

// Implementing https://docs.bunny.net/api-reference/storage
func GetFile(path string) {}

func UploadFile(path string, contentType string, data []byte) {}
