package main

func int64Ptr(i int64) *int64 { return &i }

func PanicIfErr(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}
