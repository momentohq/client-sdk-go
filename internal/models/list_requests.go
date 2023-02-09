package models

import (
	"github.com/momentohq/client-sdk-go/utils"
)

type ListFetchRequest struct {
	CacheName string
	ListName  string
}

type ListLengthRequest struct {
	CacheName string
	ListName  string
}

type ListPushFrontRequest struct {
	CacheName          string
	ListName           string
	Value              []byte
	TruncateBackToSize uint32
	CollectionTtl      utils.CollectionTtl
}

type ListPushBackRequest struct {
	CacheName           string
	ListName            string
	Value               []byte
	TruncateFrontToSize uint32
	CollectionTtl       utils.CollectionTtl
}
