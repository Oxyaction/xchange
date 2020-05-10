-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS account (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  balance integer NOT NULL DEFAULT 0 CHECK(balance >= 0)
);

CREATE TABLE IF NOT EXISTS asset (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR (255) NOT NULL,
  balance integer NOT NULL DEFAULT 0 CHECK(balance >= 0)
);

CREATE TABLE IF NOT EXISTS account_asset (
  account_id uuid NOT NULL,
  asset_id uuid NOT NULL,
  balance integer NOT NULL DEFAULT 0 CHECK(balance >= 0),
  PRIMARY KEY (account_id, asset_id),
  CONSTRAINT account_asset_account_id_fkey FOREIGN KEY (account_id)
    REFERENCES account (id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT account_asset_asset_id_fkey FOREIGN KEY (asset_id)
    REFERENCES asset (id) MATCH SIMPLE
    ON UPDATE NO ACTION ON DELETE NO ACTION
);

---- create above / drop below ----

drop table account_asset;
drop table asset;
drop table account;
