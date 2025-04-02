const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("electronAPI", {
    openPage: (page) => ipcRenderer.send("open-page", page),
    openPageMain: () => ipcRenderer.send("open-page-main"),


    createUpdateConductor: (conductor) => ipcRenderer.invoke("create-update-conductor", conductor),
    agregarDetalleRelacion: (detalleRelacion) => ipcRenderer.invoke("agregar-detalle-relacion", detalleRelacion),
    agregarRelacion: (relacion) => ipcRenderer.invoke("agregar-relacion", relacion),
    guardarPlantilla: (titulo, fuente, datos) => ipcRenderer.invoke("guardar-plantilla", titulo, fuente, datos),
    createUpdateVehiculo: (vehiculo) => ipcRenderer.invoke("create-update-vehiculo", vehiculo),
    createUpdateEmpresa: (empresa) => ipcRenderer.invoke("create-update-empresa", empresa),
    actualizarPlantilla: (requestData) => ipcRenderer.invoke("actualizar-plantilla", requestData),
    guardarHistorial:(id) => ipcRenderer.invoke("guardar-historial", id),

    ultimaRelacion: () => ipcRenderer.invoke("ultima-relacion"),
    obtenerVehiculo: (patente) => ipcRenderer.invoke("obtener-vehiculo", patente),
    obtenerConductor: (documento) => ipcRenderer.invoke("obtener-conductor", documento),
    obtenerVehiculoConductor: (nroAuto) => ipcRenderer.invoke("obtener-vehiculo-conductor", nroAuto),
    buscarPlantilla: (id) => ipcRenderer.invoke("buscar-plantilla", id),
    listarPlantillas:() => ipcRenderer.invoke("listar-plantillas"),
    agregarVehiculo: (vehiculo) => ipcRenderer.invoke("agregar-vehiculo", vehiculo),

    savePreferences:(tarjetaAT, tarjetaTA,tarjetaCabezeraAT, tarjetaCabezeraTA, relacionAT, relacionTA, relacionCabezeraAT, relacionCabezeraTA, impresoraActual) => ipcRenderer.invoke("save-preferencias", tarjetaAT, tarjetaTA, tarjetaCabezeraAT, tarjetaCabezeraTA, relacionAT, relacionTA, relacionCabezeraAT, relacionCabezeraTA, impresoraActual),
    getPreferences: () => ipcRenderer.invoke("get-preferencias"),
    historial:(page, page_size) => ipcRenderer.invoke("historial", {page, page_size}),
    eliminarPlantilla: (id) => ipcRenderer.invoke("eliminar-plantilla", id),
    eliminarConductor:(dni) => ipcRenderer.invoke("eliminar-conductor", dni),
    obtenerConductoresPaginados:({page, pageSize, is_conductor, documento, nombre}) => ipcRenderer.invoke("obtener-conductores-paginados",{page, pageSize, is_conductor, documento, nombre}),
    eliminarVehiculo:(patente) => ipcRenderer.invoke("eliminar-vehiculo", patente),
    obtenerVehiculosPaginados:({page, page_size, patente}) => ipcRenderer.invoke("obtener-vehiculos-paginados",{page, page_size, patente}),
    listarImpresoras:() => ipcRenderer.invoke("listar-impresoras"),
    imprimirPDF:(labels, data, fontSize, printerName) => ipcRenderer.invoke("imprimir-pdf", labels, data, fontSize, printerName),
    deleteRelacionFechas:(inicio, fin) => ipcRenderer.invoke("delete-relacion-fechas", inicio, fin),
});

// document.addEventListener("DOMContentLoaded", () => {
//     const fields = document.querySelectorAll("input, textarea");
//     function enforceUppercase(field) {
//         const currentValue = field.value;
//         if (currentValue !== currentValue.toUpperCase()) {
//             const start = field.selectionStart;
//             const end = field.selectionEnd;
//             requestAnimationFrame(() => {
//                 field.value = currentValue.toUpperCase();
//                 field.setSelectionRange(start, end);
//             });
//         }
//     }
//     fields.forEach((field) => {
//         field.addEventListener("input", () => enforceUppercase(field));
//     });
//
//     setInterval(() => {
//         fields.forEach((field) => enforceUppercase(field));
//     }, 0); // Intervalo optimizado
// });

