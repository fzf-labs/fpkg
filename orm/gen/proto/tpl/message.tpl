//{{.tableNameComment}}信息
message {{.upperTableName}}Info {
{{.Info}}
}

//请求-{{.tableNameComment}}-创建一条数据
message {{.upperTableName}}StoreReq {
{{.StoreReq}}
}

//响应-{{.tableNameComment}}-创建一条数据
message {{.upperTableName}}StoreReply {
{{.StoreReply}}
}

//请求-{{.tableNameComment}}-删除多条数据
message {{.upperTableName}}DelReq {
{{.DelReq}}
}

//响应-{{.tableNameComment}}-删除多条数据
message {{.upperTableName}}DelReply {}

//请求-{{.tableNameComment}}-单条数据查询
message {{.upperTableName}}OneReq {
{{.OneReq}}
}

//响应-{{.tableNameComment}}-单条数据查询
message {{.upperTableName}}OneReply {
  {{.upperTableName}}Info info = 1;
}

//请求-{{.tableNameComment}}-列表数据查询
message {{.upperTableName}}ListReq {
  paginator.PaginatorReq paginator = 1; //分页
}

//响应-{{.tableNameComment}}-列表数据查询
message {{.upperTableName}}ListReply {
  paginator.PaginatorReply paginator = 1; // 分页
  repeated {{.upperTableName}}Info list = 2; // 列表数据
}
