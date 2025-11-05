package models

// TiktokResponses represents the root response structure from TikTok API
type TiktokResponses struct {
	LiveRoom    TiktokLiveRoom    `json:"LiveRoom"`
	CurrentRoom TiktokCurrentRoom `json:"CurrentRoom"`
}

// TiktokLiveRoom represents the LiveRoom section
type TiktokLiveRoom struct {
	LoadingState       TiktokLoadingState     `json:"loadingState"`
	NeedLogin          bool                   `json:"needLogin"`
	ShowLiveGate       bool                   `json:"showLiveGate"`
	IsAgeGateRoom      bool                   `json:"isAgeGateRoom"`
	RecommendLiveRooms []any                  `json:"recommendLiveRooms"`
	LiveRoomStatus     int                    `json:"liveRoomStatus"`
	LiveRoomUserInfo   TiktokLiveRoomUserInfo `json:"liveRoomUserInfo"`
}

// TiktokLoadingState represents loading state fields
type TiktokLoadingState struct {
	GetRecommendLive int `json:"getRecommendLive"`
	GetUserInfo      int `json:"getUserInfo"`
	GetUserStat      int `json:"getUserStat"`
}

// TiktokLiveRoomUserInfo represents user info in live room
type TiktokLiveRoomUserInfo struct {
	User     TiktokUser     `json:"user"`
	Stats    TiktokStats    `json:"stats"`
	LiveRoom TiktokRoomInfo `json:"liveRoom"`
}

// TiktokUser represents user information
type TiktokUser struct {
	AvatarLarger string `json:"avatarLarger"`
	AvatarMedium string `json:"avatarMedium"`
	AvatarThumb  string `json:"avatarThumb"`
	ID           string `json:"id"`
	Nickname     string `json:"nickname"`
	SecUid       string `json:"secUid"`
	Secret       bool   `json:"secret"`
	UniqueId     string `json:"uniqueId"`
	Verified     bool   `json:"verified"`
	RoomId       string `json:"roomId"`
	Signature    string `json:"signature"`
	Status       int    `json:"status"`
	FollowStatus int    `json:"followStatus"`
}

// TiktokStats represents user statistics
type TiktokStats struct {
	FollowingCount int `json:"followingCount"`
	FollowerCount  int `json:"followerCount"`
}

// TiktokRoomInfo represents live room information
type TiktokRoomInfo struct {
	CoverUrl          string              `json:"coverUrl"`
	SquareCoverImg    string              `json:"squareCoverImg"`
	Title             string              `json:"title"`
	StartTime         int64               `json:"startTime"`
	Status            int                 `json:"status"`
	PaidEvent         TiktokPaidEvent     `json:"paidEvent"`
	LiveSubOnly       int                 `json:"liveSubOnly"`
	LiveRoomMode      int                 `json:"liveRoomMode"`
	GameTagId         int                 `json:"gameTagId"`
	LiveRoomStats     TiktokLiveRoomStats `json:"liveRoomStats"`
	StreamData        TiktokStreamData    `json:"streamData"`
	StreamId          string              `json:"streamId"`
	MultiStreamScene  int                 `json:"multiStreamScene"`
	MultiStreamSource int                 `json:"multiStreamSource"`
	HevcStreamData    TiktokStreamData    `json:"hevcStreamData"`
}

// TiktokPaidEvent represents paid event information
type TiktokPaidEvent struct {
	EventId  int `json:"event_id"`
	PaidType int `json:"paid_type"`
}

// TiktokLiveRoomStats represents live room statistics
type TiktokLiveRoomStats struct {
	EnterCount int `json:"enterCount"`
	UserCount  int `json:"userCount"`
}

// TiktokStreamData represents stream data structure
type TiktokStreamData struct {
	PullData TiktokPullData `json:"pull_data"`
	PushData TiktokPushData `json:"push_data"`
}

// TiktokPullData represents pull data structure
type TiktokPullData struct {
	Options    TiktokStreamOptions `json:"options"`
	StreamData string              `json:"stream_data"`
}

// TiktokStreamOptions represents stream quality options
type TiktokStreamOptions struct {
	DefaultPreviewQuality TiktokQualityOption   `json:"default_preview_quality"`
	DefaultQuality        TiktokQualityOption   `json:"default_quality"`
	Qualities             []TiktokQualityOption `json:"qualities"`
	ShowQualityButton     bool                  `json:"show_quality_button"`
}

// TiktokQualityOption represents a quality option
type TiktokQualityOption struct {
	IconType   int    `json:"icon_type"`
	Level      int    `json:"level"`
	Name       string `json:"name"`
	Resolution string `json:"resolution"`
	SdkKey     string `json:"sdk_key"`
	VCodec     string `json:"v_codec"`
}

// TiktokPushData represents push data structure
type TiktokPushData struct {
	PushStreamLevel  int                    `json:"push_stream_level"`
	ResolutionParams map[string]interface{} `json:"resolution_params"`
	StreamData       string                 `json:"stream_data"`
}

// TiktokCurrentRoom represents the CurrentRoom section
type TiktokCurrentRoom struct {
	LoadingState      TiktokCurrentRoomLoadingState `json:"loadingState"`
	RoomInfo          interface{}                   `json:"roomInfo"`
	AnchorId          string                        `json:"anchorId"`
	SecAnchorId       string                        `json:"secAnchorId"`
	AnchorUniqueId    string                        `json:"anchorUniqueId"`
	RoomId            string                        `json:"roomId"`
	HotLiveRoomInfo   interface{}                   `json:"hotLiveRoomInfo"`
	LiveType          string                        `json:"liveType"`
	ReportLinkType    string                        `json:"reportLinkType"`
	EnterRoomWithSSR  bool                          `json:"enterRoomWithSSR"`
	PlayMode          string                        `json:"playMode"`
	IsGuestConnection bool                          `json:"isGuestConnection"`
	IsMultiGuestRoom  bool                          `json:"isMultiGuestRoom"`
	ShowLiveChat      bool                          `json:"showLiveChat"`
	EnableChat        bool                          `json:"enableChat"`
	IsAnswerRoom      bool                          `json:"isAnswerRoom"`
	IsGateRoom        bool                          `json:"isGateRoom"`
	RequestId         string                        `json:"requestId"`
	NtpDiff           int                           `json:"ntpDiff"`
	FollowStatusMap   map[string]interface{}        `json:"followStatusMap"`
}

// TiktokCurrentRoomLoadingState represents loading state for current room
type TiktokCurrentRoomLoadingState struct {
	EnterRoom int `json:"enterRoom"`
}

type TiktokVideoQualityInfo struct {
	URL        string
	VBitrate   int
	Resolution [2]int
}
