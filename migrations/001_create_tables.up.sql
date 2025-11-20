CREATE TABLE IF NOT EXISTS users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username VARCHAR(50) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS movies (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(255) NOT NULL,
                        description TEXT,
                        release_date DATE,
                        director VARCHAR(100),
                        actors TEXT,
                        genres VARCHAR(200),
                        UNIQUE(title, release_date)
);
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);
CREATE INDEX IF NOT EXISTS idx_movies_release_date ON movies(release_date);

CREATE TABLE IF NOT EXISTS user_movies (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
                             list_type VARCHAR(20) CHECK (list_type IN ('favorite', 'watchlist', NULL)),
                             user_rating INTEGER CHECK (user_rating >= 1 AND user_rating <= 10 OR user_rating = 0),
                             UNIQUE(user_id, movie_id)
);
CREATE INDEX IF NOT EXISTS idx_user_movies_user_id ON user_movies(user_id);