package User

import (
	"CLASE4/Gestion"
	"CLASE4/Structs"
	"CLASE4/Utilities"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
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
	var login bool = false

	//Verificar si el usuario ya esa logeado en alguna particion

	for _, partitions := range mountedPartitions { // Recorre todas las particiones montadas
		for _, partition := range partitions { // Recorre todas las particiones montadas
			if partition.ID == id && partition.LoggedIn { // Si la particion ya tiene un usuario logeado
				fmt.Println("Ya existe un usuario logueado!")
				return
			}
			if partition.ID == id { // Encuentra la particion correcta
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}
	// Si no se encontró la partición montada, se detiene el proceso
	if !partitionFound {
		fmt.Println("Error: No se encontro ninguna particion montada con el ID proporcionado.")
		return
	}

	//Abrir el archivo del sistema de archivos binario
	file, err := Utilities.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close() //cierra el archivo al final de la ejecucion

	var TempMBR Structs.MBR
	// Leer el MBR (Master Boot Record) del archivo binario
	if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	//Imprimir la informacion del MBR
	Structs.PrintMBR(TempMBR)
	fmt.Println("-----------------------------------")

	var index int = -1
	// Buscar la particion en el MBR por su ID

	for i := 0; i < 4; i++ {
		if TempMBR.Partitions[i].Size != 0 { // Verifica que la particion tenga un tamaño mayor a 0
			if strings.Contains(string(TempMBR.Partitions[i].Id[:]), id) { // Compara el ID con el nombre de la particion
				fmt.Println("Particion encontrada en el MBR")
				if TempMBR.Partitions[i].Status == '1' { // Verifica que la particion esta montada
					fmt.Println("Particion montada")
					index = i
				} else {
					fmt.Println("Error: La particion no esta montada")
					return
				}
				break
			}
		}
	}

	//Si se encontro la particion, imprimir su informacion
	if index != -1 {
		Structs.PrintPartition(TempMBR.Partitions[index])
	} else {
		fmt.Println("Error: No se encontro la particion en el MBR")
		return
	}

	var tempSuperblock Structs.Superblock // Variable para almacenar el superbloque
	// Leer el superbloque de la particion
	if err := Utilities.ReadObject(file, &tempSuperblock, int64(TempMBR.Partitions[index].Start)); err != nil { // Lee el superbloque desde el inicio de la particion con el ID proporcionado ademas del tamaño del MBR desde el inicio
		fmt.Println("Error: No se pudo leer el Superbloque:", err)
		return
	}

	//Buscar el archivo de usuarios "/users.txt" en el sistema de archivos
	indexInode := InitSearch("/user.txt", file, tempSuperblock) // Obtener el índice del inodo del archivo de usuarios

	var crrInode Structs.Inode // Variable para almacenar el inodo actual
	// Leer el inodo del archivo ""user.txt""
	if err := Utilities.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(Structs.Inode{})))); err != nil { // Lee el inodo desde el inicio de la tabla de inodos
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	//Obtener el contenido del archivo users.txt desde los bloques del inodo
	data := GetInodeFileData(crrInode, file, tempSuperblock) // Obtener los datos del archivo desde los bloques del inodo

	//Dividir el contenido del archivo en lineas
	lines := strings.Split(data, "\n") // Dividir el contenido en líneas

	//Iterar a traves de las lineas para verificar las credenciales
	for _, line := range lines {
		words := strings.Split(line, ",") // Dividir la línea en palabras

		//Si la linea tiene 5 elementos, verificar si el usuario y contraseña coinciden
		if len(words) == 5 { // Verifica que la línea tenga 5 elementos (ID, Tipo, Grupo, Usuario, Contraseña)
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) { // Verifica si el usuario y la contraseña coinciden con los datos del archivo
				login = true
				break
			}
		}
	}

	//Imprimir la informacion del inodo
	fmt.Println("Inode", crrInode.I_block)

	//Si las credenciales son correctas, marcar la particion como logeada
	if login {
		fmt.Println("Usuario Logeado con exito")
		Gestion.MarkPartitionAsLoggedIn(id) // Marcar la partición como logeada
	}
	fmt.Println("===========FIN LOGIN============")
}

func InitSearch(path string, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("==========START BUSQUEDA INICIAL============")
	fmt.Println("path:", path)

	//Dividir la ruta en carpetas usando "/" como separador
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:] // Omitir el primer elemento vacío si la ruta empieza con "/"

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath { // Iterar sobre cada paso en la ruta
		fmt.Println("Step:", step) // Imprimir el paso actual
	}

	var Inode0 Structs.Inode

	//Leer el inodo raiz (primer inodo del sistema de archivos)
	if err := Utilities.ReadObject(file, &Inode0, int64(tempSuperblock.S_inode_start)); err != nil {
		return -1 // Si hay un error al leer el inodo, retornar -1
	}

	fmt.Println("===========END BUSQUEDA INICIAL============")

	//Llamar a la funcion que busca el inodo del archivo segun la ruta
	return SearchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// Stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1 // Obtener el índice del último elemento
	last := (*s)[lastIndex]  // Obtener el último elemento
	*s = (*s)[:lastIndex]    // Remover el último elemento del slice
	return last              // Retornar el último elemento
}

func SearchInodeByPath(StepsPath []string, Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) int32 {
	fmt.Println("===========START Busqueda Inodo Por Path============")

	index := int32(0) // Contador de bloques procesados en el inodo actual

	//Extrae el primer elemento del path y elimina los espacios en blanco
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("=========================SearchedName: ", SearchedName)

	// Iterar sobre los bloques del inodo
	for _, block := range Inode.I_block {
		if block != -1 { // Si el bloque es válido (no está vacío)
			if index < 13 { // Manejo de bloques directos (0-12)
				var crrFolderBlock Structs.Folderblock

				// Leer el bloque de carpeta desde el archivo binario
				if err := Utilities.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Folderblock{})))); err != nil {
					return -1
				}

				// Buscar el archivo/directorio dentro del bloque de carpeta
				for _, folder := range crrFolderBlock.B_content {
					fmt.Println("Folder === Name:", string(folder.B_name[:]), "B_inodo", folder.B_inodo)

					// Si el nombre del archivo o directorio coincide
					if strings.Contains(string(folder.B_name[:]), SearchedName) {
						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)

						if len(StepsPath) == 0 { // Si llegamos al final de la ruta
							fmt.Println("Folder found======")
							return folder.B_inodo // Retornar índice del inodo encontrado
						} else { // Si aún hay más niveles en la ruta, seguir buscando
							fmt.Println("NextInode======")
							var NextInode Structs.Inode

							// Leer el siguiente inodo desde el archivo binario
							if err := Utilities.ReadObject(file, &NextInode, int64(tempSuperblock.S_inode_start+folder.B_inodo*int32(binary.Size(Structs.Inode{})))); err != nil {
								return -1
							}

							// Llamada recursiva para seguir con la búsqueda
							return SearchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}
			} else {
				fmt.Println("Manejo de bloques indirectos no implementado") // Falta implementar acceso a bloques indirectos
			}
		}
		index++ // Incrementar índice para saber si son bloques directos o indirectos
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0 // Si no se encontró, retornar 0
}

func GetInodeFileData(Inode Structs.Inode, file *os.File, tempSuperblock Structs.Superblock) string {
	fmt.Println("===========START CONTENIDO DEL BLOQUE============")
	index := int32(0) // Contador de bloques procesados en el inodo
	//Define content as a string
	var content string

	// Iterar sobre los bloques del inodo
	for _, block := range Inode.I_block {
		if block != -1 { // Si el bloque es válido (no está vacío)
			//Dentro de los primeros 12 bloques, son bloques directos
			if index < 13 {
				var crrFileBlock Structs.Fileblock
				// Leer el bloque de archivo desde el archivo binario
				if err := Utilities.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_block_start+block*int32(binary.Size(Structs.Fileblock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_content[:]) // Concatenar el contenido del bloque al contenido total

			} else {
				fmt.Println("indirectos") // Falta implementar acceso a bloques indirectos
			}
		}
		index++ // Incrementar índice para saber si son bloques directos o indirectos
	}
	fmt.Println("===========END CONTENIDO DEL BLOQUE============")
	return content // Retornar el contenido completo del archivo
}

// MKUSER
func AppendToFileBlock(inode *Structs.Inode, newData string, file *os.File, superblock Structs.Superblock) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.I_block)*binary.Size(Structs.Fileblock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock Structs.Fileblock
	copy(updatedFileBlock.B_content[:], fullData)
	if err := Utilities.WriteObject(file, updatedFileBlock, int64(superblock.S_block_start+inode.I_block[0]*int32(binary.Size(Structs.Fileblock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_size = int32(len(fullData))
	if err := Utilities.WriteObject(file, *inode, int64(superblock.S_inode_start+inode.I_block[0]*int32(binary.Size(Structs.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}
