package main

import (
	"diktok/config"
	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/favorite/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://favorite", constant.ServiceName+".favorite")
	// defer shutdown()
	database.InitMySQL()
	storage.FavoriteRedis = cache.InitRedis(config.System.Redis.FavoriteDB)
	storage.VideoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)
	rpc.InitServerWithNacos(constant.FavoriteAddr, constant.FavoriteService, pbfavorite.RegisterFavoriteServer, &FavoriteService{})
}
