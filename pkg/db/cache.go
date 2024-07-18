package db

import cache2 "github.com/daodao97/goadmin/pkg/cache"

var cache cache2.Cache

func SetCache(c cache2.Cache) {
	cache = c
}
