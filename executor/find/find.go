package find

type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// val(i) less than val return -1;  val(i) equal val return 0; val(i) greet than val return 1;
	Compare(i int, val interface{}) int
}

func Find(data Interface, input interface{}) int {
	n := data.Len()
	return binaryFind(data, 0, n-1, input)
}

func binaryFind(data Interface, leftIndex int, rightIndex int, findVal interface{}) int {
	//判断 leftIndex 是否大于 rightIndex
	if leftIndex > rightIndex {
		return -1
	}
	//先找到 中间的下标
	middle := (leftIndex + rightIndex) / 2
	if data.Compare(middle, findVal) > 0 { //(*arr)[middle] > findVal {
		//说明我们要查找的数，应该在
		return binaryFind(data, leftIndex, middle-1, findVal)
	} else if data.Compare(middle, findVal) < 0 { // (*arr)[middle] < findVal {
		//说明我们要查找的数，应该在 middel+1 --- rightIndex
		return binaryFind(data, middle+1, rightIndex, findVal)
	} else {
		return middle
	}
}
