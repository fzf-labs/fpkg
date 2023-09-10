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
-- Name: sys_api; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sys_api (
    id uuid NOT NULL,
    permission_id uuid NOT NULL,
    method character varying(32) NOT NULL,
    path character varying(255) NOT NULL,
    "desc" character varying(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.sys_api OWNER TO postgres;

--
-- Name: TABLE sys_api; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sys_api IS '系统-接口';


--
-- Name: COLUMN sys_api.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.id IS '编号';


--
-- Name: COLUMN sys_api.permission_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.permission_id IS '权限Id';


--
-- Name: COLUMN sys_api.method; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.method IS '方法';


--
-- Name: COLUMN sys_api.path; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.path IS '路径';


--
-- Name: COLUMN sys_api."desc"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api."desc" IS '描述';


--
-- Name: COLUMN sys_api.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.created_at IS '创建时间';


--
-- Name: COLUMN sys_api.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.updated_at IS '更新时间';


--
-- Name: COLUMN sys_api.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_api.deleted_at IS '删除时间';


--
-- Name: sys_api sys_api_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_api
    ADD CONSTRAINT sys_api_pkey PRIMARY KEY (id);


--
-- Name: sys_api_permission_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX sys_api_permission_id_idx ON public.sys_api USING btree (permission_id);


--
-- PostgreSQL database dump complete
--

