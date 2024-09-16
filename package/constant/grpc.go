package constant

// grpc
var (
	UserService     = "diktok/user"
	VideoService    = "diktok/video"
	RalationService = "diktok/relation"
	CommentService  = "diktok/comment"
	MessageService  = "diktok/message"
	FavoriteService = "diktok/favorite"
	VideoAddr       = "video:8010" //docker 内使用video:8010 本地使用127.0.0.1 否则会阻塞
	UserAddr        = "user:8020"
	RelationAddr    = "relation:8030"
	MessageAddr     = "message:8040"
	FavoriteAddr    = "favorite:8050"
	CommentAddr     = "comment:8060"
)
