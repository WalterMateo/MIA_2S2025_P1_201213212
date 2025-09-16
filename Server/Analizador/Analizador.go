package Analizador

import (
	"CLASE4/FileSystem"
	"CLASE4/Gestion"
	"CLASE4/User"
	"bufio" // Para leer la entrada del usuario
	"flag"  // Para manejar parametros y opciones de comandos
	"fmt"   // para imprimir

	//"io/fs"
	"os"      // Para manejar archivos y errores ----> segun el aux para ingresar mediante consola
	"regexp"  //buscar y extraer parametros de la entrada
	"strings" // Para manipular cadenas de texto
)

// ER  mkdisk -size=3000 -unit=K -fit=BF -path=/home/cerezo/Disks/disk1.bin
var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input

	/* DEspues de procesar la entrada:
	comand sera "mkdisk"
	params sera "-size=3000 -unit=K -fit=BF -path=/home/cerezo/Disks/disk1.bin".*/

}

func Analizador() {

	for true {
		var input string
		fmt.Println(">>>>>>>>>>>>>>>>>>> INGRESE UN COMANDO  <<<<<<<<<<<<<<<<<<<<<<<<")
		fmt.Println("Ingrese Comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()

		command, params := getCommandAndParams(input)

		fmt.Println("comando: ", command, "  ||||||  ", "parametro: ", params)

		AnalyzeCommand(command, params)

	}
}

func AnalyzeCommand(command string, params string) {

	if strings.Contains(command, "mkdisk") {
		fn_mkdisk(params)
	} else if strings.Contains(command, "fdisk") {
		fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else if strings.Contains(command, "rmdisk") {
		fn_rmdisk(params)
	} else if strings.Contains(command, "login") {
		fn_login(params)
	} else if strings.Contains(command, "salir") {
		fmt.Println("Saliedo del programa... ")
		os.Exit(0) // Termina la ejecucion del programa
	} else {
		fmt.Println("Error: Commando invalido o no encontrado")
	}

}

func fn_mkdisk(params string) {

	//definir flag  que es para segmentar  mas todavia nuestro codigo
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño del disco")
	fit := fs.String("fit", "ff", "Algoritmo de ajuste")
	unit := fs.String("unit", "m", "Unidad del disco")
	path := fs.String("path", "", "Ruta del archivo de disco")

	//parseamos la flag porque
	fs.Parse(os.Args[1:])
	//-------------------------------------------extrae y asigna  los valores de los parametros
	matches := re.FindAllStringSubmatch(params, -1)

	// process the input

	for _, match := range matches {

		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //strings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas match[2]: Captura y guarda el valor del flag (por ejemplo, "3000", "k", "bf", "/home/cerezo/Disks/disk1.bin")

		flagValue = strings.Trim(flagValue, "\"") //Elimina las comillas de la entrada ejemplo "size" -> size

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: flag no encontrado")

		}
	}
	///-----------------------------------------------------
	/*
			Primera Iteración :
		    flagName es "size".
		    flagValue es "3000".
		    El switch encuentra que "size" es un flag reconocido, por lo que se ejecuta fs.Set("size", "3000").
		    Esto asigna el valor 3000 al flag size.

	*/
	//Validaciones>>>>>>>>>>>>>>>>>>>>>>>>>>>
	if *size <= 0 {
		fmt.Println("Error: El tamaño del disco debe ser mayor que cero.")
		return
	}
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf' o 'ff'.")
		return
	}
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: La unidad debe ser 'k' o 'm'.")
		return
	}
	if *path == "" {
		fmt.Println("Error: La ruta debe estar vacía.")
		return
	}
	//llamamos a la funcion
	Gestion.Mkdisk(*size, *fit, *unit, *path)
}

func fn_fdisk(input string) {

	//definir flags
	//flag.ExitOnError hace que el programa termine si hay un error al analizar las flags

	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño de la particion")
	path := fs.String("path", "", "Ruta del archivo de disco")
	name := fs.String("name", "", "Nombre de la particion")
	unit := fs.String("unit", "m", "Unidad de la particion")
	type_ := fs.String("type", "p", "Tipo de particion")
	fit := fs.String("fit", "", "Algoritmo de ajuste") //dejar el fit vacio por defecto

	//Parsear  los flags
	fs.Parse(os.Args[1:])

	//Encontrar los flags en el input
	matches := re.FindAllStringSubmatch(input, -1)

	//procesar el input
	for _, match := range matches {

		flagName := match[1]                   // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2]) //strings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas match[2]: Captura y guarda el valor del flag (por ejemplo, "3000", "k", "bf", "/home/cerezo/Disks/disk1.bin")

		flagValue = strings.Trim(flagValue, "\"") //Elimina las comillas de la entrada ejemplo "size" -> size

		switch flagName {
		case "size", "fit", "unit", "path", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: flag no encontrado")
		}
	}

	// Validaciones
	if *size <= 0 {
		fmt.Println("Error: El tamaño de la partición debe ser mayor que cero.")
		return
	}
	if *path == "" {
		fmt.Println("Error: La ruta no debe estar vacía.")
		return
	}

	//Si no se proporciono un fit, usar el valor predeterminado "w"
	if *fit == "" {
		*fit = "wf"
	}

	//Validar fit (b/w/f)
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: El ajuste debe ser 'bf', 'ff' o 'wf'.")
		return
	}
	if *unit != "k" && *unit != "m" && *unit != "b" {
		fmt.Println("Error: La unidad debe ser 'b', 'k' o 'm'.")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: El tipo debe ser 'p', 'e' o 'l'.")
		return
	}

	//llamamos a la funcion
	Gestion.Fdisk(*size, *path, *name, *unit, *type_, *fit)
}

func fn_mount(params string) {
	//definir flags
	fs := flag.NewFlagSet("mount", flag.ExitOnError)           //
	path := fs.String("path", "", "Ruta del archivo de disco") //
	name := fs.String("name", "", "Nombre de la particion")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]                      // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := strings.ToLower(match[2])    //strings.ToLower(match[2]): Captura y guarda el valor del flag, asegurándose de que esté en minúsculas match[2]: Captura y guarda el valor del flag (por ejemplo, "3000", "k", "bf", "/home/cerezo/Disks/disk1.bin")
		flagValue = strings.Trim(flagValue, "\"") //Elimina las comillas de la entrada ejemplo "size" -> size
		fs.Set(flagName, flagValue)
	}
	// Validaciones
	if *path == "" || *name == "" {
		fmt.Println("Error: Path y Name son obligatorios.")
		return
	}

	//convertir el nombre a minusculas antes de pasarlo al Mount
	lowercaseName := strings.ToLower(*name)
	Gestion.Mount(*path, lowercaseName)
}

func fn_mkfs(input string) {
	//definir flags
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "Identificador de la particion")
	type_ := fs.String("type", "full", "Tipo de formateo")
	fs_ := fs.String("fs", "2fs", "Fs") //-----------------> para el segundo proyecto

	// Parse the imput strings, not os.Args
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: flag no encontrado")
		}
	}

	// Validaciones
	if *id == "" {
		fmt.Println("Error: El id es un parametro obligatorio.")
		return
	}
	if *type_ == "" {
		fmt.Println("Error: El tipo de formateo es un parametro obligatorio.")
		return
	}

	//Llamar a la funcion
	FileSystem.Mkfs(*id, *type_, *fs_) //-----------------> para el segundo proyecto

}

func fn_rmdisk(params string) {
	//definir flags
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco a eliminar")
	// Parse the imput strings, not os.Args
	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: flag no encontrado")
		}
	}

	// Validaciones
	if *path == "" {
		fmt.Println("Error: El path es un parametro obligatorio.")
		return
	}

	// Llamamos a la funcion para eliminar el disco
	err := os.Remove(*path)
	if err != nil {
		fmt.Printf("Error: no se pude eliminar el disco en la ruta %s: %v\n", *path, err)
		return
	}
	fmt.Printf("Disco en la ruta %s eliminado correctamente. \n", *path)

}

func fn_login(input string) {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "ID de la partición")
	// Parse the imput strings, not os.Args
	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	//procesar el input
	for _, match := range matches {
		flagName := match[1]  // match[1]: Captura y guarda el nombre del flag (por ejemplo, "size", "unit", "fit", "path")
		flagValue := match[2] // match[2]: Captura y guarda el valor del flag (por ejemplo, "3000", "k", "bf", "/home/cerezo/Disks/disk1.bin")

		flagValue = strings.Trim(flagValue, "\"") //Elimina las comillas de la entrada ejemplo "size" -> size

		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: flag no encontrado")
		}
	}

	// Validaciones
	/*
		if *user == "" || *pass == "" || *id == "" {
			fmt.Println("Error: user, pass e id son parametros obligatorios.")
			return
		}
	*/

	User.Login(*user, *pass, *id)

}
