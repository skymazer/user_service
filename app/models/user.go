package models

type IdType uint64

type User struct {
	Id   IdType
	Name string
	Mail string
}
