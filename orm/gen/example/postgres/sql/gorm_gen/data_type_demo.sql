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
-- Name: data_type_demo; Type: TABLE; Schema: public; Owner: postgres
--

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


ALTER TABLE public.data_type_demo OWNER TO postgres;

--
-- Name: TABLE data_type_demo; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.data_type_demo IS '数据类型示例';


--
-- Name: COLUMN data_type_demo.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.id IS 'ID';


--
-- Name: COLUMN data_type_demo.data_type_bool; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_bool IS '数据类型 bool';


--
-- Name: COLUMN data_type_demo.data_type_int2; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_int2 IS '数据类型 int2';


--
-- Name: COLUMN data_type_demo.data_type_int8; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_int8 IS '数据类型 int8';


--
-- Name: COLUMN data_type_demo.data_type_varchar; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_varchar IS '数据类型 varchar';


--
-- Name: COLUMN data_type_demo.data_type_text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_text IS '数据类型 text';


--
-- Name: COLUMN data_type_demo.data_type_json; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_json IS '数据类型 json';


--
-- Name: COLUMN data_type_demo.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.created_at IS '创建时间';


--
-- Name: COLUMN data_type_demo.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.updated_at IS '更新时间';


--
-- Name: COLUMN data_type_demo.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.deleted_at IS '删除时间';


--
-- Name: COLUMN data_type_demo.data_type_time_null; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_time_null IS '数据类型 time null';


--
-- Name: COLUMN data_type_demo.data_type_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_time IS '数据类型 time not null';


--
-- Name: COLUMN data_type_demo.data_type_jsonb; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo.data_type_jsonb IS '数据类型 jsonb';


--
-- Name: COLUMN data_type_demo._id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo._id IS '验证下划线';


--
-- Name: COLUMN data_type_demo."cacheKey"; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.data_type_demo."cacheKey" IS '特殊保留字段名称';


--
-- Name: data_type_demo data_type_demo_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.data_type_demo
    ADD CONSTRAINT data_type_demo_pkey PRIMARY KEY (id);


--
-- Name: data_type_demo__id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX data_type_demo__id_idx ON public.data_type_demo USING btree (_id);


--
-- Name: data_type_demo_batch_api_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX data_type_demo_batch_api_idx ON public.data_type_demo USING btree (batch_api);


--
-- Name: data_type_demo_cacheKey_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "data_type_demo_cacheKey_idx" ON public.data_type_demo USING btree ("cacheKey");


--
-- Name: data_type_demo_data_type_bytea_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX data_type_demo_data_type_bytea_idx ON public.data_type_demo USING btree (data_type_bytea);


--
-- Name: data_type_demo_data_type_jsonb_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX data_type_demo_data_type_jsonb_idx ON public.data_type_demo USING btree (data_type_jsonb);


--
-- Name: data_type_demo_data_type_time_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX data_type_demo_data_type_time_idx ON public.data_type_demo USING btree (data_type_time);


--
-- Name: data_type_demo_data_type_time_null_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX data_type_demo_data_type_time_null_idx ON public.data_type_demo USING btree (data_type_time_null);


--
-- Name: data_type_demo_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX data_type_demo_deleted_at_idx ON public.data_type_demo USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

