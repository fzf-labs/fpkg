CREATE TABLE public.admin_role_demo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pid uuid NOT NULL,
    name character varying(50) NOT NULL,
    remark character varying(200),
    status smallint NOT NULL,
    sort bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);
COMMENT ON TABLE public.admin_role_demo IS '系统-角色';
COMMENT ON COLUMN public.admin_role_demo.id IS '编号';
COMMENT ON COLUMN public.admin_role_demo.pid IS '父级id';
COMMENT ON COLUMN public.admin_role_demo.name IS '名称';
COMMENT ON COLUMN public.admin_role_demo.remark IS '备注';
COMMENT ON COLUMN public.admin_role_demo.status IS '0=禁用 1=开启';
COMMENT ON COLUMN public.admin_role_demo.sort IS '排序值';
COMMENT ON COLUMN public.admin_role_demo.created_at IS '创建时间';
COMMENT ON COLUMN public.admin_role_demo.updated_at IS '更新时间';
COMMENT ON COLUMN public.admin_role_demo.deleted_at IS '删除时间';
ALTER TABLE ONLY public.admin_role_demo
    ADD CONSTRAINT admin_role_demo_pkey PRIMARY KEY (id);
