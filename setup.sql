CREATE TABLE users (
    id            CHAR(27)     NOT NULL PRIMARY KEY,
    email         VARCHAR(100) NOT NULL,
    full_name     VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    gender        ENUM ('female', 'male') DEFAULT 'female',
    age           INT                     DEFAULT 18,
    location      POINT,
    CONSTRAINT email UNIQUE (email),
    INDEX filters(age, gender, location)
);

CREATE TABLE auth_tokens (
    id      CHAR(27) NOT NULL PRIMARY KEY,
    user_id CHAR(27) NOT NULL,
    CONSTRAINT foreign_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE swipes (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    swiper_id    CHAR(27) NOT NULL,
    recipient_id CHAR(27) NOT NULL,
    matched      BOOL DEFAULT FALSE,
    CONSTRAINT foreign_swiper_id FOREIGN KEY (swiper_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT foreign_recipient_id FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE INDEX swiper_recipient_unique (swiper_id, recipient_id)
);