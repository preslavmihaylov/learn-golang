package bot

type BlackjackStats struct {
	Balance   int
	HandsWon  int
	HandsLost int
	HandsTied int
}

func (bs *BlackjackStats) Accumulate(other *BlackjackStats) {
	bs.Balance += other.Balance
	bs.HandsWon += other.HandsWon
	bs.HandsTied += other.HandsTied
	bs.HandsLost += other.HandsLost
}
