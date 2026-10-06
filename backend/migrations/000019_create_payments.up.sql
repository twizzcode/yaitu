CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL DEFAULT 'midtrans',
    order_id TEXT NOT NULL UNIQUE,
    transaction_id TEXT,
    payment_type TEXT,
    gross_amount INTEGER NOT NULL CHECK (gross_amount >= 0),
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (
            status IN (
                'pending',
                'settlement',
                'capture',
                'expire',
                'cancel',
                'deny',
                'failure',
                'refund'
            )
        ),
    qr_url TEXT,
    expires_at TIMESTAMPTZ,
    raw_charge JSONB,
    raw_notification JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX payments_booking_id_idx
ON payments(booking_id);

CREATE INDEX payments_user_id_idx
ON payments(user_id);

CREATE INDEX payments_status_idx
ON payments(status);
