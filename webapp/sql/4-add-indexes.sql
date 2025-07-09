-- Additional indexes for performance optimization

-- For ride matching: find unmatched rides quickly
ALTER TABLE rides ADD INDEX idx_rides_chair_id_created_at (chair_id, created_at);

-- For user ride queries
ALTER TABLE rides ADD INDEX idx_rides_user_id (user_id);

-- For chair ride queries  
ALTER TABLE rides ADD INDEX idx_rides_chair_id (chair_id);

-- For finding latest ride status efficiently
ALTER TABLE ride_statuses ADD INDEX idx_ride_statuses_created_at (ride_id, created_at DESC);

-- For status filtering
ALTER TABLE ride_statuses ADD INDEX idx_ride_statuses_status (status);

-- For chair activity queries
ALTER TABLE chairs ADD INDEX idx_chairs_is_active (is_active);

-- For latest chair location queries
ALTER TABLE chair_locations ADD INDEX idx_chair_locations_created_at (chair_id, created_at DESC);

-- For authentication queries
ALTER TABLE users ADD INDEX idx_users_access_token (access_token);
ALTER TABLE chairs ADD INDEX idx_chairs_access_token (access_token);
ALTER TABLE owners ADD INDEX idx_owners_access_token (access_token);

-- For coupon queries
ALTER TABLE coupons ADD INDEX idx_coupons_used_by (used_by);

-- Composite index for ride completion checks
ALTER TABLE ride_statuses ADD INDEX idx_ride_statuses_composite (ride_id, status, chair_sent_at);