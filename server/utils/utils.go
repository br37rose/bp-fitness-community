package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func SplitString(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func RemoveObjectIDFromArray(arr []primitive.ObjectID, element primitive.ObjectID) []primitive.ObjectID {
	var newArr []primitive.ObjectID
	for _, item := range arr {
		if item != element {
			newArr = append(newArr, item)
		}
	}
	return newArr
}

func RemoveElementFromArray(arr []string, element string) []string {
	var newArr []string
	for _, item := range arr {
		if item != element {
			newArr = append(newArr, item)
		}
	}
	return newArr
}
