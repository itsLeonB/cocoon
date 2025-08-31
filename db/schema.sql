--
-- PostgreSQL database dump
--

\restrict dTxp4tVQPAEmQumAD0zuWmct421rM79ddcbQrhnbedD1ILeO6Y0fzBJqzt41Qfu

-- Dumped from database version 17.5 (1b53132)
-- Dumped by pg_dump version 17.6 (Ubuntu 17.6-1.pgdg24.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE IF EXISTS ONLY public.friendships DROP CONSTRAINT IF EXISTS friendships_profile_id2_fkey;
ALTER TABLE IF EXISTS ONLY public.friendships DROP CONSTRAINT IF EXISTS friendships_profile_id1_fkey;
DROP INDEX IF EXISTS public.user_profiles_user_id_idx;
DROP INDEX IF EXISTS public.user_profiles_name_idx;
DROP INDEX IF EXISTS public.friendships_type_idx;
DROP INDEX IF EXISTS public.friendships_profile_id2_idx;
DROP INDEX IF EXISTS public.friendships_profile_id1_idx;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.user_profiles DROP CONSTRAINT IF EXISTS user_profiles_pkey;
ALTER TABLE IF EXISTS ONLY public.friendships DROP CONSTRAINT IF EXISTS unique_friendship;
ALTER TABLE IF EXISTS ONLY public.friendships DROP CONSTRAINT IF EXISTS friendships_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_profiles;
DROP TABLE IF EXISTS public.friendships;
DROP TYPE IF EXISTS public.friendship_type;
--
-- Name: friendship_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.friendship_type AS ENUM (
    'REAL',
    'ANON'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: friendships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.friendships (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    profile_id1 uuid NOT NULL,
    profile_id2 uuid NOT NULL,
    type public.friendship_type NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT profile_order CHECK ((profile_id1 < profile_id2))
);


--
-- Name: user_profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_profiles (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: COLUMN user_profiles.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_profiles.user_id IS 'Nullable. Can be NULL for peers who do not have an account in the app';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: friendships friendships_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_pkey PRIMARY KEY (id);


--
-- Name: friendships unique_friendship; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT unique_friendship UNIQUE (profile_id1, profile_id2);


--
-- Name: user_profiles user_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_profiles
    ADD CONSTRAINT user_profiles_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: friendships_profile_id1_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX friendships_profile_id1_idx ON public.friendships USING btree (profile_id1);


--
-- Name: friendships_profile_id2_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX friendships_profile_id2_idx ON public.friendships USING btree (profile_id2);


--
-- Name: friendships_type_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX friendships_type_idx ON public.friendships USING btree (type);


--
-- Name: user_profiles_name_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX user_profiles_name_idx ON public.user_profiles USING btree (name);


--
-- Name: user_profiles_user_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX user_profiles_user_id_idx ON public.user_profiles USING btree (user_id);


--
-- Name: friendships friendships_profile_id1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_profile_id1_fkey FOREIGN KEY (profile_id1) REFERENCES public.user_profiles(id);


--
-- Name: friendships friendships_profile_id2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_profile_id2_fkey FOREIGN KEY (profile_id2) REFERENCES public.user_profiles(id);


--
-- PostgreSQL database dump complete
--

\unrestrict dTxp4tVQPAEmQumAD0zuWmct421rM79ddcbQrhnbedD1ILeO6Y0fzBJqzt41Qfu

