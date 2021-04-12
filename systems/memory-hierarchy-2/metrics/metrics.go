package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type UserId int
type UserMap map[UserId]*User

type Address struct {
	fullAddress string
	zip         int
}

type Payment struct {
	amount uint64
	time   time.Time
}

type User struct {
	id       UserId
	name     string
	age      int
	address  Address
	payments []Payment
}

type UserData struct {
	ages           []int
	paymentAmounts []uint64
}

func AverageAge(users *UserData) float64 {
	totalAge := 0
	for _, age := range users.ages {
		totalAge += age
	}
	return float64(totalAge) / float64(len(users.ages))
}

func AveragePaymentAmount(users *UserData) float64 {
	average, count := 0.0, 0.0
	for _, amt := range users.paymentAmounts {
		count += 1
		average += (float64(amt) - average) / count
	}
	return average / 100
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users *UserData) float64 {
	mean := AveragePaymentAmount(users) * 100
	squaredDiffs := 0.0
	for _, amt := range users.paymentAmounts {
		diff := float64(amt) - mean
		squaredDiffs += diff * diff
	}
	return math.Sqrt(squaredDiffs / 10000 / float64(len(users.paymentAmounts)))
}

func LoadData() (UserMap, UserData) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	numUsers := len(userLines)
	users := make(UserMap, numUsers)
	userAges := make([]int, numUsers)
	for i, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		name := line[1]
		age, _ := strconv.Atoi(line[2])
		userAges[i] = age
		address := line[3]
		zip, _ := strconv.Atoi(line[3])
		users[UserId(id)] = &User{UserId(id), name, age, Address{address, zip}, []Payment{}}
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	numPayments := len(paymentLines)
	paymentAmounts := make([]uint64, numPayments)
	for i, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		datetime, _ := time.Parse(time.RFC3339, line[1])
		paymentAmounts[i] = uint64(paymentCents)
		users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
			uint64(paymentCents),
			datetime,
		})
	}

	userData := UserData{ages: userAges, paymentAmounts: paymentAmounts}
	return users, userData
}
