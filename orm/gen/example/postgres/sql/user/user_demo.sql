
CREATE TABLE public.user_demo (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
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


ALTER TABLE public.user_demo OWNER TO postgres;

--
-- Name: TABLE user_demo; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.user_demo IS '用户信息表';


--
-- Name: COLUMN user_demo.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.id IS 'ID';


--
-- Name: COLUMN user_demo.uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.uid IS 'uid';


--
-- Name: COLUMN user_demo.username; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.username IS '用户账号';


--
-- Name: COLUMN user_demo.password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.password IS '密码';


--
-- Name: COLUMN user_demo.nickname; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.nickname IS '用户昵称';


--
-- Name: COLUMN user_demo.remark; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.remark IS '备注';


--
-- Name: COLUMN user_demo.dept_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.dept_id IS '部门ID';


--
-- Name: COLUMN user_demo.post_ids; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.post_ids IS '岗位编号数组';


--
-- Name: COLUMN user_demo.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.email IS '用户邮箱';


--
-- Name: COLUMN user_demo.mobile; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.mobile IS '手机号码';


--
-- Name: COLUMN user_demo.sex; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.sex IS '用户性别';


--
-- Name: COLUMN user_demo.avatar; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.avatar IS '头像地址';


--
-- Name: COLUMN user_demo.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.status IS '帐号状态（0正常 -1停用）';


--
-- Name: COLUMN user_demo.login_ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.login_ip IS '最后登录IP';


--
-- Name: COLUMN user_demo.login_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.login_date IS '最后登录时间';


--
-- Name: COLUMN user_demo.tenant_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.tenant_id IS '租户编号';


--
-- Name: COLUMN user_demo.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.created_at IS '创建时间';


--
-- Name: COLUMN user_demo.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.updated_at IS '更新时间';


--
-- Name: COLUMN user_demo.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_demo.deleted_at IS '删除时间';


--
-- Name: user_demo user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_demo
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user_tenant_id_dept_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_tenant_id_dept_id_idx ON public.user_demo USING btree (tenant_id, dept_id);


--
-- Name: user_uid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_uid_idx ON public.user_demo USING btree (uid);


--
-- Name: user_uid_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX user_uid_status_idx ON public.user_demo USING btree (uid, status);


--
-- Name: user_username_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_username_idx ON public.user_demo USING btree (username);


--
-- PostgreSQL database dump complete
--

