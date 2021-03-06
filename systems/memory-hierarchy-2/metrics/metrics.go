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
	paymentAmounts []uint32
}

func AverageAge(users *UserData) float64 {
	totalAge := 0
	for _, age := range users.ages {
		totalAge += age
	}
	return float64(totalAge) / float64(len(users.ages))
}

func AveragePaymentAmount(users *UserData) float64 {
	var total1, total2 uint64
	pmts := users.paymentAmounts
	nPmts := len(pmts)
	i := 0
	for i < nPmts-1 {
		total1 += uint64(pmts[i])
		total2 += uint64(pmts[i+1])
		i += 2
	}

	for i < nPmts {
		total1 += uint64(pmts[i])
		i++
	}

	return float64(total1+total2) / 100 / float64(nPmts)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users *UserData) float64 {
	mean := AveragePaymentAmount(users) * 100
	squaredDiffs1, squaredDiffs2 := 0.0, 0.0
	pmts := users.paymentAmounts
	nPmts := len(pmts)

	i := 0
	for i < nPmts-1 {
		diff1 := float64(pmts[i]) - mean
		diff2 := float64(pmts[i+1]) - mean
		squaredDiffs1 += diff1 * diff1
		squaredDiffs2 += diff2 * diff2
		i += 2
	}

	for i < nPmts {
		diff1 := float64(pmts[i]) - mean
		squaredDiffs1 += diff1 * diff1
		i++
	}

	squaredDiffs := (squaredDiffs1 + squaredDiffs2) / 10000
	return math.Sqrt(squaredDiffs / float64(len(users.paymentAmounts)))
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
	paymentAmounts := make([]uint32, numPayments)
	for i, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		datetime, _ := time.Parse(time.RFC3339, line[1])
		paymentAmounts[i] = uint32(paymentCents)
		users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
			uint64(paymentCents),
			datetime,
		})
	}

	userData := UserData{ages: userAges, paymentAmounts: paymentAmounts}
	return users, userData
}
