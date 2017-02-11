CREATE TABLE user_info (
    uname varchar(64) NOT NULL,
    passwd varchar(64) NOT NULL,
    time_created int(11) NOT NULL,
    PRIMARY KEY (uname)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
