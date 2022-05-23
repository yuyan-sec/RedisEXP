package pkg

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
)

var (
	ctx      = context.Background()
	FileName = "out.json"
)

type String struct {
	Key string `json:"key"`
	Val string `json:"val"`
}
type HashData struct {
	Field string `json:"field"`
	Val   string `json:"val"`
}
type Hash struct {
	Key  string     `json:"key"`
	Data []HashData `json:"data"`
}
type ZsetData struct {
	Member string  `json:"member"`
	Score  float64 `json:"score"`
}
type Zset struct {
	Key  string     `json:"key"`
	Data []ZsetData `json:"data"`
}
type Set struct {
	Key  string   `json:"key"`
	Data []string `json:"data"`
}
type List struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type RedisInfo struct {
	String []String `json:"string"`
	Hash   []Hash   `json:"hash"`
	Set    []Set    `json:"set"`
	List   []List   `json:"list"`
	Zset   []Zset   `json:"zset"`
}

// Redis 导出
func handle_export() {
	redisInfo := RedisInfo{}
	var cursor uint64
	keys, cursor, err := Rdb.Scan(ctx, cursor, "*", 100).Result()
	if err != nil {
		Err(err)
		return
	}

	for _, key := range keys {
		sType, err := Rdb.Type(ctx, key).Result()
		if err != nil {
			Err(err)
			return
		}

		switch sType {
		case "string":
			val, _ := Rdb.Get(ctx, key).Result()
			strInfo := String{
				Key: key,
				Val: val,
			}
			redisInfo.String = append(redisInfo.String, strInfo)

		case "list":
			val, _ := Rdb.LPop(ctx, key).Result()
			listInfo := List{
				Key:  key,
				Data: val,
			}

			redisInfo.List = append(redisInfo.List, listInfo)

		case "hash":
			val, _ := Rdb.HGetAll(ctx, key).Result()
			hashInfo := Hash{
				Key: key,
			}

			for k, v := range val {
				data := HashData{
					Field: k,
					Val:   v,
				}
				hashInfo.Data = append(hashInfo.Data, data)
			}
			redisInfo.Hash = append(redisInfo.Hash, hashInfo)

		case "set":
			val, _ := Rdb.SMembers(ctx, key).Result()
			setInfo := Set{
				Key:  key,
				Data: val,
			}
			redisInfo.Set = append(redisInfo.Set, setInfo)

		case "zset":

			val, _ := Rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
			zsetInfo := Zset{
				Key: key,
			}
			zs := []ZsetData{}

			for _, z := range val {
				zs = append(zs, ZsetData{
					Member: z.Member.(string),
					Score:  z.Score,
				})

			}
			zsetInfo.Data = zs
			redisInfo.Zset = append(redisInfo.Zset, zsetInfo)
		}

	}
	bs, err := json.Marshal(redisInfo)
	if err != nil {
		Info("序列化成json失败" + err.Error())
	}
	err = ioutil.WriteFile(FileName, bs, 0644)
	if err != nil {
		Info("保存到文件失败" + err.Error())
	}

	Success(FileName + "  导出成功")
}

// Redis 导入
func handle_import() {

	bs, err := ioutil.ReadFile(FileName)
	if err != nil {
		Info("读取文件失败" + err.Error())
	}
	redis_info := RedisInfo{}
	err = json.Unmarshal(bs, &redis_info)
	if err != nil {
		Info("不是合法的json文件" + err.Error())
	}
	//string
	for _, v := range redis_info.String {
		Rdb.Set(ctx, v.Key, v.Val, 0)
	}
	//hash
	for _, v := range redis_info.Hash {
		maps := map[string]string{}
		for _, d := range v.Data {
			maps[d.Field] = d.Val
		}
		Rdb.HMSet(ctx, v.Key, maps)
	}
	//set
	for _, v := range redis_info.Set {

		for _, v1 := range v.Data {
			Rdb.SAdd(ctx, v.Key, v1)
		}
	}
	//zset

	for _, v := range redis_info.Zset {

		zs := []*redis.Z{}
		for _, v1 := range v.Data {
			zs = append(zs, &redis.Z{
				Member: v1.Member,
				Score:  v1.Score,
			})
		}

		Rdb.ZAdd(ctx, v.Key, zs...)
	}
	//list
	for _, v := range redis_info.List {
		for _, v1 := range v.Data {
			Rdb.RPush(ctx, v.Key, v1)
		}
	}

	Success(FileName + "  导入成功")

}
