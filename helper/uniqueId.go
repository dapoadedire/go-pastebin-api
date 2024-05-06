package helper

import (
	"math/rand"
	"github.com/google/uuid"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func randomString(n int) string {  
  // Allocate space for the random string
  b := make([]byte, n)

  // Generate the first character as a random lowercase or uppercase letter
  b[0] = alphanumeric[rand.Intn(len(alphanumeric)-26)] 

  // Generate remaining characters (alphanumeric)
  for i := 1; i < n; i++ {
    b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
  }

  // Return the string conversion of the byte slice
  return string(b)
}

func GenerateUniqueID() string {
	return randomString(2) + uuid.New().String()[0:4]
}

