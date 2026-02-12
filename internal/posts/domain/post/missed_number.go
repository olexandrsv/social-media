package post

type MissedNumber struct {
	followingID int
	number      int
}

func NewMissedNumber(followingID int, number int) *MissedNumber {
	return &MissedNumber{
		followingID: followingID,
		number:      number,
	}
}

func (m *MissedNumber) FollowingID() int {
	return m.followingID
}

func (m *MissedNumber) Number() int {
	return m.number
}
