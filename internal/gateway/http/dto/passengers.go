package dto

type PassengerBriefHTTP struct {
	PassengerID uint64 `json:"passenger_id"`
	RealName    string `json:"real_name"`
	IDCard      string `json:"id_card"`
}

type ListPassengersHTTPResp struct {
	Code       int32                 `json:"code"`
	Msg        string                `json:"msg"`
	Passengers []PassengerBriefHTTP  `json:"passengers,omitempty"`
}

