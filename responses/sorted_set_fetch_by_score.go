// These exist in case in the future we want SortedSetFetchByIndex and
// SortedSetFetchByScore return different responses.
package responses

type SortedSetFetchByScoreResponse SortedSetFetchResponse

type SortedSetFetchByScoreMiss SortedSetFetchMiss

type SortedSetFetchByScoreHit SortedSetFetchHit
