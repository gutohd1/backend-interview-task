package database

type UserModel struct {
	Id        string
	Name      string
	Likes     uint
	Gender    string
	CreatedAt uint64
	UpdatedAt uint64
	IsAactive bool
}

type DecisionModel struct {
	Id          string
	ActorId     uint
	RecipientId uint
	Liked       bool
	Created_at  uint64
	Updated_at  uint64
}

type PutDecisionEntry struct {
	ActorId     string
	RecipientId string
	Like        bool
}

type DecisionsTotalLikes struct {
	Likes int
}
