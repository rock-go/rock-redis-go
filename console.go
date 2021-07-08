package redis

import "github.com/rock-go/rock/lua"

func (r *Redis) Header(out lua.Printer) {
	out.Printf("type: redis client")
	out.Printf("uptime: %s", r.uptime)
	out.Printf("version: v1.0.0")
	out.Println("")
}

func (r *Redis) Show(out lua.Printer) {
	r.Header(out)
	out.Printf("name: %s", r.C.name)
	out.Printf("addr: %s", r.C.addr)
	out.Printf("password: ********")
	out.Printf("db: %d", r.C.db)
	out.Printf("pool_size: %d", r.C.poolSize)
	out.Printf("max_age_size: %d", r.C.maxConnAge)
	out.Println("")
}

func (r *Redis) Help(out lua.Printer) {
	r.Header(out)
	out.Println(".start() 启动")
	out.Println(".close() 关闭")
}
