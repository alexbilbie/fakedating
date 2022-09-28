package payload

import (
	"github.com/segmentio/ksuid"
)

type SwipeRequest struct {
	Recipient ksuid.KSUID
	Matched   bool
}

type SwipeResponse struct {
	MutualMatch bool
}
