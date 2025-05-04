-- Add status column to listings table
ALTER TABLE listings ADD COLUMN status TEXT NOT NULL DEFAULT 'available';

-- Create purchases table to track buying history
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES listings(id) ON DELETE CASCADE,
    buyer_id INT REFERENCES users(id) ON DELETE CASCADE,
    seller_id INT REFERENCES users(id) ON DELETE CASCADE,
    price DECIMAL NOT NULL,
    purchased_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (listing_id)
);

-- Add index for better performance
CREATE INDEX purchases_buyer_id_idx ON purchases(buyer_id);
CREATE INDEX purchases_seller_id_idx ON purchases(seller_id);
CREATE INDEX listings_status_idx ON listings(status);

-- Add comment about possible status values
COMMENT ON COLUMN listings.status IS 'Possible values: available, sold'; 