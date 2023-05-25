package responses

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ItemType int32

const (
	ItemTypeScalar     ItemType = 0
	ItemTypeDictionary ItemType = 1
	ItemTypeSet        ItemType = 2
	ItemTypeList       ItemType = 3
	ItemTypeSortedSet  ItemType = 4
)

// ItemGetTypeResponse is the base response type for an item get type request.
type ItemGetTypeResponse interface {
	isItemGetTypeResponse()
}

// ItemGetTypeHit hit response to an item get type api request.
type ItemGetTypeHit struct {
	value ItemType
}

func (r *ItemGetTypeHit) isItemGetTypeResponse() {}

// Type returns the ItemType representation of the item type.
func (r *ItemGetTypeHit) Type() ItemType {
	return r.value
}

// NewItemGetTypeHit returns a new ItemGetTypeHit containing the item type.
func NewItemGetTypeHit(value pb.XItemGetTypeResponse_ItemType) *ItemGetTypeHit {
	return &ItemGetTypeHit{value: ItemType(value)}
}

// ItemGetTypeMiss miss response to an item get type api request.
type ItemGetTypeMiss struct{}

func (r *ItemGetTypeMiss) isItemGetTypeResponse() {}
