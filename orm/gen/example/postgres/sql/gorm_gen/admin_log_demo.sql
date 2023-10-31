CREATE TABLE public.admin_log_demo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    admin_id uuid NOT NULL,
    ip character varying(32) NOT NULL,
    uri character varying(200) NOT NULL,
    useragent character varying(255),
    header json,
    req json,
    resp json,
    created_at timestamp with time zone NOT NULL
);
COMMENT ON TABLE public.admin_log_demo IS '系统-日志';
COMMENT ON COLUMN public.admin_log_demo.id IS '编号';
COMMENT ON COLUMN public.admin_log_demo.admin_id IS '管理员ID';
COMMENT ON COLUMN public.admin_log_demo.ip IS 'ip';
COMMENT ON COLUMN public.admin_log_demo.uri IS '请求路径';
COMMENT ON COLUMN public.admin_log_demo.useragent IS '浏览器标识';
COMMENT ON COLUMN public.admin_log_demo.header IS 'header';
COMMENT ON COLUMN public.admin_log_demo.req IS '请求数据';
COMMENT ON COLUMN public.admin_log_demo.resp IS '响应数据';
COMMENT ON COLUMN public.admin_log_demo.created_at IS '创建时间';
ALTER TABLE ONLY public.admin_log_demo
    ADD CONSTRAINT admin_log_demo_pkey PRIMARY KEY (id);
