package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomDinoName() string {
	dinos := [11]string{"Abelisaurus", "Euoplocephalus", "Liliensternus", "Troodon", "Ultrasauros", "Coelophysis", "Gallimimus", "Halticosaurus", "Daemonosaurus", "Dilophosaurus", "Magyarosaurus"}
	// to avoid to get every time the same number
	rand.Seed(time.Now().UnixNano())

	// Shuffle dinos array
	rand.Shuffle(len(dinos), func(i, j int) { dinos[i], dinos[j] = dinos[j], dinos[i] })

	return dinos[0]
}
