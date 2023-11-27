package utils

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
	//"time"

	"github.com/irononet/publisher/structs"
)



// LoadUserRecords loads user rows from a csv file 
// and returns an array of User objects
func LoadUserRecords(filepath string) ([]*structs.User, error){
	file, err := os.Open(filepath)  
	if err != nil{
		return nil, errors.New("could not find the filepath") 
	}

	// read the contents of the file 
	reader := csv.NewReader(file) 

	rows , err := reader.ReadAll() 
	if err != nil{
		return nil, errors.New("could not read csv records")
	}

	var records []*structs.User 

	for _, row := range rows{
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil{
			log.Println("error converting ID to int:", err) 
			continue
		}

		// createdAt, err := strconv.ParseInt(row[4], 10, 64)
		// if err != nil{
		// 	log.Println("error parsing createdAt field ", err)
		// 	continue 
		// }

		// deletedAt, err := strconv.ParseInt(row[5], 10, 64) 
		// if err != nil{
		// 	log.Println("error parsing deletedAt field: ", err)
		// 	continue 
		// }

		// mergedAt, err := strconv.ParseInt(row[6], 10, 64) 
		// if err != nil{
		// 	log.Println("error parsing mergedAt field", err)
		// 	continue 
		// }

		parentUserId, err := strconv.ParseInt(row[7], 10, 64) 
		if err != nil{
			log.Println("error parsing")
			continue 
		}

		user := &structs.User{
			ID: int(id), 
			FirstName: row[1], 
			LastName: row[2], 
			EmailAddress: row[3],
			CreatedAt: row[4],
			DeletedAt: row[5], 
			MergedAt: row[6],
			ParentUserId: int(parentUserId) ,
		}

		records = append(records, user)
	}

	return records, nil 
}