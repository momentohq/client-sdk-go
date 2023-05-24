package responses

import (
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ItemType int32

const (
	Scalar     ItemType = ItemType(pb.XItemGetTypeResponse_SCALAR)
	Dictionary ItemType = ItemType(pb.XItemGetTypeResponse_DICTIONARY)
	Set        ItemType = ItemType(pb.XItemGetTypeResponse_SET)
	List       ItemType = ItemType(pb.XItemGetTypeResponse_LIST)
	SortedSet  ItemType = ItemType(pb.XItemGetTypeResponse_SORTED_SET)
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
