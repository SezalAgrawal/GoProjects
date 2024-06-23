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

ALTER TABLE IF EXISTS ONLY public.installment_payments DROP CONSTRAINT IF EXISTS fk_rails_cf5f2f8949;
ALTER TABLE IF EXISTS ONLY public.loans DROP CONSTRAINT IF EXISTS fk_rails_c15c911198;
ALTER TABLE IF EXISTS ONLY public.installment_payments DROP CONSTRAINT IF EXISTS fk_rails_ba6cce9ae1;
ALTER TABLE IF EXISTS ONLY public.access_tokens DROP CONSTRAINT IF EXISTS fk_rails_96fc070778;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS fk_rails_3369e0d5fc;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS fk_rails_318345354e;
ALTER TABLE IF EXISTS ONLY public.installments DROP CONSTRAINT IF EXISTS fk_rails_1f10be3d66;
DROP INDEX IF EXISTS public.index_users_on_name_and_hashed_password;
DROP INDEX IF EXISTS public.index_user_roles_on_user_id;
DROP INDEX IF EXISTS public.index_loans_on_user_id;
DROP INDEX IF EXISTS public.index_installments_on_loan_id;
DROP INDEX IF EXISTS public.index_installment_payments_on_loan_id;
DROP INDEX IF EXISTS public.index_installment_payments_on_installment_id;
DROP INDEX IF EXISTS public.index_access_tokens_on_user_id;
DROP INDEX IF EXISTS public.index_access_tokens_on_token;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.roles DROP CONSTRAINT IF EXISTS roles_pkey;
ALTER TABLE IF EXISTS ONLY public.loans DROP CONSTRAINT IF EXISTS loans_pkey;
ALTER TABLE IF EXISTS ONLY public.installments DROP CONSTRAINT IF EXISTS installments_pkey;
ALTER TABLE IF EXISTS ONLY public.installment_payments DROP CONSTRAINT IF EXISTS installment_payments_pkey;
ALTER TABLE IF EXISTS ONLY public.ar_internal_metadata DROP CONSTRAINT IF EXISTS ar_internal_metadata_pkey;
ALTER TABLE IF EXISTS ONLY public.access_tokens DROP CONSTRAINT IF EXISTS access_tokens_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_roles;
DROP TABLE IF EXISTS public.schema_migrations;
DROP TABLE IF EXISTS public.roles;
DROP TABLE IF EXISTS public.loans;
DROP TABLE IF EXISTS public.installments;
DROP TABLE IF EXISTS public.installment_payments;
DROP TABLE IF EXISTS public.ar_internal_metadata;
DROP TABLE IF EXISTS public.access_tokens;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: access_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.access_tokens (
    id character varying NOT NULL,
    user_id character varying NOT NULL,
    token character varying NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    deleted_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: installment_payments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.installment_payments (
    id character varying NOT NULL,
    loan_id character varying NOT NULL,
    installment_id character varying NOT NULL,
    amount_in_cents integer NOT NULL,
    one_time_settlement_id character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: installments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.installments (
    id character varying NOT NULL,
    loan_id character varying NOT NULL,
    amount_in_cents integer NOT NULL,
    serial_no integer NOT NULL,
    status character varying NOT NULL,
    due_date timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: loans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loans (
    id character varying NOT NULL,
    user_id character varying NOT NULL,
    amount_in_cents integer NOT NULL,
    term integer NOT NULL,
    frequency_in_days integer NOT NULL,
    status character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id character varying NOT NULL,
    name character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    id character varying NOT NULL,
    user_id character varying NOT NULL,
    role_id character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id character varying NOT NULL,
    name character varying NOT NULL,
    hashed_password character varying NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


--
-- Name: access_tokens access_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_tokens
    ADD CONSTRAINT access_tokens_pkey PRIMARY KEY (id);


--
-- Name: ar_internal_metadata ar_internal_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ar_internal_metadata
    ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);


--
-- Name: installment_payments installment_payments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installment_payments
    ADD CONSTRAINT installment_payments_pkey PRIMARY KEY (id);


--
-- Name: installments installments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installments
    ADD CONSTRAINT installments_pkey PRIMARY KEY (id);


--
-- Name: loans loans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: index_access_tokens_on_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_access_tokens_on_token ON public.access_tokens USING btree (token);


--
-- Name: index_access_tokens_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_access_tokens_on_user_id ON public.access_tokens USING btree (user_id);


--
-- Name: index_installment_payments_on_installment_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_installment_payments_on_installment_id ON public.installment_payments USING btree (installment_id);


--
-- Name: index_installment_payments_on_loan_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_installment_payments_on_loan_id ON public.installment_payments USING btree (loan_id);


--
-- Name: index_installments_on_loan_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_installments_on_loan_id ON public.installments USING btree (loan_id);


--
-- Name: index_loans_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_loans_on_user_id ON public.loans USING btree (user_id);


--
-- Name: index_user_roles_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_user_roles_on_user_id ON public.user_roles USING btree (user_id);


--
-- Name: index_users_on_name_and_hashed_password; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_name_and_hashed_password ON public.users USING btree (name, hashed_password);


--
-- Name: installments fk_rails_1f10be3d66; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installments
    ADD CONSTRAINT fk_rails_1f10be3d66 FOREIGN KEY (loan_id) REFERENCES public.loans(id);


--
-- Name: user_roles fk_rails_318345354e; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_rails_318345354e FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_roles fk_rails_3369e0d5fc; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT fk_rails_3369e0d5fc FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: access_tokens fk_rails_96fc070778; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_tokens
    ADD CONSTRAINT fk_rails_96fc070778 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: installment_payments fk_rails_ba6cce9ae1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installment_payments
    ADD CONSTRAINT fk_rails_ba6cce9ae1 FOREIGN KEY (loan_id) REFERENCES public.loans(id);


--
-- Name: loans fk_rails_c15c911198; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT fk_rails_c15c911198 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: installment_payments fk_rails_cf5f2f8949; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.installment_payments
    ADD CONSTRAINT fk_rails_cf5f2f8949 FOREIGN KEY (installment_id) REFERENCES public.installments(id);


--
-- PostgreSQL database dump complete
--

SET search_path TO "$user", public;

INSERT INTO "schema_migrations" (version) VALUES
('20240623113056'),
('20240623113057'),
('20240623113058'),
('20240623113059'),
('20240623113060'),
('20240623113061'),
('20240623113062');


