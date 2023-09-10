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
-- Name: sys_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sys_log (
    id uuid NOT NULL,
    admin_id uuid NOT NULL,
    ip character varying(32) NOT NULL,
    uri character varying(200) NOT NULL,
    useragent character varying(255),
    header json,
    req json,
    resp json,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.sys_log OWNER TO postgres;

--
-- Name: TABLE sys_log; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sys_log IS '系统-日志';


--
-- Name: COLUMN sys_log.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.id IS '编号';


--
-- Name: COLUMN sys_log.admin_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.admin_id IS '管理员ID';


--
-- Name: COLUMN sys_log.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.ip IS 'ip';


--
-- Name: COLUMN sys_log.uri; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.uri IS '请求路径';


--
-- Name: COLUMN sys_log.useragent; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.useragent IS '浏览器标识';


--
-- Name: COLUMN sys_log.header; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.header IS 'header';


--
-- Name: COLUMN sys_log.req; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.req IS '请求数据';


--
-- Name: COLUMN sys_log.resp; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.resp IS '响应数据';


--
-- Name: COLUMN sys_log.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_log.created_at IS '创建时间';


--
-- Name: sys_log sys_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_log
    ADD CONSTRAINT sys_log_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

