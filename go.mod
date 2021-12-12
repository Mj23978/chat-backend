module github.com/mj23978/chat-backend

go 1.16

replace (
	go.etcd.io/etcd/api/v3 v3.5.0-pre => go.etcd.io/etcd/api/v3 v3.0.0-20210107172604-c632042bb96c
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre => go.etcd.io/etcd/pkg/v3 v3.0.0-20210107172604-c632042bb96c
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gocql/gocql v0.0.0-20211015133455-b225f9b53fa1
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/copier v0.3.2
	github.com/mj23978/chat-backend-x v0.0.0-20211201052328-b035391b4b1c
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.20.0
	github.com/scylladb/gocqlx/v2 v2.6.0
	github.com/spf13/viper v1.9.0
	github.com/thoas/go-funk v0.9.1
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/dealancer/validate.v2 v2.1.0
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.4
)
