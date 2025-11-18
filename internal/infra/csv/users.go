package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"processor/internal/domain"
)

type UsersCsvI interface {
	WriteAll(users []domain.User) error
}

type usersCsv struct {
	writer *csv.Writer
}

func NewUsersCsv(w io.Writer) usersCsv {
	return usersCsv{
		writer: csv.NewWriter(w),
	}
}

func (c *usersCsv) WriteAll(users []domain.User) error {
	records := make([][]string, 0, len(users)+1)

	for _, user := range users {
		records = append(records, []string{
			fmt.Sprintf("%d", user.Id),
			user.Name,
			user.Email,
			fmt.Sprintf("%d", user.Age),
		})
	}

	err := c.writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("write all users to csv: %w", err)
	}

	return nil
}
