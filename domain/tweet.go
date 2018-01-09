package domain

import "time"

type Tweet struct {
	ID   uint64
	From UserID
	At   time.Time
	Text string
}

// TweetWithUsername shadows From field with username of a tweeter.
type TweetWithUsername struct {
	Tweet
	From string
}
