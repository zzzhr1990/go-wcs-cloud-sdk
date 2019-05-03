package wcserror

// github.com/zzzhr1990/go-wcs-cloud-sdk/wcserror
import (
	"errors"
)

var (
	// ErrFileExists when file exist
	ErrFileExists = errors.New("file already exists")
)
