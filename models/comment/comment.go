package comment

type Comment struct {
	Id         string `json:"id,omitempty" bson:"_id,omitempty"`
	IdPembuat  string `json:"idpembuat,omitempty" bson:"idpembuat,omitempty"`
	IsiComment string `json:"isicomment,omitempty" bson:"isicomment,omitempty"`
	TglComment string `json:"tglcomment,omitempty" bson:"tglcomment,omitempty"`
}
