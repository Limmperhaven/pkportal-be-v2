DROP TABLE IF EXISTS public.user_profile_subjects;
DROP TABLE IF EXISTS public.user_profiles;
DROP TABLE IF EXISTS public.profile_subjects;
DROP TABLE IF EXISTS public.user_test_dates;
DROP TABLE IF EXISTS public.common_locations;
DROP TABLE IF EXISTS public.user_foreign_languages;
DROP TABLE IF EXISTS public.foreign_languages;
DROP TABLE IF EXISTS public.subjects;
DROP TABLE IF EXISTS public.profiles;
DROP TABLE IF EXISTS public.user_screenshots;
DROP TABLE IF EXISTS public.test_dates;
DROP TABLE IF EXISTS public.user_statuses;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.statuses;

DROP TYPE IF EXISTS screenshot_type;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_gender;
DROP TYPE IF EXISTS test_date_pub_status;

DROP EXTENSION IF EXISTS "uuid-ossp";