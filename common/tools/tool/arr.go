package tool

func RemoveDuplicatesAndEmpty(arr []string) (ret []string){
	for i:=0; i < len(arr); i++{
		if (i > 0 && arr[i-1] == arr[i]) || len(arr[i])==0{
			continue
		}
		ret = append(ret, arr[i])
	}
	return
}