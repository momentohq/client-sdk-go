package momento

type ListCachesRequest struct {
	// Token to continue paginating through the list. It's used to handle large paginated lists.
	NextToken string
}
