package codegen

import (
	"hash"
	"hash/fnv"
)

const (
	TypeID_VOID    = "void"
	TypeID_STRING  = "string"
	TypeID_INTEGER = "integer"
	TypeID_FLOAT   = "float"
	TypeID_BOOLEAN = "boolean"
	TypeID_USER    = "user"
)

// DynamicType specifies a dynamic type that is specified by the user
type DynamicType struct {
	TypeID    string `json:"typeID"`
	Reference string `json:"reference"`
}

// ITypeMapper provides an interface for mapping dynamic types into language-specific types.
type ITypeMapper interface {
	Convert(dtype DynamicType) (string, error)
}

type UserTypeCache struct {
	cache  map[string]DynamicType
	hasher hash.Hash32
}

func NewUserTypeCache() *UserTypeCache {
	hasher := fnv.New32()
	return &UserTypeCache{
		cache:  make(map[string]DynamicType),
		hasher: hasher,
	}
}

func (userTypeCache *UserTypeCache) getEndpointKey(endpoint APIEndpoint) string {
	bytes := userTypeCache.hasher.Sum([]byte(endpoint.Name))
	bytes = userTypeCache.hasher.Sum([]byte(endpoint.Method))
	userTypeCache.hasher.Reset()

	return string(bytes)
}

func (userTypeCache *UserTypeCache) CacheEndpointInputType(endpoint APIEndpoint, inputType DynamicType) {
	key := userTypeCache.getEndpointKey(endpoint)
	userTypeCache.cache[key] = inputType
}

func (userTypeCache *UserTypeCache) GetEndpointInputType(endpoint APIEndpoint) (DynamicType, bool) {
	key := userTypeCache.getEndpointKey(endpoint)
	rtype, ok := userTypeCache.cache[key]
	return rtype, ok
}
