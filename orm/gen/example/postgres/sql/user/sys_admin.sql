CREATE TABLE public.sys_admin (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username character varying(50) DEFAULT ''::character varying NOT NULL,
    password character varying(128) NOT NULL,
    nickname character varying(50) NOT NULL,
    avatar character varying(255),
    gender smallint DEFAULT 0 NOT NULL,
    email character varying(50),
    mobile character varying(15),
    job_id uuid,
    dept_id uuid,
    role_ids json,
    salt character varying(32) NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    motto character varying(255),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.sys_admin OWNER TO postgres;

--
-- Name: TABLE sys_admin; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sys_admin IS '系统-用户';


--
-- Name: COLUMN sys_admin.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.id IS '编号';


--
-- Name: COLUMN sys_admin.username; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.username IS '用户名';


--
-- Name: COLUMN sys_admin.password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.password IS '密码';


--
-- Name: COLUMN sys_admin.nickname; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.nickname IS '昵称';


--
-- Name: COLUMN sys_admin.avatar; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.avatar IS '头像';


--
-- Name: COLUMN sys_admin.gender; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.gender IS '0=保密 1=女 2=男';


--
-- Name: COLUMN sys_admin.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.email IS '邮件';


--
-- Name: COLUMN sys_admin.mobile; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.mobile IS '手机号';


--
-- Name: COLUMN sys_admin.job_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.job_id IS '岗位';


--
-- Name: COLUMN sys_admin.dept_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.dept_id IS '部门';


--
-- Name: COLUMN sys_admin.role_ids; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.role_ids IS '角色集';


--
-- Name: COLUMN sys_admin.salt; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.salt IS '盐值';


--
-- Name: COLUMN sys_admin.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.status IS '0=禁用 1=开启';


--
-- Name: COLUMN sys_admin.motto; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.motto IS '个性签名';


--
-- Name: COLUMN sys_admin.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.created_at IS '创建时间';


--
-- Name: COLUMN sys_admin.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.updated_at IS '更新时间';


--
-- Name: COLUMN sys_admin.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_admin.deleted_at IS '删除时间';


--
-- Name: sys_admin sys_admin_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_admin
    ADD CONSTRAINT sys_admin_pkey PRIMARY KEY (id);


--
-- Name: sys_admin_username_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX sys_admin_username_idx ON public.sys_admin USING btree (username);


--
-- PostgreSQL database dump complete
--

