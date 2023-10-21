// 变量的命名一律使用小驼峰命名法，例如：firstName、lastName等。
// 后缀定义:请求req,响应reply


service {{.upperTableName}} {
  //{{.tableNameComment}}-创建一条数据
  rpc {{.upperTableName}}Store({{.upperTableName}}StoreReq) returns ({{.upperTableName}}StoreReply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}_store"
      body: "*"
    };
  }
  //{{.tableNameComment}}-删除多条数据
  rpc {{.upperTableName}}Del({{.upperTableName}}DelReq) returns ({{.upperTableName}}DelReply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}_del"
      body: "*"
    };
  }
  //{{.tableNameComment}}-单条数据查询
  rpc {{.upperTableName}}One({{.upperTableName}}OneReq) returns ({{.upperTableName}}OneReply) {
    option (google.api.http) = {get: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}_info"};
  }
  //{{.tableNameComment}}-列表数据查询
  rpc {{.upperTableName}}List({{.upperTableName}}ListReq) returns ({{.upperTableName}}ListReply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}_list",
      body: "*"
    };
  }
}