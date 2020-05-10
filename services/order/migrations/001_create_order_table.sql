-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS sell_order (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  asset_id uuid NOT NULL,
  amount INTEGER NOT NULL CHECK(amount > 0),
  price INTEGER NOT NULL CHECK(price > 0),
  seller_id INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS buy_order (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  asset_id uuid NOT NULL,
  amount INTEGER NOT NULL CHECK(amount > 0),
  price INTEGER NOT NULL CHECK(price > 0),
  buyer_id INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sell_order_history (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  asset_id uuid NOT NULL,
  amount INTEGER NOT NULL CHECK(amount > 0),
  price INTEGER NOT NULL CHECK(price > 0),
  seller_id INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS buy_order_history (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  asset_id uuid NOT NULL,
  amount INTEGER NOT NULL CHECK(amount > 0),
  price INTEGER NOT NULL CHECK(price > 0),
  sell_order_history_id uuid NOT NULL,
  buyer_id INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT sell_order_history_id_fkey FOREIGN KEY (sell_order_history_id)
    REFERENCES sell_order_history (id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION
);

---- create above / drop below ----

drop table buy_order_history;
drop table sell_order_history;
drop table sell_order;
drop table buy_order;
