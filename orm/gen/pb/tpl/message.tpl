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
  int32 page = 1 [(validate.rules).int32 = {gte: 1}]; //页码
  int32 pageSize = 2 [(validate.rules).int32 = {
    gte: 1,
    lte: 1000
  }]; //页数
}

//响应-{{.tableNameComment}}-列表数据查询
message {{.upperTableName}}ListReply {
  int32 page = 1; // 页码
  int32 pageSize = 2; // 页数
  int32 total = 3;// 总数
  repeated {{.upperTableName}}Info list = 4; //列表数据
}
