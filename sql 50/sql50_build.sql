create table student(
 sno varchar(10) primary key,
 sname varchar(20),
 sage int, 
 ssex varchar(5)
);

create table teacher(
 tno varchar(10) primary key, 
 tname varchar(20)
);

create table course(
 cno varchar(10),
 cname varchar(20), 
 tno varchar(20), 
 constraint pk_course primary key (cno,tno)
);

create table sc(
 sno varchar(10),
 cno varchar(10), 
 score DECIMAL(4,2),
 constraint pk_sc primary key (sno, cno)
);