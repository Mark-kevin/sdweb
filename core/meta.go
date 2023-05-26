package core

type SystemMeta struct {
	DiskStorage  map[string]string
	MemorySystem map[string]string
	CpuSystem    map[string]string
	CpuUsage     map[string]string
	LoadAverage  map[string]string
}

type FileMeta struct {
	Id    string
	Name  string
	Type  string
	Size  string
	IsDel string
}
