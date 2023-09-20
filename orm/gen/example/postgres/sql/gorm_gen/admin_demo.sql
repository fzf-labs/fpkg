--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-2.pgdg110+2)
-- Dumped by pg_dump version 15.4 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admin_demo; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admin_demo (
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


ALTER TABLE public.admin_demo OWNER TO postgres;

--
-- Name: TABLE admin_demo; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.admin_demo IS '系统-用户';


--
-- Name: COLUMN admin_demo.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.id IS '编号';


--
-- Name: COLUMN admin_demo.username; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.username IS '用户名';


--
-- Name: COLUMN admin_demo.password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.password IS '密码';


--
-- Name: COLUMN admin_demo.nickname; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.nickname IS '昵称';


--
-- Name: COLUMN admin_demo.avatar; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.avatar IS '头像';


--
-- Name: COLUMN admin_demo.gender; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.gender IS '0=保密 1=女 2=男';


--
-- Name: COLUMN admin_demo.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.email IS '邮件';


--
-- Name: COLUMN admin_demo.mobile; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.mobile IS '手机号';


--
-- Name: COLUMN admin_demo.job_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.job_id IS '岗位';


--
-- Name: COLUMN admin_demo.dept_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.dept_id IS '部门';


--
-- Name: COLUMN admin_demo.role_ids; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.role_ids IS '角色集';


--
-- Name: COLUMN admin_demo.salt; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.salt IS '盐值';


--
-- Name: COLUMN admin_demo.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.status IS '0=禁用 1=开启';


--
-- Name: COLUMN admin_demo.motto; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.motto IS '个性签名';


--
-- Name: COLUMN admin_demo.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.created_at IS '创建时间';


--
-- Name: COLUMN admin_demo.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.updated_at IS '更新时间';


--
-- Name: COLUMN admin_demo.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_demo.deleted_at IS '删除时间';


--
-- Name: admin_demo sys_admin_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin_demo
    ADD CONSTRAINT sys_admin_pkey PRIMARY KEY (id);


--
-- Name: sys_admin_username_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX sys_admin_username_idx ON public.admin_demo USING btree (username);


--
-- PostgreSQL database dump complete
--

