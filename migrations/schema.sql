--
-- PostgreSQL database dump
--

-- Dumped from database version 13.4 (Ubuntu 13.4-1.pgdg20.04+1)
-- Dumped by pg_dump version 13.4 (Ubuntu 13.4-1.pgdg20.04+1)

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
-- Name: genres; Type: TABLE; Schema: public; Owner: movie_user
--

CREATE TABLE public.genres (
    id integer NOT NULL,
    genre_name character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    json_name character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.genres OWNER TO movie_user;

--
-- Name: genres_id_seq; Type: SEQUENCE; Schema: public; Owner: movie_user
--

CREATE SEQUENCE public.genres_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.genres_id_seq OWNER TO movie_user;

--
-- Name: genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: movie_user
--

ALTER SEQUENCE public.genres_id_seq OWNED BY public.genres.id;


--
-- Name: movies; Type: TABLE; Schema: public; Owner: movie_user
--

CREATE TABLE public.movies (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    description character varying(1000) NOT NULL,
    year integer DEFAULT 0 NOT NULL,
    release_date timestamp without time zone NOT NULL,
    runtime integer DEFAULT 0 NOT NULL,
    rating integer DEFAULT 0 NOT NULL,
    mpaa_rating character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    poster character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.movies OWNER TO movie_user;

--
-- Name: movies_genres; Type: TABLE; Schema: public; Owner: movie_user
--

CREATE TABLE public.movies_genres (
    id integer NOT NULL,
    movie_id integer DEFAULT 0 NOT NULL,
    genre_id integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.movies_genres OWNER TO movie_user;

--
-- Name: movies_genres_id_seq; Type: SEQUENCE; Schema: public; Owner: movie_user
--

CREATE SEQUENCE public.movies_genres_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.movies_genres_id_seq OWNER TO movie_user;

--
-- Name: movies_genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: movie_user
--

ALTER SEQUENCE public.movies_genres_id_seq OWNED BY public.movies_genres.id;


--
-- Name: movies_id_seq; Type: SEQUENCE; Schema: public; Owner: movie_user
--

CREATE SEQUENCE public.movies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.movies_id_seq OWNER TO movie_user;

--
-- Name: movies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: movie_user
--

ALTER SEQUENCE public.movies_id_seq OWNED BY public.movies.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: movie_user
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO movie_user;

--
-- Name: genres id; Type: DEFAULT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.genres ALTER COLUMN id SET DEFAULT nextval('public.genres_id_seq'::regclass);


--
-- Name: movies id; Type: DEFAULT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies ALTER COLUMN id SET DEFAULT nextval('public.movies_id_seq'::regclass);


--
-- Name: movies_genres id; Type: DEFAULT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies_genres ALTER COLUMN id SET DEFAULT nextval('public.movies_genres_id_seq'::regclass);


--
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);


--
-- Name: movies_genres movies_genres_pkey; Type: CONSTRAINT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT movies_genres_pkey PRIMARY KEY (id);


--
-- Name: movies movies_pkey; Type: CONSTRAINT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: movie_user
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: movies_genres movies_genres_genres_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT movies_genres_genres_id_fk FOREIGN KEY (genre_id) REFERENCES public.genres(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: movies_genres movies_genres_movies_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: movie_user
--

ALTER TABLE ONLY public.movies_genres
    ADD CONSTRAINT movies_genres_movies_id_fk FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

