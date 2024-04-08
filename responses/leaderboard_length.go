package responses

type LeaderboardLengthResponse interface {
	isLeaderboardLengthResponse()
}

type LeaderboardLengthSuccess struct {
	length uint32
}

func (LeaderboardLengthSuccess) isLeaderboardLengthResponse() {}

// NewLeaderboardLengthSuccess returns a new LeaderboardLengthSuccess containing the supplied elements.
func NewLeaderboardLengthSuccess(length uint32) *LeaderboardLengthSuccess {
	return &LeaderboardLengthSuccess{length: length}
}

func (s LeaderboardLengthSuccess) Length() uint32 {
	return s.length
}
