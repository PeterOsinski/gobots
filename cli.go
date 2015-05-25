package main

import "flag"

type ConfigValues struct {
	workerNum 	int
	file			string
}

func parseCli() *ConfigValues{
	
	cf := new(ConfigValues)
	
	workerNum := flag.Int("worker_num", 10, "number of pararrel requests")
	file := flag.String("file", "./list", "Path to file with list of hostnames")
	
	flag.Parse()
	
	cf.workerNum = *workerNum
	cf.file = *file
	
	Logger.Printf("Config: %+v\n", cf)
	
	return cf
}