CREATE TABLE users (

id UUID PRIMARY KEY,

username TEXT UNIQUE NOT NULL,

email TEXT UNIQUE NOT NULL,

password_hash TEXT NOT NULL,

first_name TEXT,

last_name TEXT,

role TEXT,

is_active BOOLEAN DEFAULT TRUE,

ldap_user BOOLEAN DEFAULT FALSE,

last_login TIMESTAMP,

manager_id UUID,

created_at TIMESTAMP DEFAULT NOW(),

updated_at TIMESTAMP DEFAULT NOW()

);