package cache

// TODO 目前各个服务之间的redis任然是强耦合
// 我认为需要分开 耦合不好 数据库也是
// 也应该分的清楚一些
func InitRedis() {
	InitUserRedis()
	InitVideoRedis()
	InitRelationRedis()
	InitFavoriteRedis()
	InitCommentRedis()
}
