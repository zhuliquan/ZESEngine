-- 创建一个搜索引擎数据库
create database if not exists zesengine_db default character set utf8;
use zesengine_db;

-- 创建关键字字典
create table if not exists key_word (
    wordId int primary key auto_increment,
    word varchar(1000) unique,
    avaliable bit(1) default 1
);

-- 创建原始的文档字典
create table if not exists real_doc (
    realDocId int primary key auto_increment,
    filePath varchar(1000) unique,
    avaliable bit(1) default 1
);

-- 创建一个正排索引的表
create table if not exists forward_index (
    docId int primary key auto_increment,
    content varchar(10000),
    avaliable bit(1) default 1,
    realDocId int,
    foreign key(realDocId) references real_doc(realDocId)
);

-- 创建 wordId -- DocId 的联系表
create table if not exists post_item (
    linkId int primary key auto_increment,
    docId int,
    wordId int, 
    termFrequence int default 0, -- 词频
    wordWeight float default 0.0, -- 单词权重
    avaliable bit(1) default 1,
    foreign key(wordId) references key_word(wordId),
    foreign key(docId) references forward_idx(docId)
);