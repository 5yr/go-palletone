package dagconfig

var (
	SConfig Sconfig
)

var DefaultConfig = Config{
	DbPath: "db/leveldb",
	DbName: "palletone",
}

// key := strings.ToLower(typ.Name()) 大写统一转小写
type Config struct {
	DbPath     string
	DbName     string
	DbUser     string
	DbPassword string
	DbCache    int
	DbHandles  int

	// cache
	CacheSource string

	//redis
	RedisAddr   string
	RedisPwd    string
	RedisPrefix string
	RedisDb     int
}

type Sconfig struct {
	Blight bool
}
