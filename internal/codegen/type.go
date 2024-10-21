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
	TypeID_ARRAY   = "array"
	TypeID_GENERIC = "generic"
)

// DynamicType specifies a dynamic type that is specified by the user
type DynamicType struct {
	TypeID    string        `json:"typeID"`
	Reference string        `json:"reference"`
	Inner     []DynamicType `json:"inner"`
}

func (tp DynamicType) IsVoid() bool {
	return tp.TypeID == TypeID_VOID
}

func (tp DynamicType) ArrayElementTp() DynamicType {
	if tp.TypeID != TypeID_ARRAY {
		panic("type is not an array")
	}

	if len(tp.Inner) < 1 {
		panic("array type does not have at least one inner type")
	}

	return tp.Inner[0]
}

func (tp DynamicType) TypeReferences() []string {
	references := []string{tp.Reference}
	for _, inner := range tp.Inner {
		inner := inner.TypeReferences()
		references = append(references, inner...)
	}

	return references
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
