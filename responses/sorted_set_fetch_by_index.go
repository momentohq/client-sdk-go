// These exist in case in the future we want SortedSetFetchByIndex and
// SortedSetFetchByScore return different responses.
package responses

type SortedSetFetchByIndexResponse SortedSetFetchResponse

type SortedSetFetchByIndexMiss SortedSetFetchMiss

type SortedSetFetchByIndexHit SortedSetFetchHit
