package responses

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ItemType int32

const (
	ItemTypeScalar     ItemType = ItemType(pb.XItemGetTypeResponse_SCALAR)
	ItemTypeDictionary ItemType = ItemType(pb.XItemGetTypeResponse_DICTIONARY)
	ItemTypeSet        ItemType = ItemType(pb.XItemGetTypeResponse_SET)
	ItemTypeList       ItemType = ItemType(pb.XItemGetTypeResponse_LIST)
	ItemTypeSortedSet  ItemType = ItemType(pb.XItemGetTypeResponse_SORTED_SET)
)

// ItemGetTypeResponse is the base response type for an item get type request.
type ItemGetTypeResponse interface {
	isItemGetTypeResponse()
}

// ItemGetTypeHit hit response to an item get type api request
type ItemGetTypeHit struct {
	value ItemType
}

func (r *ItemGetTypeHit) isItemGetTypeResponse() {}

// TypeString returns the string representation of the item type
func (r *ItemGetTypeHit) TypeString() string {
	return pb.XItemGetTypeResponse_ItemType_name[int32(r.value)]
}

// Type returns the ItemType representation of the item type
func (r *ItemGetTypeHit) Type() ItemType {
	return r.value
}

// NewItemGetTypeHit returns a new ItemGetTypeHit containing the item type
func NewItemGetTypeHit(value pb.XItemGetTypeResponse_ItemType) *ItemGetTypeHit {
	return &ItemGetTypeHit{value: ItemType(value)}
}

// ItemGetTypeMiss miss response to an item get type api request
type ItemGetTypeMiss struct{}

func (r *ItemGetTypeMiss) isItemGetTypeResponse() {}
