package main

import (
	"ddd/pkg/repository/mysql"
)

type ProjectConfig struct {
	Name       string `config:"name"`
	Port       string `config:"port"`
	AlarmEmail string `config:"alarm_email"`

	MySQLWriter mysql.Config `config:"repo.mysql.writer."`
	MySQLReader mysql.Config `config:"repo.mysql.reader."`

	MySQL struct {
		Writer mysql.Config `config:"repo.mysql.writer."`
		Reader mysql.Config `config:"repo.mysql.reader."`
	} `config:""`

	Repo struct {
		MySQL struct {
			Writer mysql.Config `config:"repo.mysql.writer."`
			Reader mysql.Config `config:"repo.mysql.reader."`
		} `config:""`
	}
}
