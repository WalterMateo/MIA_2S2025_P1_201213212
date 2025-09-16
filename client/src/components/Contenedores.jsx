import axios from 'axios';
import React, {useRef, useState} from 'react';

export function Contenedores(){

    const [getSalidas, setSalidas] = useState([]);
    const [getTexto, setTexto] = useState("");
    //const [getarreglo, setarreglo] = useState([]);
    //no mio
    const [getArchivo, setArchivo] = useState(null);
    
    const fileInputRef = useRef(null);

    const [otroEstado, setOtroEstado] = useState('');


   
///------------------------------------------------------------------------------
const [textopantalla, setTextopantalla] = useState("Consola de Salida");



const cambiarPlaceholder = (nuevoTexto) => {
    setTextopantalla(nuevoTexto);
  };
///------------------------------------------------------------------------------

    //no mio
    /////////////////////////////////////////////////////////////////////
    const handleFileInputChange = (event) => {
        const Archivo = event.target.files[0];
        const Leedor = new FileReader();
    
        Leedor.onload = (e) => {
          // Guarda el contenido del archivo en la variable de estado
          setArchivo(e.target.result);
        };
    
        // Lee el contenido del archivo como texto
        Leedor.readAsText(Archivo);
      };

      const handleOpenFileClick = () => {
        fileInputRef.current.click();
      };

      const handleChange2 = (event) => {
        setTexto2(event.target.value);
      };

      const handleChangeOtroEstado = (event) => {
        setOtroEstado(event.target.value);
        // Realizar otras operaciones relacionadas con el cambio de otro estado si es necesario
      };
    /////////////////////////////////////////////////////////////////////
 
    function llamadaApi(){ 
        axios.get("http://localhost:4000/salida").then(
            (response) => {
            setSalidas(response.data.salidas);
            
            //console.log(response.data.message);
            }
        )

    }
   //////////////////////////////////////////
    function analizar(){

        if(getTexto=="") return;

        console.log(getTexto);  
        
        let body = {entrada : getTexto};
        axios
            .post("http://localhost:3000/interpretar", body)
            .then((response)=> {
               //setMensaje(response.data.message);
               console.log(response.data.message);
            });
                

    }
    /////////////////////////////////////////////////////////////////
   //let contenido = "echo 5+1+2+3;"
 function POST (ruta, contenido){
    console.log(contenido + "----------------");
    return fetch(ruta, {
        method: "POST",
        body: JSON.stringify({"lenguaje":contenido}),
        headers: {
            "Access-Control-Allow-Origin": "*",
            'Content-Type': 'application/json'
        },
        })
        .then(res => res.json())
        .catch(error => {
            console.error("Error en la solicitud:", error);
        });


}
                /// Tratando de conectar con el servidor
                
                function showValue() {
                    

                    
                    POST("http://localhost:3000/interpretar", getTexto ).then(res => {
                        console.log("222222222222222222222222222222222222222222222222222222222222222222222");
                        console.log(res.mensaje);
                        let textopantalla = res.mensaje;    
                        console.log("222222222222222222222222222222222222222222222222222222222222222222222");
                        console.log(textopantalla);
                        cambiarPlaceholder(textopantalla);
                        
                        //setResultado(res.mensaje);    
                    })
                }
    ///////////////////////////////////////////////////////////




    
        /*       
        getSalidas.map(elemento => {
            textopantalla += " > "+ elemento + "\n";
        });

        */

    return (<>
        <div class= "contenedor">
            
            
                <nav className="nav">
                    <ul className="Menu">
                        <li><a href="#"  class="bold" >Menu</a>
                            <ul>
                                <li><a href="#" class="bold" >Limpiar</a></li>
                                <li><a href="#" className="bold" onClick={(e) => { e.preventDefault(); fileInputRef.current.click(); }}>Abrir</a>
  <input
    type="file"
    onChange={handleFileInputChange}
    ref={fileInputRef}
    style={{ display: 'none' }}
  />
</li>

                                
                            </ul>
                        </li>
                        <li><a href="#"  class="bold"> Ejecutar</a>
                            <ul>
                                <li><a href="#"  class="bold" onClick={showValue}>Correr Programa</a></li>
                                
                            </ul>
                        </li>
                       

                    </ul>
                </nav>
            
            <header className="Heater">
                <h1> Proyecto 1 [GoDisk]</h1>
            </header>
            
            <div className="Widget-1">
                <h2>Entrada</h2>
            <textarea onChange={(e) => {setTexto(e.target.value)}} className="areaTexto" placeholder={getArchivo}  style={{ width: '850px', height: '500px', backgroundColor: '#002b36', color: '#839496', border: '2px solid #073642',  padding: '10px', resize: 'none' }}></textarea >
               
            
            </div>
            <div className="Widget-2">
                <h2>Salida</h2>

            <textarea className="areaTexto" placeholder={textopantalla} readOnly style={{width: '850px', height: '500px', backgroundColor: '#002b36', color: '#839496', border: '2px solid #073642',  padding: '10px', resize: 'none'  }}></textarea>

            </div>
            
            
        </div >
        </>);

}
//<textarea className="areaTexto" placeholder={textopantalla} readOnly style={{width: '850px', height: '500px', backgroundColor: '#002b36', color: '#839496', border: '2px solid #073642',  padding: '10px', resize: 'none'  }}></textarea>

