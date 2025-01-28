create table public.doctor_schedules (
    id serial primary key,
    date date not null,
    start_time time not null,
    end_time time not null,
    status varchar(20) not null,
    patient_id integer references patients(id),
    doctor_id integer references doctors(id),
    created_at timestamp with time zone not null default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone not null default CURRENT_TIMESTAMP
);

create index idx_doctor_scheduls_doctor_id on doctor_schedules (doctor_id);
create index idx_doctor_scheduls_date on doctor_schedules (date);
create index idx_doctor_scheduls_status on doctor_schedules (status);