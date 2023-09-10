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
-- Name: sys_permission; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sys_permission (
    id uuid NOT NULL,
    pid uuid NOT NULL,
    type character varying(255) NOT NULL,
    title character varying(50) NOT NULL,
    name character varying(50) NOT NULL,
    path character varying(100) NOT NULL,
    icon character varying(50) NOT NULL,
    menu_type character varying(255),
    url character varying(255) NOT NULL,
    component character varying(100) NOT NULL,
    extend character varying(255) NOT NULL,
    remark character varying(255) NOT NULL,
    sort bigint NOT NULL,
    status smallint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.sys_permission OWNER TO postgres;

--
-- Name: TABLE sys_permission; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sys_permission IS '菜单和权限规则表';


--
-- Name: COLUMN sys_permission.pid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.pid IS '上级菜单';


--
-- Name: COLUMN sys_permission.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.type IS '类型:menu_dir=菜单目录,menu=菜单项,button=页面按钮';


--
-- Name: COLUMN sys_permission.title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.title IS '标题';


--
-- Name: COLUMN sys_permission.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.name IS '规则名称';


--
-- Name: COLUMN sys_permission.path; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.path IS '路由路径';


--
-- Name: COLUMN sys_permission.icon; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.icon IS '图标';


--
-- Name: COLUMN sys_permission.menu_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.menu_type IS '菜单类型:tab=选项卡,link=链接,iframe=Iframe';


--
-- Name: COLUMN sys_permission.url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.url IS 'Url';


--
-- Name: COLUMN sys_permission.component; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.component IS '组件路径';


--
-- Name: COLUMN sys_permission.extend; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.extend IS '扩展属性:none=无,add_rules_only=只添加为路由,add_menu_only=只添加为菜单';


--
-- Name: COLUMN sys_permission.remark; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.remark IS '备注';


--
-- Name: COLUMN sys_permission.sort; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.sort IS '权重(排序)';


--
-- Name: COLUMN sys_permission.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.status IS '0=禁用 1=开启';


--
-- Name: COLUMN sys_permission.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.created_at IS '创建时间';


--
-- Name: COLUMN sys_permission.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.updated_at IS '更新时间';


--
-- Name: COLUMN sys_permission.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sys_permission.deleted_at IS '删除时间';


--
-- Name: sys_permission sys_permission_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_permission
    ADD CONSTRAINT sys_permission_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

