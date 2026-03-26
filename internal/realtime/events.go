package realtime

const (
	EventJoinRoom        = "JOIN_ROOM"
	EventLeaveRoom       = "LEAVE_ROOM"
	EventSendMessage     = "SEND_MESSAGE"
	EventTypingStart     = "TYPING_START"
	EventTypingStop      = "TYPING_STOP"
	EventNewMessage      = "NEW_MESSAGE"
	EventUserJoined      = "USER_JOINED"
	EventUserLeft        = "USER_LEFT"
	EventTyping          = "TYPING"
	EventError           = "ERROR"
	EventAuthenticate    = "AUTHENTICATE"
	EventPresence        = "PRESENCE"
	EventPresenceUpdate  = "PRESENCE_UPDATE"
	EventPresenceGet     = "PRESENCE_GET"
	EventPresenceList    = "PRESENCE_LIST"
	EventVoiceJoin       = "VOICE_JOIN"
	EventVoiceLeave      = "VOICE_LEAVE"
	EventVoiceOffer      = "VOICE_OFFER"
	EventVoiceAnswer     = "VOICE_ANSWER"
	EventIceCandidate    = "ICE_CANDIDATE"
	EventVoiceUserJoined = "VOICE_USER_JOINED"
	EventVoiceUserLeft   = "VOICE_USER_LEFT"
)

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type JoinRoomPayload struct {
	ChannelID uint `json:"channel_id"`
}

type LeaveRoomPayload struct {
	ChannelID uint `json:"channel_id"`
}

type SendMessagePayload struct {
	ChannelID uint   `json:"channel_id"`
	Content   string `json:"content"`
}

type TypingPayload struct {
	ChannelID uint   `json:"channel_id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
}

type MessagePayload struct {
	ID        uint   `json:"id"`
	ChannelID uint   `json:"channel_id"`
	AuthorID  uint   `json:"author_id"`
	Content   string `json:"content"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type UserPayload struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type AuthenticatePayload struct {
	Token string `json:"token"`
}

type PresencePayload struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

type PresenceUpdatePayload struct {
	Status string `json:"status"`
}

type PresenceListPayload struct {
	Users []PresencePayload `json:"users"`
}

type PresenceUserPayload struct {
	UserID uint   `json:"user_id"`
	Status string `json:"status"`
}

type VoiceJoinPayload struct {
	ChannelID uint `json:"channel_id"`
}

type VoiceLeavePayload struct {
	ChannelID uint `json:"channel_id"`
}

type VoiceOfferPayload struct {
	ChannelID    uint   `json:"channel_id"`
	TargetUserID uint   `json:"target_user_id"`
	SDP          string `json:"sdp"`
}

type VoiceAnswerPayload struct {
	ChannelID    uint   `json:"channel_id"`
	TargetUserID uint   `json:"target_user_id"`
	SDP          string `json:"sdp"`
}

type IceCandidatePayload struct {
	ChannelID    uint   `json:"channel_id"`
	TargetUserID uint   `json:"target_user_id"`
	Candidate    string `json:"candidate"`
}

type VoiceUserPayload struct {
	ChannelID uint   `json:"channel_id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	IsMuted   bool   `json:"is_muted"`
	IsDeaf    bool   `json:"is_deaf"`
}
