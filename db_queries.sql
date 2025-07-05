-- データベース接続情報
-- Host: localhost
-- Port: 5432
-- Database: yelp_sample
-- Username: postgres
-- Password: postgres

-- 全テーブルの構造確認
\d users;
\d businesses;
\d reviews;

-- 全ユーザー一覧
SELECT * FROM users ORDER BY id;

-- 全企業一覧（評価順）
SELECT 
    id,
    name,
    category,
    rating,
    review_count,
    address,
    phone
FROM businesses 
ORDER BY rating DESC, review_count DESC;

-- 全レビュー一覧（企業名とユーザー名付き）
SELECT 
    r.id,
    b.name AS business_name,
    u.name AS user_name,
    r.rating,
    r.text,
    r.created_at
FROM reviews r
JOIN businesses b ON r.business_id = b.id
JOIN users u ON r.user_id = u.id
ORDER BY r.created_at DESC;

-- カテゴリ別企業数
SELECT 
    category,
    COUNT(*) as business_count,
    AVG(rating) as avg_rating
FROM businesses 
GROUP BY category
ORDER BY business_count DESC;

-- 評価の高い企業TOP3
SELECT 
    name,
    category,
    rating,
    review_count,
    address
FROM businesses 
WHERE review_count > 0
ORDER BY rating DESC, review_count DESC
LIMIT 3;

-- 特定企業（ラーメン山田）の詳細情報
SELECT 
    b.*,
    COUNT(r.id) as total_reviews,
    AVG(r.rating) as calculated_avg_rating
FROM businesses b
LEFT JOIN reviews r ON b.id = r.business_id
WHERE b.name = 'ラーメン山田'
GROUP BY b.id;

-- 特定企業のレビュー詳細
SELECT 
    r.id,
    u.name AS reviewer,
    r.rating,
    r.text,
    r.created_at
FROM reviews r
JOIN users u ON r.user_id = u.id
WHERE r.business_id = 1
ORDER BY r.created_at DESC;

-- ユーザー別レビュー数
SELECT 
    u.name,
    u.email,
    COUNT(r.id) as review_count,
    AVG(r.rating) as avg_rating_given
FROM users u
LEFT JOIN reviews r ON u.id = r.user_id
GROUP BY u.id, u.name, u.email
ORDER BY review_count DESC;

-- 評価分布
SELECT 
    rating,
    COUNT(*) as count
FROM reviews
GROUP BY rating
ORDER BY rating DESC;

-- 最近のレビュー（直近10件）
SELECT 
    b.name AS business_name,
    u.name AS user_name,
    r.rating,
    r.text,
    r.created_at
FROM reviews r
JOIN businesses b ON r.business_id = b.id
JOIN users u ON r.user_id = u.id
ORDER BY r.created_at DESC
LIMIT 10;