package services

import (
	"errors"

	"github.com/livekit/protocol/auth"
)

type SFUService struct {
	apiKey    string
	apiSecret string
	serverURL string
}

func NewSFUService(apiKey, apiSecret, serverURL string) (*SFUService, error) {
	if apiKey == "" || apiSecret == "" || serverURL == "" {
		return nil, errors.New("LiveKit configuration is required")
	}

	return &SFUService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		serverURL: serverURL,
	}, nil
}

func (s *SFUService) GenerateToken(userID, roomName, username string) (string, error) {
	token := auth.NewAccessToken(s.apiKey, s.apiSecret)
	videoGrant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	token.AddGrant(videoGrant)
	token.SetIdentity(userID)
	token.SetName(username)

	return token.ToJWT()
}

func (s *SFUService) CreateRoom(roomName string, maxParticipants int) error {
	return errors.New("room creation not implemented - use LiveKit dashboard or API directly")
}

func (s *SFUService) DeleteRoom(roomName string) error {
	return errors.New("room deletion not implemented - use LiveKit dashboard or API directly")
}

func (s *SFUService) GetParticipants(roomName string) ([]Participant, error) {
	return []Participant{}, nil
}

func (s *SFUService) IsEnabled() bool {
	return s.apiKey != "" && s.apiSecret != "" && s.serverURL != ""
}

type Participant struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMuted  bool   `json:"is_muted"`
	IsScreen bool   `json:"is_screen"`
}

type JoinVoiceSFUInput struct {
	ChannelID uint
	UserID    uint
	Username  string
}

type JoinVoiceSFUResult struct {
	RoomName     string        `json:"room_name"`
	Token        string        `json:"token"`
	ServerURL    string        `json:"server_url"`
	Participants []Participant `json:"participants"`
}

func (s *SFUService) JoinVoice(input JoinVoiceSFUInput) (*JoinVoiceSFUResult, error) {
	roomName := s.getRoomName(input.ChannelID)

	token, err := s.GenerateToken(string(rune(input.UserID)), roomName, input.Username)
	if err != nil {
		return nil, err
	}

	return &JoinVoiceSFUResult{
		RoomName:     roomName,
		Token:        token,
		ServerURL:    s.serverURL,
		Participants: []Participant{},
	}, nil
}

func (s *SFUService) LeaveVoice(channelID, userID uint) error {
	return nil
}

func (s *SFUService) getRoomName(channelID uint) string {
	return "voice-channel-" + string(rune(channelID))
}

func (s *SFUService) GetSFUInfo() (string, string) {
	return s.serverURL, s.apiKey
}

func (s *SFUService) ListRooms() ([]string, error) {
	return []string{}, nil
}
