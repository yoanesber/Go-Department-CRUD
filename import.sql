-- Description: SQL script to import initial user data into the database.
INSERT INTO users (username,"password",email,firstname,lastname,is_enabled,is_account_non_expired,is_account_non_locked,is_credentials_non_expired,is_deleted,account_expiration_date,credentials_expiration_date,user_type,last_login,created_by,updated_by) VALUES
	 ('admin','$2a$10$eP5Sddi7Q5Jv6seppeF93.XsWGY8r4PnsqprWGb5AxsZ9TpwULIGa','admin@mygmail.com','Admin','Admin',true,true,true,true,false,'2025-04-23 21:52:38.000','2025-02-28 01:58:35.000','USER_ACCOUNT','2025-02-11 22:54:32.000',0,0),
	 ('userone','$2a$10$eP5Sddi7Q5Jv6seppeF93.XsWGY8r4PnsqprWGb5AxsZ9TpwULIGa','userone@mygmail.com','User','One',true,true,true,true,false,'2025-07-14 19:50:56.000','2025-05-11 22:57:25.000','USER_ACCOUNT','2025-02-10 14:53:04.000',1,1);


-- Description: SQL script to import initial role data into the database.
INSERT INTO roles ("name") VALUES
	 ('ROLE_USER'),
	 ('ROLE_MODERATOR'),
	 ('ROLE_ADMIN');

-- Description: SQL script to import initial user-role mapping data into the database.
INSERT INTO user_roles (user_id,role_id) VALUES
	 (1,3),
	 (2,1);

-- Description: SQL script to import initial department data into the database.
INSERT INTO department (id,dept_name,active,created_by,updated_by) VALUES
	 ('d001','Marketing',true,1,1),
	 ('d002','Finance',true,1,1),
	 ('d003','Human Resources',true,1,1),
	 ('d004','Production',true,1,1),
	 ('d005','Development',true,1,1),
	 ('d006','Quality Management',true,1,1),
	 ('d007','Sales',true,1,1),
	 ('d008','Research',true,1,1),
	 ('d009','Customer Service',true,1,1),
	 ('d010','Information Technology',true,1,1);