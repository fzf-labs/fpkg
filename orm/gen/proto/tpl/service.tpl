service {{.upperTableName}} {
  //{{.tableNameComment}}-创建一条数据
  rpc Create{{.upperTableName}}(Create{{.upperTableName}}Req) returns (Create{{.upperTableName}}Reply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}/create"
      body: "*"
    };
  }
  //{{.tableNameComment}}-更新一条数据
  rpc Update{{.upperTableName}}(Update{{.upperTableName}}Req) returns (Update{{.upperTableName}}Reply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}/update"
      body: "*"
    };
  }
  //{{.tableNameComment}}-删除多条数据
  rpc Delete{{.upperTableName}}(Delete{{.upperTableName}}Req) returns (Delete{{.upperTableName}}Reply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}/delete"
      body: "*"
    };
  }
  //{{.tableNameComment}}-单条数据查询
  rpc Get{{.upperTableName}}(Get{{.upperTableName}}Req) returns (Get{{.upperTableName}}Reply) {
    option (google.api.http) = {get: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}/get"};
  }
  //{{.tableNameComment}}-列表数据查询
  rpc List{{.upperTableName}}(List{{.upperTableName}}Req) returns (List{{.upperTableName}}Reply) {
    option (google.api.http) = {
      post: "/{{.tableNameUnderScore}}/v1/{{.tableNameUnderScore}}/list",
      body: "*"
    };
  }
}