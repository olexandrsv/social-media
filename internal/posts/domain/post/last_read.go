package post

type LastRead struct {
	followingID    int
	lastReadPostID string
}

func NewLastReadPost(followingID int, lastReadPostID string) *LastRead {
	return &LastRead{
		followingID:    followingID,
		lastReadPostID: lastReadPostID,
	}
}

func (l *LastRead) FollowingID() int {
	return l.followingID
}

func (l *LastRead) LastReadPostID() string {
	return l.lastReadPostID
}
