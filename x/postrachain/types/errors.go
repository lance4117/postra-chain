package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/postrachain module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrPostNotFound  = errors.Register(ModuleName, 1101, "post not found")
	ErrInvalidTitle  = errors.Register(ModuleName, 1102, "invalid title")
	ErrInvalidContentURI  = errors.Register(ModuleName, 1103, "invalid content uri")
	ErrInvalidContentHash = errors.Register(ModuleName, 1104, "invalid content hash")
)
