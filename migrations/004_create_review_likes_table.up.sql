CREATE TABLE IF NOT EXISTS review_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    UNIQUE(user_id, review_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_likes_user_id ON review_likes(user_id);

CREATE INDEX IF NOT EXISTS idx_review_likes_review_id ON review_likes(review_id);