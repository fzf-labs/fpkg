//{{.tableNameComment}}信息
message {{.upperTableName}}Info {
  {{.info}}
}

//请求-{{.tableNameComment}}-创建一条数据
message Create{{.upperTableName}}Req {
  {{.createReq}}
}

//响应-{{.tableNameComment}}-创建一条数据
message Create{{.upperTableName}}Reply {
  {{.createReply}}
}

//请求-{{.tableNameComment}}-更新一条数据
message Update{{.upperTableName}}Req {
  {{.updateReq}}
}

//响应-{{.tableNameComment}}-更新一条数据
message Update{{.upperTableName}}Reply {}

//请求-{{.tableNameComment}}-删除多条数据
message Delete{{.upperTableName}}Req {
  {{.deleteReq}}
}

//响应-{{.tableNameComment}}-删除多条数据
message Delete{{.upperTableName}}Reply {}

//请求-{{.tableNameComment}}-单条数据查询
message Get{{.upperTableName}}Req {
  {{.getReq}}
}

//响应-{{.tableNameComment}}-单条数据查询
message Get{{.upperTableName}}Reply {
  {{.upperTableName}}Info info = 1;
}

//请求-{{.tableNameComment}}-列表数据查询
message List{{.upperTableName}}Req {
  paginator.PaginatorReq paginator = 1; //分页
}

//响应-{{.tableNameComment}}-列表数据查询
message List{{.upperTableName}}Reply {
  paginator.PaginatorReply paginator = 1; // 分页
  repeated {{.upperTableName}}Info list = 2; // 列表数据
}
