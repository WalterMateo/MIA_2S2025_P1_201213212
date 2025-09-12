package Structs

import (
	"fmt"
)

type MBR struct {
	MbrSize      int32        // 4 bytes // int32 va desde -2,147,483,648 hasta 2,147,483,647
	CreationDate [10]byte     // 10 bytes
	Signature    int32        // 4 bytes
	Fit          [1]byte      // 1 byte
	Partitions   [4]Partition // 4 particiones de 32 bytes cada una = 128 bytes
}

// Las particiones no se colocan aqui porque son otro tipo de estructuras

func PrintMBR(data MBR) {
	// Imprime los datos del MBR, formateando la fecha de creación, el ajuste (fit) y el tamaño total del disco
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.CreationDate[:]), string(data.Fit[:]), data.MbrSize))

	// Recorre las 4 particiones y las imprime
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}

}

type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       int32
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

func PrintPartition(data Partition) {
	fmt.Println(
		//Usa fmt.Sprintf para formatear la informacion de la particion en un solo string
		fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s",
			//string(data.Name[:]) convierte el array de bytes [16]bytes a una cadena de texto
			string(data.Name[:]), string(data.Type[:]), data.Start, data.Size,
			//data.Start y data.Size se imprimen como enteros (int32)
			string(data.Status[:]), string(data.Id[:])))
}

type EBR struct {
	PartMount byte
	PartFit   byte
	PartStart int32
	PartSize  int32
	PartNext  int32
	PartName  [16]byte
}

func PrintEBR(data EBR) {
	fmt.Println(fmt.Sprintf("Name: %s, fit: %c, start: %d, size: %d, next: %d, mount: %c",
		string(data.PartName[:]),
		data.PartFit,
		data.PartStart,
		data.PartSize,
		data.PartNext,
		data.PartMount))
}

type Superblock struct {
	S_filesystem_type   int32    // Guarda el número que identifica el sistema de archivos utilizado
	S_inodes_count      int32    // Cantidad total de inodos en el sistema de archivos
	S_blocks_count      int32    // Cantidad total de bloques en el sistema de archivos
	S_free_blocks_count int32    // Cantidad de bloques libres disponibles
	S_free_inodes_count int32    // Cantidad de inodos libres disponibles
	S_mtime             [17]byte // Fecha y hora de la última montura del sistema de archivos
	S_umtime            [17]byte // Fecha y hora de la última desmontura del sistema de archivos
	S_mnt_count         int32    // Número de veces que se ha montado el sistema de archivos
	S_magic             int32    // Número mágico que identifica el sistema de archivos
	S_inode_size        int32    // Tamaño del inodo (podria ser en bytes)=
	S_block_size        int32    // Tamaño del bloque (podria ser en bytes)=
	S_first_ino         int32    // Primer inodo libre (direccion del inodo)
	S_first_block       int32    // Primer bloque libre (direccion del bloque)
	S_bm_inode_start    int32    // Guardara el inicio del bitmap de inodos
	S_bm_block_start    int32    // Guardara el inicio del bitmap de bloques
	S_inode_start       int32    // Guardara el inicio de la tabla de inodos
	S_block_start       int32    // Guardara el inicio de la tabla de bloques
}

func PrintSuperblock(sb Superblock) {
	fmt.Println("====== Superblock ======")
	fmt.Printf("S_filesystem_type: %d\n", sb.S_filesystem_type)
	fmt.Printf("S_inodes_count: %d\n", sb.S_inodes_count)
	fmt.Printf("S_blocks_count: %d\n", sb.S_blocks_count)
	fmt.Printf("S_free_blocks_count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("S_free_inodes_count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("S_mtime: %s\n", string(sb.S_mtime[:]))
	fmt.Printf("S_umtime: %s\n", string(sb.S_umtime[:]))
	fmt.Printf("S_mnt_count: %d\n", sb.S_mnt_count)
	fmt.Printf("S_magic: 0x%X\n", sb.S_magic) // Usamos 0x%X para mostrarlo en formato hexadecimal
	fmt.Printf("S_inode_size: %d\n", sb.S_inode_size)
	fmt.Printf("S_block_size: %d\n", sb.S_block_size)
	fmt.Printf("S_fist_ino: %d\n", sb.S_first_ino)
	fmt.Printf("S_first_blo: %d\n", sb.S_first_block)
	fmt.Printf("S_bm_inode_start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("S_bm_block_start: %d\n", sb.S_bm_block_start)
	fmt.Printf("S_inode_start: %d\n", sb.S_inode_start)
	fmt.Printf("S_block_start: %d\n", sb.S_block_start)
	fmt.Println("========================")
}

type Inode struct {
	I_uid   int32     // Identificador del usuario propietario del archivo o directorio
	I_gid   int32     // Identificador del grupo propietario del archivo o directorio
	I_size  int32     // Tamaño del archivo o directorio en bytes
	I_atime [17]byte  // Fecha y hora del último acceso al archivo o directorio
	I_ctime [17]byte  // Fecha y hora de la creación del archivo o directorio
	I_mtime [17]byte  // Fecha y hora de la última modificación del archivo o directorio
	I_block [15]int32 // Punteros a los bloques de datos (12 directos, 1 indirecto simple, 1 indirecto doble, 1 indirecto triple)
	I_type  [1]byte   // Tipo de archivo ('0' para archivo, '1' para directorio)
	I_perm  [3]byte   // Permisos de acceso al archivo o directorio (lectura, escritura, ejecución)
}

func PrintInode(inode Inode) {
	fmt.Println("====== Inode ======")
	fmt.Printf("I_uid: %d\n", inode.I_uid)
	fmt.Printf("I_gid: %d\n", inode.I_gid)
	fmt.Printf("I_size: %d\n", inode.I_size)
	fmt.Printf("I_atime: %s\n", string(inode.I_atime[:]))
	fmt.Printf("I_ctime: %s\n", string(inode.I_ctime[:]))
	fmt.Printf("I_mtime: %s\n", string(inode.I_mtime[:]))
	fmt.Printf("I_type: %s\n", string(inode.I_type[:]))
	fmt.Printf("I_perm: %s\n", string(inode.I_perm[:]))
	fmt.Printf("I_block: %v\n", inode.I_block)
	fmt.Println("===================")
}

type Folderblock struct {
	B_content [4]Content // Cada bloque de carpeta puede contener hasta 4 entradas (archivos o subcarpetas)
}

func PrintFolderblock(folderblock Folderblock) {
	fmt.Println("====== Folder Block ======")
	for i, content := range folderblock.B_content {
		fmt.Printf("Content %d: Name: %s, Inodo: %d\n", i, string(content.B_name[:]), content.B_inodo)
	}
	fmt.Println("========================")
}

type Content struct {
	B_name  [12]byte // Este campo almacena el nombre de un archivo o carpeta
	B_inodo int32    //Se utiliza un arreglo de bytes ([12]byte) para asegurar que el nombre tenga un tamaño fijo y sea fácil de manejar
}

/*
func PrintContent(content Content) {
	fmt.Println("====== Content ======")
	fmt.Printf("B_name: %s\n", string(content.B_name[:]))
	fmt.Printf("B_inodo: %d\n", content.B_inodo)
	fmt.Println("=====================")
}
*/

type Fileblock struct {
	B_content [64]byte // Cada bloque de archivo puede almacenar hasta 64 bytes de datos
}

func PrintFileblock(fileblock Fileblock) {
	fmt.Println("====== File Block ======")
	fmt.Printf("B_content: %s\n", string(fileblock.B_content[:]))
	fmt.Println("========================")
}

type PointerBlock struct {
	B_pointers [16]int32 //Este campo es un apuntador al inodo asociado al archivo o carpeta.
	//Cada apuntador es un int32 (4 bytes), lo que permite direccionar hasta 2^32 bloques
}

func PrintPointerBlock(pointerblock PointerBlock) {
	fmt.Println("====== Pointer Block ======")
	for i, pointer := range pointerblock.B_pointers {
		fmt.Printf("Pointer %d: %d\n", i, pointer)
	}
	fmt.Println("===========================")
}
