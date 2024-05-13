-- Если база данных уже существует, не создавайте её заново
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_database WHERE datname = 'url_shorter'
    ) THEN
        CREATE DATABASE your_database;
    END IF;
END $$;
