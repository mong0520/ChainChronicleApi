package models

type GeneralResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error_code"`
}

type Banner struct {
	Start     int         `json:"start"`
	End       int         `json:"end"`
	Priority  int         `json:"priority"`
	Name      string      `json:"name"`
	URL       string      `json:"url"`
	Type      string      `json:"type"`
	PlaceID   int         `json:"place_id,omitempty"`
	AreaID    int         `json:"area_id,omitempty"`
	TabID     int         `json:"tab_id,omitempty"`
	InfoURL   string      `json:"info_url"`
	GachaType string      `json:"gacha_type,omitempty"`
	GachaID   int         `json:"gacha_id,omitempty"`
	GachaInfo interface{} `json:"gacha_info,omitempty"`
	BannerURL string      `json:"banner_url,omitempty"`
}
type Events struct {
	Deceive     int `json:"deceive"`
	EventPortal []struct {
		ID     int `json:"_id"`
		Events []struct {
			Banner []Banner `json:"banner"`
			Icon   string   `json:"icon"`
			Img    []struct {
				Start int    `json:"start"`
				End   int    `json:"end"`
				Chara string `json:"chara"`
			} `json:"img"`
			URL      string `json:"url"`
			End      int    `json:"end"`
			Start    int    `json:"start"`
			Priority int    `json:"priority"`
			Category string `json:"category"`
		} `json:"events"`
		Timestamp string `json:"timestamp"`
	} `json:"event_portal"`
	Res int `json:"res"`
}
