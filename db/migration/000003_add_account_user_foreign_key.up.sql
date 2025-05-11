ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");