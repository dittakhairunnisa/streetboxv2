package model

// ReqUpdateAppSettingByKey ..
type ReqUpdateAppSettingByKey struct {
	Value  string `json:"value"`
}

// AppSetting ..
type AppSetting struct {
	ID   	int64  `json:"id"`
	Key     string `json:"key"`
	Value 	string `json:"value"`
}
