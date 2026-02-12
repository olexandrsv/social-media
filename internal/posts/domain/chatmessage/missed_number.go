package chatmessage

type MissedNumber struct {
	chatID int
	number int
}

func NewMissedNumber(chatID, number int) *MissedNumber {
	return &MissedNumber{
		chatID: chatID,
		number: number,
	}
}

func (m *MissedNumber) ChatID() int {
	return m.chatID
}

func (m *MissedNumber) Number() int {
	return m.number
}
