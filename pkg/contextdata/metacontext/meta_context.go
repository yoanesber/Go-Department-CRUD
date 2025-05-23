package metacontext

import (
	"context"
	"fmt"
)

// This struct defines the RequestMeta struct
//
//	It can be used to store metadata about the request
type RequestMeta struct {
	UserID   int64
	UserName string
	Email    string
	Roles    []string
}

// This struct defines the requestMetaKeyType struct
//
//	It is used as a key for storing and retrieving RequestMeta from the context
type requestMetaKeyType struct{}

// Define a key for storing RequestMeta in the context
var requestMetaKey = requestMetaKeyType{}

// GetValueFromContext retrieves a value from the context using the provided key.
// It returns the value and an error if the key does not exist in the context.
func GetValueFromContext(ctx context.Context, key string) (interface{}, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, fmt.Errorf("key %s not found in context", key)
	}
	return value, nil
}

// InjectRequestMeta injects the RequestMeta into the context.
// This function is used to add metadata to the context for later retrieval
func InjectRequestMeta(ctx context.Context, meta RequestMeta) context.Context {
	return context.WithValue(ctx, requestMetaKey, meta)
}

// ExtractRequestMeta retrieves the RequestMeta from the context.
// This function is used to access the metadata stored in the context
func ExtractRequestMeta(ctx context.Context) (RequestMeta, bool) {
	meta, ok := ctx.Value(requestMetaKey).(RequestMeta)
	return meta, ok
}
