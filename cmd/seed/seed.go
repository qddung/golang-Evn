package seed

import (
	"fmt"

	"github.com/homework/lab/internal/infrastructure"
	"github.com/homework/lab/internal/test/data/fixture"
)

func main() {
	db := infrastructure.CreateSqlClient()
	u := fixture.NewUserTestCase(nil)
	u.SetDB(db)
	err := u.GenerateData()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate data: %v", err))
	}
}
