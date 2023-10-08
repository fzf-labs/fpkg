CREATE TABLE public.data_type_demo (
    id uuid NOT NULL,
    data_type_bool boolean,
    data_type_int2 smallint,
    data_type_int8 bigint,
    data_type_varchar character varying DEFAULT 'test'::character varying,
    data_type_text text,
    data_type_json json,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    data_type_time_null timestamp with time zone,
    data_type_time timestamp with time zone NOT NULL,
    data_type_jsonb jsonb,
    data_type_date date,
    data_type_float4 real,
    data_type_float8 double precision,
    _id uuid,
    "cacheKey" character varying,
    data_type_timestamp timestamp without time zone,
    data_type_bytea bytea,
    data_type_numeric numeric,
    data_type_interval interval,
    batch_api character varying
);
COMMENT ON TABLE public.data_type_demo IS '数据类型示例';
COMMENT ON COLUMN public.data_type_demo.id IS 'ID';
COMMENT ON COLUMN public.data_type_demo.data_type_bool IS '数据类型 bool';
COMMENT ON COLUMN public.data_type_demo.data_type_int2 IS '数据类型 int2';
COMMENT ON COLUMN public.data_type_demo.data_type_int8 IS '数据类型 int8';
COMMENT ON COLUMN public.data_type_demo.data_type_varchar IS '数据类型 varchar';
COMMENT ON COLUMN public.data_type_demo.data_type_text IS '数据类型 text';
COMMENT ON COLUMN public.data_type_demo.data_type_json IS '数据类型 json';
COMMENT ON COLUMN public.data_type_demo.created_at IS '创建时间';
COMMENT ON COLUMN public.data_type_demo.updated_at IS '更新时间';
COMMENT ON COLUMN public.data_type_demo.deleted_at IS '删除时间';
COMMENT ON COLUMN public.data_type_demo.data_type_time_null IS '数据类型 time null';
COMMENT ON COLUMN public.data_type_demo.data_type_time IS '数据类型 time not null';
COMMENT ON COLUMN public.data_type_demo.data_type_jsonb IS '数据类型 jsonb';
COMMENT ON COLUMN public.data_type_demo._id IS '验证下划线';
COMMENT ON COLUMN public.data_type_demo."cacheKey" IS '特殊保留字段名称';
ALTER TABLE ONLY public.data_type_demo
    ADD CONSTRAINT data_type_demo_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX data_type_demo__id_idx ON public.data_type_demo USING btree (_id);
CREATE INDEX data_type_demo_batch_api_idx ON public.data_type_demo USING btree (batch_api);
CREATE INDEX "data_type_demo_cacheKey_idx" ON public.data_type_demo USING btree ("cacheKey");
CREATE INDEX data_type_demo_data_type_bytea_idx ON public.data_type_demo USING btree (data_type_bytea);
CREATE INDEX data_type_demo_data_type_jsonb_idx ON public.data_type_demo USING btree (data_type_jsonb);
CREATE INDEX data_type_demo_data_type_time_idx ON public.data_type_demo USING btree (data_type_time);
CREATE UNIQUE INDEX data_type_demo_data_type_time_null_idx ON public.data_type_demo USING btree (data_type_time_null);
CREATE INDEX data_type_demo_deleted_at_idx ON public.data_type_demo USING btree (deleted_at);
