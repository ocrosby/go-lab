package util

func max(num1 int, num2 int) int {
	var result int

	if (num1 > num2) {
		result = num1
	} else {
		result = num2
	}
	
	return result
}
