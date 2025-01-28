--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2
-- Dumped by pg_dump version 17.2

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

--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


--
-- Name: appointment_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.appointment_status AS ENUM (
    'confirmed',
    'canceled',
    'completed'
);


ALTER TYPE public.appointment_status OWNER TO postgres;

--
-- Name: set_location_trigger(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.set_location_trigger() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

begin 

new.location := ST_Point(new.longitude::double precision, new.latitude::double precision, 4326);

return new;

end;

$$;


ALTER FUNCTION public.set_location_trigger() OWNER TO postgres;

--
-- Name: set_working_hours(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.set_working_hours() RETURNS trigger
    LANGUAGE plpgsql
    AS $$



begin

    new.working_hours := extract(hour from new.end_time) - extract(hour from new.start_time);

    return new;

end;

$$;


ALTER FUNCTION public.set_working_hours() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: appointments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.appointments (
    id integer NOT NULL,
    doctor_id integer,
    patient_id integer,
    ap_date date NOT NULL,
    ap_time time without time zone NOT NULL
);


ALTER TABLE public.appointments OWNER TO postgres;

--
-- Name: appointments_details; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.appointments_details (
    id integer NOT NULL,
    main_complaint text NOT NULL,
    presnet_illness text NOT NULL,
    past_illness text NOT NULL,
    family_history text NOT NULL,
    drug_history text NOT NULL,
    allergies text NOT NULL,
    doctor_id integer NOT NULL,
    patient_id integer NOT NULL,
    appointment_id integer NOT NULL
);


ALTER TABLE public.appointments_details OWNER TO postgres;

--
-- Name: appointments_details_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.appointments_details_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.appointments_details_id_seq OWNER TO postgres;

--
-- Name: appointments_details_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.appointments_details_id_seq OWNED BY public.appointments_details.id;


--
-- Name: appointments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.appointments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.appointments_id_seq OWNER TO postgres;

--
-- Name: appointments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.appointments_id_seq OWNED BY public.appointments.id;


--
-- Name: doctors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.doctors (
    id integer NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL,
    register_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.doctors OWNER TO postgres;

--
-- Name: doctors_availability; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.doctors_availability (
    id integer NOT NULL,
    doctor_id integer NOT NULL,
    start_time time with time zone NOT NULL,
    end_time time with time zone NOT NULL,
    date date NOT NULL,
    is_booked boolean DEFAULT false
);


ALTER TABLE public.doctors_availability OWNER TO postgres;

--
-- Name: doctors_availability_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.doctors_availability_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.doctors_availability_id_seq OWNER TO postgres;

--
-- Name: doctors_availability_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.doctors_availability_id_seq OWNED BY public.doctors_availability.id;


--
-- Name: doctors_clinic; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.doctors_clinic (
    id integer NOT NULL,
    governorate character varying NOT NULL,
    city character varying NOT NULL,
    address text NOT NULL,
    latitude text NOT NULL,
    longitude text NOT NULL,
    start_time time without time zone NOT NULL, -- state time should be with time zone
    end_time time without time zone NOT NULL,
    working_hours integer NOT NULL,
    holiday integer[] NOT NULL,
    location public.geometry(Point,4326) NOT NULL,
    price_per_appointment integer,
    patient_number_per_day integer,
    doctor_id integer NOT NULL,
    appointment_duration integer DEFAULT 15 NOT NULL
);


ALTER TABLE public.doctors_clinic OWNER TO postgres;

--
-- Name: doctors_clinic_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.doctors_clinic_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.doctors_clinic_id_seq OWNER TO postgres;

--
-- Name: doctors_clinic_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.doctors_clinic_id_seq OWNED BY public.doctors_clinic.id;


--
-- Name: doctors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.doctors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.doctors_id_seq OWNER TO postgres;

--
-- Name: doctors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.doctors_id_seq OWNED BY public.doctors.id;


--
-- Name: doctors_info; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.doctors_info (
    id integer NOT NULL,
    age integer NOT NULL,
    doctor_id integer NOT NULL,
    speciality character varying(20) NOT NULL,
    city character varying(20) NOT NULL,
    phone character varying(20) NOT NULL,
    img_url text,
    professional_statement character varying,
    experience character varying,
    satisfaction_score integer,
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    governorate character varying NOT NULL,
    CONSTRAINT doctors_info_satisfaction_score_check CHECK (((satisfaction_score > 0) AND (satisfaction_score <= 10)))
);


ALTER TABLE public.doctors_info OWNER TO postgres;

--
-- Name: doctors_info_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.doctors_info_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.doctors_info_id_seq OWNER TO postgres;

--
-- Name: doctors_info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.doctors_info_id_seq OWNED BY public.doctors_info.id;


--
-- Name: patients; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.patients (
    id integer NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL,
    register_date timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.patients OWNER TO postgres;

--
-- Name: patients_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.patients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.patients_id_seq OWNER TO postgres;

--
-- Name: patients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.patients_id_seq OWNED BY public.patients.id;


--
-- Name: patients_info; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.patients_info (
    id integer NOT NULL,
    age integer NOT NULL,
    patient_id integer,
    disease character varying(20) NOT NULL,
    location public.geometry(Point,4326),
    first_name character varying NOT NULL,
    last_name character varying NOT NULL,
    past_medical_history text,
    family_history text,
    drug_history text,
    allergies text
);


ALTER TABLE public.patients_info OWNER TO postgres;

--
-- Name: patients_info_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.patients_info_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.patients_info_id_seq OWNER TO postgres;

--
-- Name: patients_info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.patients_info_id_seq OWNED BY public.patients_info.id;


--
-- Name: questions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.questions (
    id integer NOT NULL,
    question text NOT NULL,
    answer text,
    doctor_id integer NOT NULL,
    patient_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    answered boolean DEFAULT false
);


ALTER TABLE public.questions OWNER TO postgres;

--
-- Name: questions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.questions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.questions_id_seq OWNER TO postgres;

--
-- Name: questions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.questions_id_seq OWNED BY public.questions.id;


--
-- Name: ratings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ratings (
    id integer NOT NULL,
    patient_id integer NOT NULL,
    doctor_id integer NOT NULL,
    rating integer NOT NULL,
    review text NOT NULL,
    CONSTRAINT ratings_rating_check CHECK (((rating > 0) AND (rating <= 10)))
);


ALTER TABLE public.ratings OWNER TO postgres;

--
-- Name: ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ratings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ratings_id_seq OWNER TO postgres;

--
-- Name: ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ratings_id_seq OWNED BY public.ratings.id;


--
-- Name: appointments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments ALTER COLUMN id SET DEFAULT nextval('public.appointments_id_seq'::regclass);


--
-- Name: appointments_details id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments_details ALTER COLUMN id SET DEFAULT nextval('public.appointments_details_id_seq'::regclass);


--
-- Name: doctors id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors ALTER COLUMN id SET DEFAULT nextval('public.doctors_id_seq'::regclass);


--
-- Name: doctors_availability id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_availability ALTER COLUMN id SET DEFAULT nextval('public.doctors_availability_id_seq'::regclass);


--
-- Name: doctors_clinic id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_clinic ALTER COLUMN id SET DEFAULT nextval('public.doctors_clinic_id_seq'::regclass);


--
-- Name: doctors_info id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_info ALTER COLUMN id SET DEFAULT nextval('public.doctors_info_id_seq'::regclass);


--
-- Name: patients id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients ALTER COLUMN id SET DEFAULT nextval('public.patients_id_seq'::regclass);


--
-- Name: patients_info id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients_info ALTER COLUMN id SET DEFAULT nextval('public.patients_info_id_seq'::regclass);


--
-- Name: questions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions ALTER COLUMN id SET DEFAULT nextval('public.questions_id_seq'::regclass);


--
-- Name: ratings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings ALTER COLUMN id SET DEFAULT nextval('public.ratings_id_seq'::regclass);


--
-- Name: appointments_details appointments_details_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments_details
    ADD CONSTRAINT appointments_details_pkey PRIMARY KEY (id);


--
-- Name: appointments appointments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_pkey PRIMARY KEY (id);


--
-- Name: doctors_info doctor; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_info
    ADD CONSTRAINT doctor UNIQUE (doctor_id);


--
-- Name: doctors_clinic doctor_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_clinic
    ADD CONSTRAINT doctor_id UNIQUE (doctor_id);


--
-- Name: doctors_availability doctors_availability_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_availability
    ADD CONSTRAINT doctors_availability_pkey PRIMARY KEY (id);


--
-- Name: doctors_clinic doctors_clinic_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_clinic
    ADD CONSTRAINT doctors_clinic_pkey PRIMARY KEY (id);


--
-- Name: doctors doctors_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors
    ADD CONSTRAINT doctors_email_key UNIQUE (email);


--
-- Name: doctors_info doctors_info_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_info
    ADD CONSTRAINT doctors_info_pkey PRIMARY KEY (id);


--
-- Name: doctors doctors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors
    ADD CONSTRAINT doctors_pkey PRIMARY KEY (id);


--
-- Name: patients patients_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients
    ADD CONSTRAINT patients_email_key UNIQUE (email);


--
-- Name: patients_info patients_info_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients_info
    ADD CONSTRAINT patients_info_pkey PRIMARY KEY (id);


--
-- Name: patients patients_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients
    ADD CONSTRAINT patients_pkey PRIMARY KEY (id);


--
-- Name: questions questions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pkey PRIMARY KEY (id);


--
-- Name: ratings ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_pkey PRIMARY KEY (id);


--
-- Name: doctors_clinic set_location_trigger; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_location_trigger BEFORE INSERT ON public.doctors_clinic FOR EACH ROW EXECUTE FUNCTION public.set_location_trigger();


--
-- Name: doctors_clinic set_working_hours; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_working_hours BEFORE INSERT ON public.doctors_clinic FOR EACH ROW EXECUTE FUNCTION public.set_working_hours();


--
-- Name: appointments_details appointments_details_appointment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments_details
    ADD CONSTRAINT appointments_details_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES public.appointments(id);


--
-- Name: appointments_details appointments_details_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments_details
    ADD CONSTRAINT appointments_details_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: appointments_details appointments_details_patient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments_details
    ADD CONSTRAINT appointments_details_patient_id_fkey FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: appointments appointments_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: appointments appointments_patient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_patient_id_fkey FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: doctors_availability doctors_availability_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_availability
    ADD CONSTRAINT doctors_availability_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: doctors_clinic doctors_clinic_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_clinic
    ADD CONSTRAINT doctors_clinic_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: doctors_info doctors_info_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.doctors_info
    ADD CONSTRAINT doctors_info_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: patients_info patients_info_patient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.patients_info
    ADD CONSTRAINT patients_info_patient_id_fkey FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: questions questions_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: questions questions_patient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_patient_id_fkey FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- Name: ratings ratings_doctor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);


--
-- Name: ratings ratings_patient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_patient_id_fkey FOREIGN KEY (patient_id) REFERENCES public.patients(id);


--
-- PostgreSQL database dump complete
--

