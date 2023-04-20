ALTER TABLE IF EXISTS public.user_test_dates
    DROP CONSTRAINT user_test_dates_pkey;
ALTER TABLE IF EXISTS public.user_test_dates
    ADD PRIMARY KEY (user_id, education_year);

DROP TABLE IF EXISTS public.user_exam_results;