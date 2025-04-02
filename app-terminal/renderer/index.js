const {createApp, ref, onMounted, watch, computed, onUnmounted} = Vue;
import {printTemplate, alertaError, alertaSuccess, alertaWarning} from './utils.js';

const app = createApp({
    setup() {
        const empresa = ref({
            Nauto: "",//NO SE USA
            Nombre: "",
            Resolucion: "",
            Documento: "",
            Permiso: "",
        })
        // Datos del vehículo
        const vehiculo = ref({
            Patente: "",
            Tipo: "",
            Modelo: "",
            Motor: "",
            Marca: "",
            Anio: "",
            Pais: "",
            Empresa: "",
            Chasis: "",
            NroAuto: "",
            Costo: "",
        });

        const buscarVehiculoPorPatente = async () => {
            if (!vehiculo.value.Patente || vehiculo.value.Patente == "") {
                return "PATENTE VACIA";
            }
            try {
                const resultado = await window.electronAPI.obtenerVehiculo(
                    vehiculo.value.Patente
                ); // Llama a la función del repositorio
                if (!resultado.error || resultado.error === "") {
                    // Actualiza las propiedades del objeto vehiculo
                    vehiculo.value = {
                        ...vehiculo.value, // Mantén las propiedades existentes (por si acaso)
                        ...resultado, // Sobrescribe con los nuevos datos
                    };

                    return ""
                }

                return "No se encontró el vehiculo";
            } catch (error) {
                await alertaError(`Error al obtener el vehículo con la patente ${vehiculo.value.Patente}`,
                    error);
                return error
            }
        };

        // Datos del conductor
        const conductor = ref({
            Documento: "",
            Licencia: "",
            Nombre: "",
            ApellidoPa: "",
            ApellidoMa: "",
            Direccion: "",
            Nacionalidad: "",
            Residencia: "",
            Profesion: "",
            FechaNac: "",
            EstCivil: "",
            Sexo: "",
            Mva1: "",
            Mva2: "",
            Tipoa: "",
            Nauto: "",
        });

        const buscarConductorPorDocumento = async () => {
            try {
                const resultado = await window.electronAPI.obtenerConductor(
                    conductor.value.Documento
                ); // Llama a la función del repositorio
                if (!resultado.error || resultado.error === "") {
                    // Actualiza las propiedades del objeto vehiculo
                    conductor.value = {
                        ...resultado.data, // Sobrescribe con los nuevos datos
                    };

                    return ""
                }
            } catch (error) {
                return error
            }
        };

        const buscarVehiculoConConductorPorNroAuto = async () => {
            try {
                const resultado = await window.electronAPI.obtenerVehiculoConductor(
                    vehiculo.value.NroAuto
                );


                if (!resultado.error || resultado.error === "") {
                    let empresaData = {
                        Nombre: resultado.data.empresa.nombre,
                        Resolucion: resultado.data.empresa.resolucion,
                        Documento: resultado.data.empresa.documento,
                        Permiso: resultado.data.empresa.permiso,
                    }

                    vehiculo.value = {...vehiculo.value, ...resultado.data.vehiculo};
                    conductor.value = {...conductor.value, ...resultado.data.conductor};
                    empresa.value = {...empresa.value, ...empresaData};
                    return ""
                }
            } catch (error) {
                console.error(
                    `Error al obtener el vehículo con conductor para el nroAuto ${vehiculo.value.NroAuto}`,
                    error
                );
                return error
            }
        };

        // Datos del pasajero
        const pasajero = ref({
            Documento: "",
            Licencia: "",
            Nombre: "",
            ApellidoPa: "",
            ApellidoMa: "",
            Direccion: "",
            Nacionalidad: "",
            Residencia: "",
            Profesion: "",
            FechaNac: "",
            EstCivil: "",
            Sexo: "",
            Mva1: "",
            Mva2: "",
            Tipoa: "",
            Nauto: "",
            TipoDocumento: "",
        });

        const buscarPasajeroPorDocumento = async () => {
            try {
                const resultado = await window.electronAPI.obtenerConductor(
                    pasajero.value.Documento
                ); // Llama a la función del repositorio
                if (!resultado.error || resultado.error === "") {
                    // Actualiza las propiedades del objeto vehiculo
                    pasajero.value = {
                        //...pasajero.value, // Mantén las propiedades existentes (por si acaso)
                        ...resultado.data, // Sobrescribe con los nuevos datos
                    };

                    await agregarDato();

                    return ""
                }
            } catch (error) {
                console.error(
                    `Error al obtener el pasajero con el documento ${pasajero.value.Documento}`,
                    error
                );
                return error
            }
        };

        const listarUltimaRelacion = async () => {
            try {
                // Llamada a la función que obtiene la última relación
                const resultado = await window.electronAPI.ultimaRelacion();
                // Verifica si el resultado existe y tiene IdRelacion
                if (resultado && resultado.IdRelacion) {
                    // Convierte IdRelacion a número y suma 1
                    const idRelacion = Number(resultado.IdRelacion);

                    // Asegúrate de que IdRelacion sea un número válido
                    if (!isNaN(idRelacion)) {
                        relacionInfo.value.NroRelacion = String(idRelacion + 1);
                    } else {
                        console.error(
                            "IdRelacion no es un número válido:",
                            resultado.IdRelacion
                        );
                        relacionInfo.value.NroRelacion = "1"; // Valor predeterminado si IdRelacion no es válido
                    }
                } else {
                    // Si no hay resultado o IdRelacion es undefined/null, asigna un valor predeterminado
                    relacionInfo.value.NroRelacion = "1";
                }
            } catch (error) {
                console.error("Error al obtener la última relación:", error);
            }
        };

        // Datos falsos para el listado
        const listado = ref([]);

        // Información adicional
        const relacionInfo = ref({
            IsFecha : true,
            Procedencia: true,
            Destino:true,
            FechaSalida: new Date().toLocaleDateString("es-ES", {
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
            }),
            FechaRetorno: new Date().toLocaleDateString("es-ES", {
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
            }),
            NroRelacion: "",
            Fecha: new Date().toLocaleDateString("es-ES", {
                day: "2-digit",
                month: "2-digit",
                year: "numeric",
            }),
            NroPasajero: computed(() => listado.value.length || 0),
        });

        // Función para agregar un nuevo registro
        const agregarDato = async () => {
            listado.value.push({
                Nro: (listado.value.length + 1).toString().padStart(2, "0"), // Temporalmente asigna el número
                Documento: pasajero.value.Documento,
                TipoDocumento: pasajero.value.TipoDocumento || "",
                Licencia: pasajero.value.Licencia || "",
                Nauto: pasajero.value.Nauto || "",
                Residencia: pasajero.value.Residencia || "",
                Nombres: pasajero.value.Nombre || "",
                ApellidoP: pasajero.value.ApellidoPa || "",
                ApellidoM: pasajero.value.ApellidoMa || "",
                Direccion: pasajero.value.Direccion || "",
                Nacionalidad: pasajero.value.Nacionalidad || "",
                Profesion: pasajero.value.Profesion || "",
                FechaNac: pasajero.value.FechaNac || "",
                EstCivil: pasajero.value.EstCivil || "",
                Sexo: pasajero.value.Sexo || "",
                Mva1: pasajero.value.Mva1 || "",
                TarPeru: "✔",
                TarChile: "✔",
            });

            // Limpia el formulario
            pasajero.value = {
                Documento: "",
                TipoDocumento: "",
                Licencia: "",
                Nauto: "",
                Residencia: "",
                Nombre: "",
                ApellidoPa: "",
                ApellidoMa: "",
                Direccion: "",
                Nacionalidad: "",
                Profesion: "",
                FechaNac: "",
                EstCivil: "",
                Sexo: "",
                Mva1: "",
            };
        };

        //COSAS DEL MODAL
        const showModal = ref(false);
        const showModalTarjetas = ref(false);
        const tasks = ref([
            {text: "Creando o actualizando pasajeros...", done: false, details: ""},
            {text: "Creando relación...", done: false, details: ""},
            {
                text: "Agregando detalles de la relación...",
                done: false,
                details: "",
            },
        ]);

        const crearOActualizarPasajeros = async () => {
            for (const pasajeroItem of listado.value) {
                let conductor = {};
                conductor.Documento = pasajeroItem.Documento;
                conductor.TipoDocumento = pasajeroItem.TipoDocumento;
                conductor.Licencia = pasajeroItem.Licencia.trim() !== "" ? pasajeroItem.Licencia : "";
                conductor.Nombre = pasajeroItem.Nombres;
                conductor.ApellidoPa = pasajeroItem.ApellidoP;
                conductor.ApellidoMa = pasajeroItem.ApellidoM;
                conductor.Direccion = pasajeroItem.Direccion;
                conductor.Nacionalidad = pasajeroItem.Nacionalidad;
                conductor.Residencia = pasajeroItem.Residencia;
                conductor.Profesion = pasajeroItem.Profesion;
                conductor.FechaNac = pasajeroItem.FechaNac;
                conductor.EstCivil = pasajeroItem.EstCivil;
                conductor.Sexo = pasajeroItem.Sexo;
                conductor.Mva1 = pasajeroItem.Mva1; // TODO:
                conductor.Mva2 = pasajeroItem.Mva1;
                conductor.Tipoa = pasajeroItem.Mva1;
                conductor.Nauto = pasajeroItem.Nauto.trim() !== "" ? pasajeroItem.Nauto : "";

                if (!conductor.Documento || conductor.Documento.trim() === "") {
                    showError(`DOCUMENTO vacio - pasajero no actualizado`)
                    continue;
                }

                let response = await window.electronAPI.createUpdateConductor(
                    conductor
                );
                if (response.error !== "") {
                    console.error(response.error);
                } else {
                    tasks.value[0].details += `✓ ${pasajeroItem.Nombres} ${pasajeroItem.ApellidoP}\n`;
                }
            }
        };
        const crearRelacion = async () => {
            let relacion = {};
            relacion.IdRelacion = relacionInfo.value.NroRelacion;
            relacion.FkPatente = vehiculo.value.Patente;
            relacion.FkDocumento = conductor.value.Documento;
            relacion.Fesal = relacionInfo.value.FechaSalida;
            relacion.Fere = relacionInfo.value.FechaRetorno;
            //relacion.Prot = relacionInfo.value.Procedencia ? "TACNA" : "";//TODO cambiar
            //relacion.Desa = relacionInfo.value.Destino ? "ARICA" : "";
            relacion.Prot = "TACNA"
            relacion.Desa=  "ARICA"
            relacion.Proa = "ARICA";
            relacion.Dest = "TACNA";

            let response = await window.electronAPI.agregarRelacion(relacion);
            if (response.error !== "") {
                console.log(response);
                console.log(
                    `Error al crear la relación con ID ${relacionInfo.value.NroRelacion}`
                );
                throw new Error(`Error al crear la relación`);
            } else {
                tasks.value[1].details = `✓ Relación creada con ID ${relacionInfo.value.NroRelacion}`;
            }
        };
        const crearDetalleRelacion = async () => {
            for (const pasajeroItem of listado.value) {
                let detalleRelacion = {};
                detalleRelacion.FkRelacion = relacionInfo.value.NroRelacion;
                detalleRelacion.FkDocumento = pasajeroItem.Documento;
                detalleRelacion.FkNombre = pasajeroItem.Nombres;
                detalleRelacion.FkApellidoPa = pasajeroItem.ApellidoP;
                detalleRelacion.FkApellidoMa = pasajeroItem.ApellidoM;
                detalleRelacion.FkDireccion = pasajeroItem.Direccion;
                detalleRelacion.FkNacionalidad = pasajeroItem.Nacionalidad;
                detalleRelacion.FkResidencia = pasajeroItem.Residencia;
                detalleRelacion.FkProfesion = pasajeroItem.Profesion;
                detalleRelacion.FkFechanac = pasajeroItem.FechaNac;
                detalleRelacion.FkEstcivil = pasajeroItem.EstCivil;
                detalleRelacion.FkSexo = pasajeroItem.Sexo;
                detalleRelacion.FkMva1 = pasajeroItem.Mva1; //TODO
                detalleRelacion.FkMva2 = pasajeroItem.Mva1;
                detalleRelacion.FkTipoa = pasajeroItem.Mva1;

                if (!pasajeroItem.Documento || pasajeroItem.Documento.trim() === "") {
                    showError(`DOCUMENTO vacio - pasajero no actualizado`)
                    continue;
                }

                let response = await window.electronAPI.agregarDetalleRelacion(
                    detalleRelacion
                );
                if (response.error !== "") {
                    console.log(response);
                    console.log(
                        `Error al intentar agregar el detalle relación con el documento ${pasajeroItem.Documento}`
                    );
                    throw new Error(
                        `Error al procesar el detalle de relación con documento ${pasajeroItem.Documento}`
                    );
                } else {
                    tasks.value[2].details = "✓ Detalles agregados correctamente";
                }
            }
        };
        //PROCESO DE IMPRESION

        const selectOption = async (option) => {
            await startProcess(option); // Llama al proceso con la opción seleccionada
            //closeModal(); // Cierra el modal después de seleccionar //TODO
        };

        const getProcedenciaDestino = async (option) => {
            const normalizedOption = option.toUpperCase();
            if (normalizedOption.includes("T-A")) {
                return {pro: "TACNA", dest: "ARICA"};
            }

            if (normalizedOption.includes("A-T")) {
                return {pro: "ARICA", dest: "TACNA"};
            }

            return {pro: "", dest: ""};
        }

        const procedencia = ref("");
        const destino = ref("");
        const guardarHistorial = async (id) => {
            window.electronAPI.guardarHistorial(id)
        }
        const startProcess = async (option) => {
            if (!option.includes("CABEZERA")) {
                // Verifica si hay detalles de los pasajeros
                if (!listado.value || listado.value.length === 0) {
                    await alertaWarning(
                        "No se puede imprimir porque no hay pasajeros"
                    );
                    return; // Detiene el proceso si no hay detalles
                }

                if (!listado.value.some(item => item.Documento && item.Documento.trim() !== "")){
                    await alertaWarning(
                        "No se puede imprimir porque no hay pasajeros validos"
                    );
                    return; // Detiene el proceso si no hay detalles
                }
            }

            const data = await window.electronAPI.getPreferences();
            if (data[option] === undefined) {
                await alertaWarning("Necesita configurar una plantilla para esta opcion")
                return
            }

            let impresoraActual = data["IMPRESORA-ACTUAL"]
            if(!impresoraActual || impresoraActual.trim() === "") {
                await alertaWarning("Necesita configurar una impresora para esta opcion")
                return
            }

            let plantillaResponse = await window.electronAPI.buscarPlantilla(
                data[option]
            );
            if (plantillaResponse.error !== "") {
                await alertaError("Error al buscar la plantilla.");
                return;
            }
            const {pro, dest} = await getProcedenciaDestino(option)
            procedencia.value = relacionInfo.value.Procedencia ? pro : "";
            destino.value = relacionInfo.value.Destino ? dest : "";
            //Si todo está bien, inicia el proceso
            showModal.value = true;
            try {
                await crearRelacion();
                await crearOActualizarPasajeros();
                await crearDetalleRelacion();
                await guardarHistorial(relacionInfo.value.NroRelacion)
                showSuccessPrint("TODO LISTO!, IMPRIMIENDO")
                await listarUltimaRelacion()
                let dataParaImprimir = await processData(option)
                await printTemplate(plantillaResponse.data.datos, dataParaImprimir, plantillaResponse.data.fuente, impresoraActual)
                await showSuccessPrint(`NUEVO NUMERO DE RELACION : ${relacionInfo.value.NroRelacion}`)
            } catch (error) {
                showError(error);
                console.error(error);
            }
        };

        // Función para cerrar el modal
        const closeModal = () => {
            showModal.value = false;
        };

        const closeModalTarjetas = () => {
            showModalTarjetas.value = false;
        }
        // Event Listener para detectar la tecla Enter
        const handleKeyPress = (event) => {
            if (showModal.value) {
                if (event.key === "Enter") {
                    closeModal();
                }
            }
        };

        // Mover el scroll al final cuando el componente se monta
        onMounted(() => {
            listarUltimaRelacion();
            document.addEventListener("keydown", handleKeyPress);
            document.addEventListener('keydown', (event) => {
                // Detectar si se presiona F5
                if (event.key === 'F5') {
                    event.preventDefault(); // Evitar el comportamiento predeterminado
                    window.location.reload(); // Recargar la página
                }

                // Detectar si se presiona Ctrl + R (o Cmd + R en macOS)
                if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'r') {
                    event.preventDefault(); // Evitar el comportamiento predeterminado
                    window.location.reload(); // Recargar la página
                }
            });
        });

        onUnmounted(() => {
            document.removeEventListener("keydown", handleKeyPress);
        });

        function mostrarCalendario(inputId) {
            flatpickr(`#${inputId}`, {
                dateFormat: "Y-m-d",
                defaultDate: "today",
            }).open();
        }

        // Estado del menú activo
        const activeMenu = ref(null);

        // Función para alternar el menú
        const toggleMenu = (menuName) => {
            activeMenu.value = activeMenu.value === menuName ? null : menuName;
        };

        // Función para abrir una nueva ventana
        const openWindow = (file) => {
            window.electronAPI.openPage(file);
        };

        // Cerrar el menú al hacer clic fuera
        const closeMenu = (event) => {
            const menuContainers = document.querySelectorAll(".relative");
            const isClickInsideMenu = Array.from(menuContainers).some((container) =>
                container.contains(event.target)
            );
            if (!isClickInsideMenu) {
                activeMenu.value = null;
            }
        };

        // Agregar listener global para cerrar el menú al hacer clic fuera
        window.addEventListener("click", closeMenu);

        const crearOactualizarVehiculo = async () => {
            if (!vehiculo.value.NroAuto || vehiculo.value.NroAuto.trim() === "") {
                return "El campo Número de auto está vacío"
            }

            if (!vehiculo.value.Patente || vehiculo.value.Patente.trim() === "") {
                return "El campo Patente de auto está vacío"
            }
            const vehiculoData = {...vehiculo.value};
            let response = await window.electronAPI.createUpdateVehiculo(vehiculoData);
            if (response.error !== "") {
                return `Error al crear el vehiculo con patente ${vehiculo.value.Patente}`
            } else {
                return ""
            }
        };

        const crearOActualizarUnConductor = async () => {
            if (!vehiculo.value.NroAuto || vehiculo.value.NroAuto.trim() === "" ) {
                return "El campo Número de auto está vacío"
            }

            if (!conductor.value.Documento || conductor.value.Documento.trim() === "") {
                return "El campo Documento del conductor está vacío"
            }

            let conductorData = {
                Documento: conductor.value.Documento,
                Licencia: conductor.value.Licencia,
                Nombre: conductor.value.Nombre,
                ApellidoPa: conductor.value.ApellidoPa,
                ApellidoMa: conductor.value.ApellidoMa,
                Direccion: conductor.value.Direccion,
                Nacionalidad: conductor.value.Nacionalidad,
                Residencia: conductor.value.Residencia,
                Profesion: conductor.value.Profesion,
                FechaNac: conductor.value.FechaNac,
                EstCivil: conductor.value.EstCivil,
                Sexo: conductor.value.Sexo,
                Mva1: conductor.value.Mva1,
                Mva2: conductor.value.Mva1,
                Tipoa: conductor.value.Mva1,
                Nauto: vehiculo.value.NroAuto,
            };

            let response = await window.electronAPI.createUpdateConductor(
                conductorData
            );
            if (response.error !== "") {
                console.error(response.error);
                return response.error;
            }

            return "";
        };

        const crearOActualizarUnPasajero = async () => {
            if (!pasajero.value.Documento || pasajero.value.Documento.trim() === "" ) {
                return "El campo Documento de pasajero está vacío"
            }

            pasajero.value.Mva1 = "X"
            let pasajeroData = {};
            pasajeroData.Documento = pasajero.value.Documento;
            pasajeroData.Licencia = pasajero.value.Licencia.trim() !== "" ? pasajero.value.Licencia : "";
            pasajeroData.Nombre = pasajero.value.Nombre;
            pasajeroData.ApellidoPa = pasajero.value.ApellidoPa;
            pasajeroData.ApellidoMa = pasajero.value.ApellidoMa;
            pasajeroData.Direccion = pasajero.value.Direccion;
            pasajeroData.Nacionalidad = pasajero.value.Nacionalidad;
            pasajeroData.Residencia = pasajero.value.Residencia;
            pasajeroData.Profesion = pasajero.value.Profesion;
            pasajeroData.FechaNac = pasajero.value.FechaNac;
            pasajeroData.EstCivil = pasajero.value.EstCivil;
            pasajeroData.Sexo = pasajero.value.Sexo;
            pasajeroData.Mva1 = pasajero.value.Mva1; // TODO:
            pasajeroData.Mva2 = pasajero.value.Mva1;
            pasajeroData.Tipoa = pasajero.value.Mva1;
            pasajeroData.Nauto = pasajero.value.Nauto.trim() !== "" ? pasajeroData.Nauto.trim() : "";
            pasajeroData.TipoDocumento = pasajero.value.TipoDocumento;

            let response = await window.electronAPI.createUpdateConductor(
                pasajeroData
            );
            if (response.error !== "") {
                console.error(response.error);
                return response.error;
            }

            return "";
        };

        const esperaSegundoEnter = ref(false); // Variable de estado para controlar el segundo Enter

        const focusNextInput = async (currentID, nextInputIdSuccess, nextInputIDFail) => {
            if (currentID === "input-1") {
                const error = await buscarVehiculoPorPatente();
                if (error !== "") {
                    await alertaError(error);
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-6") {
                const [inputFail1, inputFail2] = nextInputIDFail.split(",");
                const error = await buscarVehiculoConConductorPorNroAuto();
                if (error !== "") {
                    await alertaError("No se encontró el número de auto en la base de datos").then(() => {
                        if(vehiculo.value.Patente.trim() === "") {
                            nextInputFocus(inputFail1);
                        }
                        else{
                            nextInputFocus(inputFail2);
                        }
                    });
                    return
                }
                await alertaSuccess("Auto encontrado en la base de datos")
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-12") {
                const error = await buscarConductorPorDocumento();
                if (error !== "") {
                    await alertaError("No se encontró al conductor en la base de datos")
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                await alertaSuccess("Conductor encontrado en la base de datos")
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-10") {
                // const errorV = await buscarVehiculoPorPatente();
                // if (errorV === "") {
                //     alert("el vehiculo ya existe en la base de datos")
                //     nextInputFocus(nextInputIdSuccess)
                //     return
                // }

                const error = await crearOactualizarVehiculo();
                if (error !== "") {
                    await nextInputFocus(nextInputIDFail)
                    await alertaError(`ERROR ${error}`)
                    return
                }
                await alertaSuccess("Se CREO/ACTUALIZO el vehiculo en la base de datos")
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-17") {
                // const errorV = await buscarConductorPorDocumento();
                // if (errorV === "") {
                //     alert("el conductor ya existe en la base de datos")
                //     nextInputFocus(nextInputIdSuccess)
                //     return
                // }

                const error = await crearOActualizarUnConductor();
                if (error !== "") {
                    await alertaError(`ERROR: ${error}`);
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                await alertaSuccess("Se CREO/ACTUALIZO el conductor en la base de datos")
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-21") {
                if (empresa.value.Nombre.trim() !== "") {
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-24") {
                const error = await crearOActualizarUnaEmpresa();
                if (error !== "") {
                    console.log(error);
                    await alertaError(`ERROR: ${error}`);
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                await alertaSuccess("Se guardo/actualizo el registro en la base de datos")
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "dni") {
                if (pasajero.value.Documento.startsWith("http://") || pasajero.value.Documento.startsWith("https://")) {
                    const urlObj = new URL(pasajero.value.Documento);
                    const run = urlObj.searchParams.get("RUN");

                    if (run) {
                        pasajero.value.Documento = run;
                    }
                }

                const error = await buscarPasajeroPorDocumento();
                if (error !== "") {
                    await alertaError("No se encontro el pasajero en la base de datos")
                    await nextInputFocus(nextInputIDFail)
                    return
                }
                showSuccessToast();
                await nextInputFocus(nextInputIdSuccess)
                return;
            }
            if (currentID === "input-35") {
                // const errorV = await buscarPasajeroPorDocumento();
                // if (errorV === "") {
                //     alert("el pasajero ya existe en la base de datos")
                //     nextInputFocus(nextInputIdSuccess)
                //     return
                // }
                if (event.key === "Enter") {
                    if (!esperaSegundoEnter.value) {
                        // Primer Enter: Activamos la espera para el segundo Enter
                        console.log("Presiona Enter nuevamente para continuar...");
                        esperaSegundoEnter.value = true;

                        // Opcional: Añadir un mensaje visual para el usuario
                        // alert("Presiona Enter nuevamente para continuar...");

                        return;
                    } else {
                        // Segundo Enter: Ejecutamos la lógica principal
                        const error = await crearOActualizarUnPasajero();
                        if (error !== "") {
                            console.log(error);
                            await alertaError(`ERROR: ${error}`);
                            await nextInputFocus(nextInputIDFail);
                            esperaSegundoEnter.value = false; // Reiniciamos el estado
                            return;
                        }
                        await alertaSuccess("Se CREO/ACTUALIZO el pasajero en la base de datos");
                        await agregarDato();
                        await nextInputFocus(nextInputIdSuccess);

                        // Reiniciamos el estado para futuras interacciones
                        esperaSegundoEnter.value = false;
                        return;
                    }
                }
            }

            esperaSegundoEnter.value = false
            const nextInput = document.getElementById(nextInputIdSuccess);
            if (nextInput) {
                nextInput.focus();
            }
        };

        const nextInputFocus = async (nextInputID) => {
            const nextInput = document.getElementById(nextInputID);
            if (nextInput) {
                setTimeout(() => {
                    nextInput.focus();
                }, 150); // Retraso de 100 ms
            }
        };

        function showSuccessToast() {
            const toast = document.getElementById("success-toast");
            toast.classList.add("show");

            // Ocultar después de 3 segundos
            setTimeout(() => {
                toast.classList.remove("show");
            }, 1500);
        }

        function showError(message) {
            const toast = document.getElementById("error-toast");
            toast.textContent = message; // Coloca el mensaje de error
            toast.classList.add("show");

            // Ocultar después de 3 segundos
            setTimeout(() => {
                toast.classList.remove("show");
            }, 3000);
        }

        function showSuccessPrint(message) {
            const toast = document.getElementById("success-toast-print");
            toast.querySelector("span").textContent = message;
            toast.classList.add("show");

            // Ocultar después de 3 segundos
            setTimeout(() => {
                toast.classList.remove("show");
            }, 5000);
        }

        async function openPageMain() {
            await crearRelacion()
            window.electronAPI.openPageMain();
        }


        const transformData = async (listado) => {
            const numeracion = Array.from({ length: listado.length }, (_, index) => index + 1);
            return {
                "PASAJERO_NUMERACION":numeracion,
                "PASAJERO_NOMBRES_Y_APELLIDOS": listado.map(pasajero => `${pasajero.Nombres || ""} ${pasajero.ApellidoP || ""} ${pasajero.ApellidoM || ""}`.trim()),
                "PASAJERO_NACIONALIDAD": listado.map(pasajero => pasajero.Nacionalidad || ""),
                "PASAJERO_ESTADO_CIVIL": listado.map(pasajero => pasajero.EstCivil || ""),
                "PASAJERO_FECHA_NACIMIENTO": listado.map(pasajero => pasajero.FechaNac || ""),
                "PASAJERO_PROFESION": listado.map(pasajero => pasajero.Profesion || ""),
                "PASAJERO_NUMERO_DOCUMENTO": listado.map(pasajero => pasajero.Documento || ""),
                "PASAJERO_TIPO_DOCUMENTO": listado.map(pasajero => pasajero.TipoDocumento || "")
            };
        };


        const processData = async (option) => {
            let dataPasajero = {}
            if (option.includes("TARJETA-A-T") || option.includes("TARJETA-T-A")) {
                let listadoFiltrado = [...listado.value];
                if (option.includes("TARJETA-T-A")) {
                    listadoFiltrado = listadoFiltrado.filter(item => item.TarPeru !== "");
                } else if (option.includes("TARJETA-A-T")) {
                    listadoFiltrado = listadoFiltrado.filter(item => item.TarChile !== "");
                }
                dataPasajero = await transformData(listadoFiltrado);
            } else {
                dataPasajero = await transformData(listado.value);
            }

            if (option.includes("CABEZERA")) {
                dataPasajero = {}
            }

            return {
                CABEZERA_FECHA: relacionInfo.value.IsFecha ? relacionInfo.value.FechaSalida : "",
                CABEZERA_PROCEDENCIA: procedencia.value,
                CABEZERA_DESTINO: destino.value,
                VEHICULO_TIPO: vehiculo.value.Tipo,
                VEHICULO_MARCA: vehiculo.value.Marca,
                VEHICULO_MODELO: vehiculo.value.Modelo,
                VEHICULO_AÑO: vehiculo.value.Anio,
                VEHICULO_MOTOR: vehiculo.value.Motor,
                VEHICULO_CHASIS: vehiculo.value.Chasis,
                VEHICULO_PLACA: vehiculo.value.Patente,
                VEHICULO_PAIS: vehiculo.value.Pais,
                VEHICULO_DI: conductor.value.Documento,
                CONDUCTOR_NOMBRE: `${conductor.value.Nombre || ""} ${conductor.value.ApellidoPa || ""} ${conductor.value.ApellidoMa || ""}`.trim(),
                CONDUCTOR_DOMICILIO: conductor.value.Direccion,
                CONDUCTOR_DOCUMENTO: conductor.value.Documento,
                CONDUCTOR_NACIONALIDAD: conductor.value.Nacionalidad,
                CONDUCTOR_PROFESION: conductor.value.Profesion,
                CONDUCTOR_BREVETE: conductor.value.Brevete,
                CONDUCTOR_FECHA_NACIMIENTO: conductor.value.FechaNac,
                ...dataPasajero,
                EMPRESA_AUTORIZADA: empresa.value.Nombre,
                EMPRESA_RESOLUCION_EXENTA: empresa.value.Resolucion,
                EMPRESA_DOCUMENTO_IDONEIDAD: empresa.value.Documento,
                EMPRESA_PERMISO_COMPLEMENTARIO: empresa.value.Permiso
            };
        };

        const deleteRow = (index) => {
            // Elimina el elemento del arreglo usando el índice
            listado.value.splice(index, 1);
            // Recalcular los números después de agregar o eliminar
            listado.value.forEach((item, index) => {
                item.Nro = (index + 1).toString().padStart(2, "0");
            });
        }
        const deleteAllRows = () => {
            // Elimina el elemento del arreglo usando el índice
            listado.value = []
        }

        const crearOActualizarUnaEmpresa = async () => {
            if (!vehiculo.value.NroAuto || vehiculo.value.NroAuto.trim() === "" ) {
                return "El campo Número de auto está vacío"
            }

            let empresaData = {
                documento: empresa.value.Documento,
                nombre: empresa.value.Nombre,
                resolucion: empresa.value.Resolucion,
                nauto: vehiculo.value.NroAuto,
                permiso: empresa.value.Permiso,
            };

            let response = await window.electronAPI.createUpdateEmpresa(
                empresaData
            );
            if (response.error !== "") {
                console.error(response.error);
                return response.error;
            }

            return "";
        };

        const agregarPasajeroAlListado = async() =>{
            const error = await crearOActualizarUnPasajero();
            if (error !== "") {
                console.log(error);
                await alertaError(`ERROR: ${error}`);
                return;
            }
            await alertaSuccess("Se CREO/ACTUALIZO el pasajero en la base de datos");
            await agregarDato();
        }
        watch(
            () => pasajero.value.Nacionalidad,
            (newValue, oldValue) => {
                pasajero.value.Residencia = newValue;
            }
        );

        return {
            agregarPasajeroAlListado,
            empresa,
            deleteAllRows,
            deleteRow,
            closeModalTarjetas,
            openPageMain,
            showError,
            showSuccessToast,
            vehiculo,
            buscarVehiculoPorPatente,
            conductor,
            buscarConductorPorDocumento,
            buscarVehiculoConConductorPorNroAuto,
            pasajero,
            buscarPasajeroPorDocumento,
            listado,
            relacionInfo,
            agregarDato,
            showModal,
            showModalTarjetas,
            tasks,
            crearOActualizarPasajeros,
            crearRelacion,
            crearDetalleRelacion,
            selectOption,
            startProcess,
            closeModal,
            mostrarCalendario,
            activeMenu,
            toggleMenu,
            openWindow,
            focusNextInput
        };
    },
});

app.mount("#app");
