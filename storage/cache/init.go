package cache

// TODO 目前各个服务之间的redis任然是强耦合
// 我认为需要分开 耦合不好
func InitRedis() {
	InitUserRedis()
	InitVideoRedis()
	InitRelationRedis()
	InitFavoriteRedis()
	InitCommentRedis()
}
