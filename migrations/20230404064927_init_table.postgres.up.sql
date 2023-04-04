CREATE TABLE users (
    "id" varchar(16) NOT NULL,
    "name" varchar(50) NOT NULL,
    "username" varchar(50) NOT NULL,
    "email" varchar(50) NOT NULL,
    "password" varchar(256) NOT NULL,
    "phone" varchar(16),
    "role_id" varchar(16) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE roles (
    "id" varchar(16) NOT NULL,
    "name" varchar(50) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE reset_token (
    "tokens" VARCHAR(255),
    "email" VARCHAR(50),
    "created_at" TIMESTAMP
);

ALTER TABLE users
    ADD CONSTRAINT "users_role_id_fkey" 
    FOREIGN KEY ("role_id") 
    REFERENCES roles("id");

INSERT INTO roles ("id", "name", "created_at", "updated_at") VALUES
    ('role001', 'Admin', NOW(), NOW()),
    ('role002', 'Client', NOW(), NOW());
