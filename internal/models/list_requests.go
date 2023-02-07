package models

type ListFetchRequest struct {
	CacheName string
	ListName  string
}

type ListLengthRequest struct {
	CacheName string
	ListName  string
}
