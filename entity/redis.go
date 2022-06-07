package entity

type RedisSetToken struct {
	Token               string
	ExpireTime          int
	DomainIdentity      int
	ServiceIdentity     int
	ApplicationIdentity int
}

type RedisApplicationPermission struct {
	AppId  int
	Method string
	Value  bool
}
