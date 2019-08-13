# inker
A simple blog backend.

## 使用

1. 重命名 `config.sample.json` 为 `config.json`，并填入配置信息。
2. `go run main.go setup` 新建用户。
3. `go run main.go` 或编译后运行

注意：输入 `key` 时请脸滚键盘，保证随机性。

## API

### GET /login

登录，返回 Token。

|   参数   | 方式  |  类型  | 说明                                                   |
| :------: | :---: | :----: | :----------------------------------------------------- |
| username | query | string | 用户名                                                 |
| password | query | string | 密码                                                   |
| remember | query |  bool  | 可空，默认为 false。是否记住此 Token（更长的有效时间） |

### GET /getArticle

获取一篇文章。

| 参数 | 方式  |  类型  | 说明   |
| :--: | :---: | :----: | :----- |
| name | query | string | 文章名 |

### GET /paginateHome

获取首页文章。

| 参数  | 方式  | 类型 | 说明                   |
| :---: | :---: | :--: | :--------------------- |
| skip  | query | int  | 可空。跳过的文章数     |
| limit | query | int  | 可空。返回的最大文章数 |

### GET /paginateSearch

搜索文章，会在文章的名称、标题、内容中搜索。

| 参数  | 方式  | 类型 | 说明                   |
| :---: | :---: | :--: | :--------------------- |
| skip  | query | int  | 可空。跳过的文章数     |
| limit | query | int  | 可空。返回的最大文章数 |

### GET /getFile

获取一个文件。

| 参数 | 方式  |  类型  | 说明   |
| :--: | :---: | :----: | ------ |
| name | query | string | 文件名 |

### GET /manage/updateUser

更新用户信息，`newUsername` 与 `newPassword` 不可都为空。

|    参数     | 方式  |  类型  | 说明                       |
| :---------: | :---: | :----: | :------------------------- |
|    token    | query | string | 通过 `/login` 生成的 Token |
| newUsername | query | string | 可空。新用户名             |
| newPassword | query | string | 可空。新密码               |

### POST /manage/newArticle

新建文章，`attr` 为自定义属性，例如可存入文章的封面 URL 和文章分类。

|  参数   | 方式  |  类型  | 说明                                |
| :-----: | :---: | :----: | :---------------------------------- |
|  token  | query | string | 通过 `/login` 生成的 Token          |
|  name   | form  | string | 新建文章名                          |
|  title  | form  | string | 新建文章标题                        |
| content | form  | string | 新建文章内容                        |
|  attr   | form  |  json  | 附加文章属性，如 `{"cover": "URL"}` |

### POST /manage/updateArticle

更新文章。

|    参数    | 方式  |  类型  | 说明                                      |
| :--------: | :---: | :----: | :---------------------------------------- |
|   token    | query | string | 通过 `/login` 生成的 Token                |
|    name    | form  | string | 需要更新的文章名                          |
| updateData | form  |  json  | 更新数据，如 `{"content": "NEW CONTENT"}` |

### GET /manage/deleteArticle

删除文章。

| 参数  | 方式  |  类型  | 说明                       |
| :---: | :---: | :----: | :------------------------- |
| token | query | string | 通过 `/login` 生成的 Token |
| name  | form  | string | 需要删除的文章名           |

### POST /manage/uploadFile

上传文件。

|  参数  | 方式  |  类型  | 说明                       |
| :----: | :---: | :----: | :------------------------- |
| token  | query | string | 通过 `/login` 生成的 Token |
| upload | form  |  file  | 需要更新的文章名           |
|  name  | form  | string | 文件名                     |

### GET /manage/deleteFile

删除文件。

| 参数  | 方式  |  类型  | 说明                       |
| :---: | :---: | :----: | :------------------------- |
| token | query | string | 通过 `/login` 生成的 Token |
| name  | form  | string | 需要删除的文件名           |

### GET /manage/paginateFile

获取文件信息列表。

| 参数  | 方式  | 类型 | 说明                   |
| :---: | :---: | :--: | :--------------------- |
| skip  | query | int  | 可空。跳过的文件数     |
| limit | query | int  | 可空。返回的最大文件数 |

## 关于 DEBUG 模式

可前往 `api/router.go` 修改 `debugMode` 全局变量的值。

注意：开启后会绕过认证中间件的验证（即不需要提供 `token`），以及会在响应中返回堆栈信息与具体异常。

