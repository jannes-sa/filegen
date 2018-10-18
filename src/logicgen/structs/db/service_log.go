package db

// ServiceLog - ServiceLogStruct
type ServiceLog struct {
	Type    string `bson:"type"`
	Req     string `bson:"req"`
	Res     string `bson:"res"`
	Errcode string `bson:"errcode"`
	JobID   string `bson:"job_id"`
}
