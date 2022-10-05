# porm
A orm based on prisma and graphql


# start 

go run github.com/prisma/prisma-client-go db push


# TODO

- [x] 程序总入口，配置等
- [x] 枚举值特性支持
- [ ] 原生查询sdk生成(先支持原生，然后再生成)
- mutation {result: queryRaw(query:"SELECT * FROM Post WHERE id= ?",parameters:"[1]") }

- [ ] 当前功能case编写
- [ ] 多opearation支持（文件命名）
- [ ] 事务支持(批量查询)
- [x] 内置graphql端点（代理~）
    - 启动prisma 引擎的指令~

- [ ] prisma 引擎集成？弄一个安装的命令
    - xx inti中判断是否下载了？
    - 生成 graphql schema
- [ ] prisma 配置端口
- [ ] 日志


bug fix

- [x] 重复的变量定义~
- [x] datatime支持
- [ ] 修复不支持数组的bug

# thx

- sqlc
- prisma
- graphql
