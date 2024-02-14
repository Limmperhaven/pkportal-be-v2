ALTER TABLE IF EXISTS public.user_test_dates
    DROP CONSTRAINT user_test_dates_pkey;
ALTER TABLE IF EXISTS public.user_test_dates
    ADD PRIMARY KEY (user_id, test_date_id);

CREATE TABLE IF NOT EXISTS public.user_exam_results
(
    user_id                bigint REFERENCES public.users (id)        NOT NULL,
    test_date_id           bigint REFERENCES public.test_dates (id)   NOT NULL,
    education_year         smallint                                   NOT NULL,
    russian_language_grade integer,
    math_grade             integer,
    foreign_language_grade integer,
    first_profile_grade    integer,
    second_profile_grade   integer,
    PRIMARY KEY (user_id, test_date_id)
);