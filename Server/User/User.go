package User

import (
	"fmt"
	"os"
	"strings"
	"Clase4/Gestion"
	"Clase4/Structs"
	"Clase4/Utilities"
	"encoding/binary"
)

func Login(user string, pass string, id string) {
	fmt.Println("===========INICIO LOGIN============")
	fmt.Println("USER:", user)
	fmt.Println("PASS:", pass)
	fmt.Println("ID:", id)

	//Obtener las partriciones montadas
	mountedPartitions := Gestion.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	login = false				
	
	//Verificar si el usuario ya esa logeado en alguna particion

	for _, partitions := range mountedPartitions {			// Recorre todas las particiones montadas
		for _, partition := range partitions {				// Recorre todas las particiones montadas
			if partition.ID == id && partition.LoggedIn {	// Si la particion ya tiene un usuario logeado
				fmt.Println("Ya existe un usuario logueado!")
				return
			}
			
	
	


}	