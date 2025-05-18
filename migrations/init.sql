CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    name          TEXT,
    city          TEXT,
    avatar        TEXT,
    is_verified   BOOLEAN   DEFAULT FALSE,
    created_at    TIMESTAMP DEFAULT NOW()
);
-- Добавляем поле last_name в таблицу users
ALTER TABLE users ADD COLUMN last_name TEXT;

CREATE TABLE two_factor_codes
(
    id         SERIAL PRIMARY KEY,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE,
    code       TEXT      NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE categories
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE listings
(
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users (id) ON DELETE CASCADE,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    price       DECIMAL NOT NULL,
    condition   TEXT    NOT NULL, -- новое, хорошее, среднее, плохое
    city        TEXT    NOT NULL,
    category_id INT REFERENCES categories (id),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

ALTER TABLE listings ADD COLUMN status TEXT NOT NULL DEFAULT 'active';

CREATE TABLE listing_images
(
    id         SERIAL PRIMARY KEY,
    listing_id INT REFERENCES listings (id) ON DELETE CASCADE,
    image_path TEXT NOT NULL,
    is_main    BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE favorites
(
    id         SERIAL PRIMARY KEY,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE,
    listing_id INT REFERENCES listings (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, listing_id)
);

CREATE TABLE chats
(
    id         SERIAL PRIMARY KEY,
    listing_id INT REFERENCES listings (id) ON DELETE CASCADE,
    buyer_id   INT REFERENCES users (id) ON DELETE CASCADE,
    seller_id  INT REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (listing_id, buyer_id, seller_id)
);

CREATE TABLE messages
(
    id         SERIAL PRIMARY KEY,
    chat_id    INT REFERENCES chats (id) ON DELETE CASCADE,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    is_read    BOOLEAN   DEFAULT FALSE
);

-- Вставляем начальные категории
INSERT INTO categories (name)
VALUES ('Диваны и кресла'),
       ('Столы и стулья'),
       ('Шкафы и комоды'),
       ('Кровати и матрасы'),
       ('Другое');

-- Create purchases table to track buying history
CREATE TABLE purchases
(
    id           SERIAL PRIMARY KEY,
    listing_id   INT REFERENCES listings (id) ON DELETE CASCADE,
    buyer_id     INT REFERENCES users (id) ON DELETE CASCADE,
    seller_id    INT REFERENCES users (id) ON DELETE CASCADE,
    price        DECIMAL NOT NULL,
    purchased_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (listing_id)
);

-- Add index for better performance
CREATE INDEX purchases_buyer_id_idx ON purchases (buyer_id);
CREATE INDEX purchases_seller_id_idx ON purchases (seller_id);
CREATE INDEX listings_status_idx ON listings (status);