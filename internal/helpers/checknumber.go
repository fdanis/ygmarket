package helpers

func CheckNumber(cardnumber string) bool {
	len := len(cardnumber)
	number := 0
	sum := 0

	for i := 0; i < len; i++ {
		number = int(cardnumber[i] - '0')
		if i%2 == len%2 {
			number *= 2
			if number > 9 {
				number -= 9
			}
		}
		sum += number
		if sum >= 10 {
			sum -= 10
		}
	}
	return sum%10 == 0
}
