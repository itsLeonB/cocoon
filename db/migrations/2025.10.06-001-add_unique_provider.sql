ALTER TABLE oauth_accounts
ADD CONSTRAINT oauth_accounts_provider_provider_id_unique
UNIQUE (provider, provider_id);
