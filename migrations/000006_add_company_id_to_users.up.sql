ALTER TABLE users ADD COLUMN company_id UUID;
ALTER TABLE users ADD CONSTRAINT fk_users_companies FOREIGN KEY (company_id) REFERENCES companies(id);
