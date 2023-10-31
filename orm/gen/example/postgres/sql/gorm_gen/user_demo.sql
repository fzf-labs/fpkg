CREATE TABLE public.user_demo (
    id bigint NOT NULL,
    uid character varying(64) NOT NULL,
    username character varying(30) NOT NULL,
    password character varying(100) NOT NULL,
    nickname character varying(30) NOT NULL,
    remark character varying(500),
    dept_id bigint,
    post_ids character varying(255),
    email character varying(50),
    mobile character varying(11),
    sex smallint,
    avatar character varying(100),
    status smallint DEFAULT 0 NOT NULL,
    login_ip character varying(50),
    login_date timestamp without time zone,
    tenant_id bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);
COMMENT ON TABLE public.user_demo IS '用户';
COMMENT ON COLUMN public.user_demo.id IS 'ID';
COMMENT ON COLUMN public.user_demo.uid IS 'uid';
COMMENT ON COLUMN public.user_demo.username IS '用户账号';
COMMENT ON COLUMN public.user_demo.password IS '密码';
COMMENT ON COLUMN public.user_demo.nickname IS '用户昵称';
COMMENT ON COLUMN public.user_demo.remark IS '备注';
COMMENT ON COLUMN public.user_demo.dept_id IS '部门ID';
COMMENT ON COLUMN public.user_demo.post_ids IS '岗位编号数组';
COMMENT ON COLUMN public.user_demo.email IS '用户邮箱';
COMMENT ON COLUMN public.user_demo.mobile IS '手机号码';
COMMENT ON COLUMN public.user_demo.sex IS '用户性别';
COMMENT ON COLUMN public.user_demo.avatar IS '头像地址';
COMMENT ON COLUMN public.user_demo.status IS '帐号状态（0正常 -1停用）';
COMMENT ON COLUMN public.user_demo.login_ip IS '最后登录IP';
COMMENT ON COLUMN public.user_demo.login_date IS '最后登录时间';
COMMENT ON COLUMN public.user_demo.tenant_id IS '租户编号';
COMMENT ON COLUMN public.user_demo.created_at IS '创建时间';
COMMENT ON COLUMN public.user_demo.updated_at IS '更新时间';
COMMENT ON COLUMN public.user_demo.deleted_at IS '删除时间';
ALTER TABLE ONLY public.user_demo
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);
CREATE INDEX user_tenant_id_dept_id_idx ON public.user_demo USING btree (tenant_id, dept_id);
CREATE UNIQUE INDEX user_uid_idx ON public.user_demo USING btree (uid);
CREATE UNIQUE INDEX user_uid_status_idx ON public.user_demo USING btree (uid, status);
CREATE INDEX user_username_idx ON public.user_demo USING btree (username);
