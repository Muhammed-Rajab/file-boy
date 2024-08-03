package codec

import "errors"

var ErrNotEncryptFile = errors.New("file not encrypted")
var ErrPathDoesNotExist = errors.New("path doesn't exist")
