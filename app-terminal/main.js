import fs from "fs";

import path from "path";
import axios from "axios";
import net from "net";
import {
    app,
    BrowserWindow,
    ipcMain,
    Menu,
    globalShortcut,
    clipboard,
} from "electron";
import {fileURLToPath} from "url";
import PDFDocument from 'pdfkit';
import {Buffer} from 'buffer';
import {spawn} from "child_process";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

let mainWindow;
let loadingWindow;
let apiProcess;

// Crear la ventana de carga
function createLoadingWindow() {
    loadingWindow = new BrowserWindow({
        width: 400,
        height: 300,
        frame: false,
        resizable: false,
        webPreferences: {
            nodeIntegration: true,
            contextIsolation: false,
        },
        icon: path.join(__dirname, 'assets', 'icon.ico') // Cambia la ruta según corresponda
    });
    loadingWindow.loadFile("renderer/loading.html");
}

//TODO
function createMainWindow() {
    mainWindow = new BrowserWindow({
        width: 1366,
        height: 768,
        webPreferences: {
            preload: path.join(__dirname, "preload.js"),
            contextIsolation: true,
            nodeIntegration: false,
            devTools: !produccion,
        },
        icon: path.join(__dirname, 'assets', 'icon.ico') // Cambia la ruta según corresponda
    });
    mainWindow.loadFile("renderer/index.html");

    if (!produccion) {
        // Configurar menú personalizado
        const template = [
            {
                label: "View",
                submenu: [
                    {
                        label: "Toggle DevTools",
                        accelerator: "CmdOrCtrl+Shift+I",
                        click: () => {
                            mainWindow.webContents.toggleDevTools();
                        },
                    },
                ],
            },
        ];
        const menu = Menu.buildFromTemplate(template);
        Menu.setApplicationMenu(menu);
    } else {
        Menu.setApplicationMenu(null);
    }

// Manejar evento de cierre
    mainWindow.on("closed", () => {
        mainWindow = null;
    });
}

// Verificar si el puerto está en uso
function isPortInUse(port) {
    return new Promise((resolve) => {
        const server = net.createServer();
        server.once("error", (err) => {
            if (err.code === "EADDRINUSE") {
                resolve(true); // Puerto en uso
            } else {
                resolve(false); // Otro error
            }
        });
        server.once("listening", () => {
            server.close(); // Cerrar el servidor temporal
            resolve(false); // Puerto no está en uso
        });
        server.listen(port);
    });
}

// Verificar el estado de la API
async function checkAPIHealth() {
    try {
        const response = await axios.get("http://127.0.0.1:8082/health", {
            timeout: 2000,
        });
        return response.status === 200; // Retorna `true` si la API responde correctamente
    } catch {
        return false; // Retorna `false` si hay un error
    }
}

const produccion = true

function startAPI() {
    let apiPath;

    if (produccion) { //TODO
        apiPath = path.join(process.resourcesPath, "api.exe");
    } else {
        apiPath = path.join(__dirname, "api");
    }

    console.log("Intentando ejecutar:", apiPath);

    // Iniciar el proceso de la API en Windows
    apiProcess = spawn(apiPath, [], {
        cwd: path.dirname(apiPath), // Establece el directorio de trabajo
        detached: true, // Permite que el proceso corra independiente
        stdio: ["ignore", "ignore", "ignore"] // Evita bloquear la app por logs
    });

    apiProcess.on("error", (err) => {
        console.error("Error al iniciar la API:", err);
    });

    apiProcess.on("exit", (code) => {
        console.log("La API se cerró con código:", code);
    });

    console.log("API iniciada correctamente.");
}

// Esperar a que la API esté lista
async function waitForAPIHealth(retries = 10, delay = 2000) {
    let attempts = 0;
    while (attempts < retries) {
        try {
            const isHealthy = await checkAPIHealth();
            if (isHealthy) {
                return true; // La API está lista
            }
        } catch {
            // Ignorar errores temporales
        }
        attempts++;
        console.log(`Intento ${attempts} de ${retries}...`);
        await new Promise((resolve) => setTimeout(resolve, delay));
    }
    throw new Error("No se pudo conectar a la API después de varios intentos.");
}

// Registrar atajos globales UNA SOLA VEZ
app.on("ready", () => {
});

// Limpiar atajos al cerrar la aplicación
app.on("will-quit", () => {
    globalShortcut.unregisterAll();
});

// Inicialización de la aplicación
app.whenReady().then(async () => {
    createLoadingWindow();

    try {
        // Verificar si el puerto 8082 está en uso
        const portInUse = await isPortInUse(8082);

        if (portInUse) {
            // Verificar si la API está activa
            const isHealthy = await checkAPIHealth();
            if (isHealthy) {
                console.log("API ya está activa.");
                loadingWindow.close();
                createMainWindow();
                return;
            } else {
                console.error(
                    "El puerto está en uso, pero la API no responde correctamente."
                );
                loadingWindow.webContents.send(
                    "loading-error",
                    "El puerto está en uso, pero la API no responde correctamente."
                );
                return;
            }
        }

        // Si el puerto no está en uso, iniciar la API
        console.log("API no está activa. Intentando iniciar...");
        startAPI();

        // Esperar a que la API esté lista
        await waitForAPIHealth();
        console.log("API iniciada correctamente.");
        loadingWindow.close();
        createMainWindow();
    } catch (error) {
        console.error("Error durante la inicialización:", error.message);
        loadingWindow.webContents.send(
            "loading-error",
            "Error al iniciar la aplicación: " + error.message
        );
    }
});

function registerShortcuts() {
    console.log("Registrando atajos...");
    globalShortcut.register("F5", () => {
        BrowserWindow.getFocusedWindow()?.reload();
    });

    globalShortcut.register("CommandOrControl+R", () => {
        BrowserWindow.getFocusedWindow()?.reload();
    });

    globalShortcut.register("CommandOrControl+C", () => {
        const focusedWindow = BrowserWindow.getFocusedWindow();
        if (focusedWindow) {
            focusedWindow.webContents
                .executeJavaScript("window.getSelection().toString()")
                .then((selectedText) => {
                    if (selectedText) {
                        clipboard.writeText(selectedText);
                        console.log("Texto copiado:", selectedText);
                    }
                });
        }
    });

    globalShortcut.register("CommandOrControl+V", () => {
        const pastedText = clipboard.readText();
        console.log("Texto pegado:", pastedText);
        const focusedWindow = BrowserWindow.getFocusedWindow();
        if (focusedWindow) {
            focusedWindow.webContents.paste();
        }
    });
}

function unregisterShortcuts() {
    console.log("Desregistrando atajos...");
    globalShortcut.unregisterAll();
}

// Registrar atajos solo cuando la app está enfocada
app.on("browser-window-focus", () => {
    unregisterShortcuts(); // Asegurar que no haya duplicados
    registerShortcuts();
});

// Desregistrar atajos cuando la app pierde el foco
app.on("browser-window-blur", unregisterShortcuts);

// Detener la API cuando se cierra la aplicación
app.on("window-all-closed", () => {
    if (apiProcess) {
        apiProcess.kill();
    }
    if (process.platform !== "darwin") {
        app.quit();
    }
});
app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
        createMainWindow();
    }
});

// Función para crear una nueva ventana con un archivo HTML específico
function createNewWindow(file) {
    const newWindow = new BrowserWindow({
        width: 1100,
        height: 600,
        webPreferences: {
            preload: path.join(__dirname, "preload.js"),
            contextIsolation: true,
            nodeIntegration: false,
            devTools: !produccion, // Habilitar DevTools
        },
        icon: path.join(__dirname, 'assets', 'icon.ico') // Cambia la ruta según corresponda
    });
    newWindow.webContents.openDevTools();
    newWindow.loadFile(`renderer/${file}`);
}

// Manejar los eventos desde el frontend
ipcMain.on("open-page", (event, page) => {
    createNewWindow(page);
});

ipcMain.on("open-page-main", () => {
    createMainWindow();
});

app.on("window-all-closed", () => {
    if (process.platform !== "darwin") app.quit();
});

app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
});

///////////
//ENDPOINTS
///////////

ipcMain.handle("create-update-conductor", async (_, conductor) => {
    try {
        const conductorUpper = Object.fromEntries(
            Object.entries(conductor).map(([key, value]) => [
                key.toUpperCase(),
                typeof value === "string" ? value.trim().toUpperCase() : value
            ])
        );

        const response = await axios.post(
            "http://127.0.0.1:8082/api/conductor/create-update",
            conductorUpper
        );
        return response.data;
    } catch (error) {
        console.error("Error enviando conductor:", error);
        return {error: error.message};
    }
});

ipcMain.handle("ultima-relacion", async () => {
    try {
        const response = await axios.get("http://127.0.0.1:8082/api/ultima-relacion");
        return response.data.data;
    } catch (error) {
        console.error("Error en ultimaRelacion:", error);
        return {error: error.message};
    }
});

ipcMain.handle("obtener-vehiculo", async (_, patente) => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/vehiculo/${patente}`
        );
        return response.data.data;
    } catch (error) {
        console.error("Error obteniendo vehículo:", error);
        return {error: error.message};
    }
});

ipcMain.handle("obtener-conductor", async (_, documento) => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/conductor/${documento}`
        );
        return response.data;
    } catch (error) {
        console.error("Error obteniendo conductor:", error);
        return {error: error.message};
    }
});

ipcMain.handle("obtener-vehiculo-conductor", async (_, nroAuto) => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/vehiculo-conductor/${nroAuto}`
        );
        return response.data;
    } catch (error) {
        console.error("Error obteniendo vehículo y conductor:", error);
        return {error: error.message};
    }
});

ipcMain.handle("agregar-detalle-relacion", async (_, detalleRelacion) => {
    try {
        const detalleRelacionUpper = Object.fromEntries(
            Object.entries(detalleRelacion).map(([key, value]) => [
                key.toUpperCase(),
                typeof value === "string" ? value.trim().toUpperCase() : value
            ])
        );
        const response = await axios.post(
            "http://127.0.0.1:8082/api/detalle-relacion",
            detalleRelacionUpper
        );
        return response.data;
    } catch (error) {
        console.error("Error agregando detalle de relación:", error);
        return {error: error.message};
    }
});

ipcMain.handle("agregar-relacion", async (_, relacion) => {
    try {
        const relacionUpper = Object.fromEntries(
            Object.entries(relacion).map(([key, value]) => [
                key.toUpperCase(),
                typeof value === "string" ? value.trim().toUpperCase() : value
            ])
        );

        const response = await axios.post(
            "http://127.0.0.1:8082/api/relacion",
            relacionUpper
        );
        return response.data;
    } catch (error) {
        console.error("Error agregando relación:", error);
        return {error: error.message};
    }
});

ipcMain.handle("guardar-plantilla", async (_, titulo, fuente, datos) => {
    try {
        if (!titulo || titulo.trim() === "") {
            return {error: "El campo 'titulo' es obligatorio y no puede estar vacío"};
        }

        const response = await axios.post(
            "http://127.0.0.1:8082/api/guardar-plantilla",
            {titulo: titulo.toUpperCase(), fuente: fuente, datos: datos}
        );
        return response.data;
    } catch (error) {
        console.error("Error agregando relación:", error);
        return {error: error.message};
    }
});

ipcMain.handle("buscar-plantilla", async (_, id) => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/buscar-plantilla/${id}`
        );
        return response.data;
    } catch (error) {
        console.error("Error agregando relación:", error);
        return {error: error.message};
    }
});

ipcMain.handle("listar-plantillas", async () => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/listar-plantillas`
        );
        return response.data;
    } catch (error) {
        console.error("Error listando plantillas:", error);
        return {error: error.message};
    }
});

ipcMain.handle("listar-impresoras", async () => {
    try {
        const response = await axios.get(
            `http://127.0.0.1:8082/api/impresoras`
        );
        return response.data;
    } catch (error) {
        console.error("Error listando impresoras:", error);
        return {error: error.message};
    }
});

ipcMain.handle("agregar-vehiculo", async (_, vehiculo) => {
    try {
        const response = await axios.post(
            "http://127.0.0.1:8082/api/vehiculo",
            vehiculo
        );
        return response.data;
    } catch (error) {
        console.error("Error agregando vehiculo:", error);
        return {error: error.message};
    }
});

ipcMain.handle("save-preferencias", async (_, tarjetaAT, tarjetaTA, tarjetaCabezeraAT, tarjetaCabezeraTA, relacionAT, relacionTA, relacionCabezeraAT, relacionCabezeraTA, impresoraActual) => {
    // Crear el objeto de preferencias, reemplazando undefined o null con cadenas vacías
    const preferences = {
        "TARJETA-A-T": tarjetaAT || "",
        "TARJETA-T-A": tarjetaTA || "",
        "TARJETA-CABEZERA-A-T": tarjetaCabezeraAT || "",
        "TARJETA-CABEZERA-T-A": tarjetaCabezeraTA || "",
        "RELACION-A-T": relacionAT || "",
        "RELACION-T-A": relacionTA || "",
        "RELACION-CABEZERA-A-T": relacionCabezeraAT || "",
        "RELACION-CABEZERA-T-A": relacionCabezeraTA || "",
        "IMPRESORA-ACTUAL": impresoraActual || "",
    };

    try {
        const response = await axios.post("http://127.0.0.1:8082/api/preferencias", preferences, {
            headers: {
                "Content-Type": "application/json",
            },
        });

        return response.data;
    } catch (error) {
        console.error("Error al guardar preferencias:", error.message);
        return {error: `Error al guardar preferencias: ${error.message}`};
    }
});

ipcMain.handle("get-preferencias", async () => {
    try {
        const response = await axios.get("http://127.0.0.1:8082/api/preferencias", {
            headers: {
                "Content-Type": "application/json",
            },
        });

        const data = response.data.data;

        const filteredData = {};
        for (const key in data) {
            if (data[key] !== "") { // Solo incluir claves con valores no vacíos
                filteredData[key] = data[key];
            }
        }

        return filteredData;
    } catch (error) {
        console.error("Error al obtener preferencias:", error.message);
        return { error: `Error al obtener preferencias: ${error.message}` };
    }
});

ipcMain.handle("historial", async (event, {page, page_size}) => {
    try {
        const response = await axios.get(`http://127.0.0.1:8082/api/historial/paginacion`, {
            params: {page, page_size}
        });
        return response.data;
    } catch (error) {
        console.error("Error historial:", error);
        return {error: error.message};
    }
});

ipcMain.handle("actualizar-plantilla", async (event, requestData) => {
    try {
        let {id, titulo, fuente, datos} = requestData;

        if (!id || id <= 0) {
            return {error: "El campo ID es obligatorio y debe ser un número entero válido"};
        }
        if (!titulo || titulo.trim() === "") {
            return {error: "El campo 'titulo' es obligatorio y no puede estar vacío"};
        }
        if (!Array.isArray(datos) || datos.length === 0) {
            return {error: "El campo 'datos' es obligatorio y no puede estar vacío"};
        }

        titulo = titulo.toUpperCase()
        const response = await axios.put(`http://127.0.0.1:8082/api/plantilla/${id}`, {
            titulo,
            fuente,
            datos,
        });

        return response.data;
    } catch (error) {
        console.error("Error al actualizar la plantilla:", error);
        return {error: error.message};
    }
});

ipcMain.handle("eliminar-plantilla", async (event, id) => {
    try {
        console.log("Eliminando plantilla con ID:", id);
        if (!id || id <= 0) {
            return {error: "El campo ID es obligatorio y debe ser un número entero válido"};
        }
        const response = await axios.delete(`http://127.0.0.1:8082/api/plantilla/${id}`);
        console.log("Respuesta del servidor:", response.data);
        return response.data;
    } catch (error) {
        console.error("Error al eliminar la plantilla:", error.response?.data || error.message);
        return {error: error.message};
    }
});

// Eliminar conductor
ipcMain.handle("eliminar-conductor", async (event, dni) => {
    try {
        console.log("Eliminando conductor con DNI:", dni);
        if (!dni || dni.trim() === "") {
            return {error: "El campo DNI es obligatorio"};
        }
        const response = await axios.delete(`http://127.0.0.1:8082/api/conductor`, {
            params: {doc: dni}
        });
        console.log("Respuesta del servidor:", response.data);
        return response.data;
    } catch (error) {
        console.error("Error al eliminar el conductor:", error.response?.data || error.message);
        return {error: error.message};
    }
});

// Obtener conductores paginados
ipcMain.handle("obtener-conductores-paginados", async (event, {page, page_size, is_conductor, documento, nombre}) => {
    try {
        const response = await axios.get(`http://127.0.0.1:8082/api/conductores/paginacion`, {
            params: {page, page_size, is_conductor, documento, nombre},
        });
        return response.data;
    } catch (error) {
        console.log(error)
        console.error("Error al obtener conductores paginados:", error.response?.data || error.message);
        return {error: error.message};
    }
});

// Eliminar vehículo
ipcMain.handle("eliminar-vehiculo", async (event, patente) => {
    try {
        console.log("Eliminando vehículo con patente:", patente);
        if (!patente) {
            return {error: "El campo patente es obligatorio"};
        }
        const response = await axios.delete(`http://127.0.0.1:8082/api/vehiculo`, {
            params: {patente: patente}
        });
        console.log("Respuesta del servidor:", response.data);
        return response.data;
    } catch (error) {
        console.error("Error al eliminar el vehículo:", error.response?.data || error.message);
        return {error: error.message};
    }
});

// Obtener vehículos paginados
ipcMain.handle("obtener-vehiculos-paginados", async (event, {page, page_size, patente}) => {
    try {
        const response = await axios.get(`http://127.0.0.1:8082/api/vehiculos/paginacion`, {
            params: {page, page_size, patente},
        });
        console.log("Respuesta del servidor:", response.data);
        return response.data;
    } catch (error) {
        console.error("Error al obtener vehículos paginados:", error.response?.data || error.message);
        return {error: error.message};
    }
});

ipcMain.handle("create-update-vehiculo", async (_, vehiculo) => {
    try {
        const vehiculoUpper = Object.fromEntries(
            Object.entries(vehiculo).map(([key, value]) => [
                key.toUpperCase(),
                typeof value === "string" ? value.trim().toUpperCase() : value
            ])
        );

        const response = await axios.post(
            "http://127.0.0.1:8082/api/vehiculo/create-update",
            vehiculoUpper
        );
        return response.data;
    } catch (error) {
        console.error("Error enviando vehiculo:", error);
        return {error: error.message};
    }
});

ipcMain.handle("create-update-empresa", async (_, empresa) => {
    try {
        const empresaUpper = Object.fromEntries(
            Object.entries(empresa).map(([key, value]) => [
                key.toUpperCase(),
                typeof value === "string" ? value.trim().toUpperCase() : value
            ])
        );
        const response = await axios.post(
            "http://127.0.0.1:8082/api/empresa",
            empresaUpper
        );
        return response.data;
    } catch (error) {
        console.error("Error guardando empresa:", error);
        return {error: error.message};
    }
});

ipcMain.handle("delete-relacion-fechas", async (_, fechaInicio, fechaFin) => {
    try {
        const response = await axios.delete("http://127.0.0.1:8082/api/historial/borrar-fechas", {
            params: { // Enviar las fechas como query params
                inicio: fechaInicio,
                fin: fechaFin,
            },
        });
        return response.data; // Retornar la respuesta del servidor
    } catch (error) {
        console.error("Error eliminando relaciones por fechas:", error);
        return {error: error.message}; // Retornar el error
    }
});

ipcMain.handle("guardar-historial", async (_, id) => {
    try {
        const response = await axios.post(`http://127.0.0.1:8082/api/historial/guardar/${id}`);
        return response.data;
    } catch (error) {
        console.error("Error guardando el historial:", error);
        return {error: error.message};
    }
});

/**
 * Función para generar un PDF en memoria y enviarlo al backend como Base64.
 */
ipcMain.handle('imprimir-pdf', async (_, labels, data, fontSize, printerName) => {
    if (!printerName) {
        return {error: 'Seleccione una impresora para imprimir'};
    }
    try {
        // Generar el PDF en memoria
        const pdfBuffer = await generatePDFInMemory(labels, data, fontSize);
        //fs.writeFileSync('output.pdf', pdfBuffer)
        // Convertir el Buffer a Base64
        const base64Data = pdfBuffer.toString('base64');

        // Enviar el archivo Base64 al backend
        const response = await axios.post('http://127.0.0.1:8082/api/pdf/upload', {
            pdf: base64Data,
            printer_name: printerName,
        });

        return response.data;
    } catch (error) {
        console.error('Error generando o enviando el PDF:', error);
        return {error: error.message};
    }
});

function generatePDFInMemory(labels, data, fontSize) {
    if (!fontSize || fontSize === "") {
        fontSize = "8";
    }

    return new Promise((resolve, reject) => {
        try {
            // Convertir fontSize a número entero
            const parsedFontSize = parseInt(fontSize, 10);
            if (isNaN(parsedFontSize) || parsedFontSize <= 0) {
                throw new Error(`El valor de fontSize (${fontSize}) no es válido. Debe ser un número positivo.`);
            }

            // Dimensiones de la hoja oficio en puntos
            const widthInPoints = 900; // Ancho de hoja oficio
            const heightInPoints = 1000; // Alto de hoja oficio

            const widthInPixels = 900;
            const ratio = widthInPoints / widthInPixels; // Relación puntos/píxeles

            // Crear un nuevo documento PDF con el tamaño personalizado
            const doc = new PDFDocument({
                size: [widthInPoints, heightInPoints], // Tamaño hoja oficio
                margin: 0,
                layout: 'portrait', // Asegúrate de que el layout sea el correcto
            });

            // Almacenar los datos del PDF en un array de buffers
            const buffers = [];
            doc.on('data', buffers.push.bind(buffers));
            doc.on('end', () => {
                // Combinar los buffers en uno solo
                const pdfBuffer = Buffer.concat(buffers);
                resolve(pdfBuffer);
            });

            // Agregar contenido al PDF basado en las etiquetas y datos
            labels.forEach((label) => {
                const key = label.text.replace(':', '').trim();
                const value = data[key];

                // Convertir left y top a números (asumiendo que están en formato "Xpx")
                const leftInPixels = parseInt(label.left.replace("px", ""), 10);
                const topInPixels = parseInt(label.top.replace("px", ""), 10);

                // Convertir las coordenadas de píxeles a puntos
                const leftInPoints = leftInPixels * ratio;
                const topInPoints = topInPixels * ratio;

                // Verificar si el label es CONDUCTOR_DOCUMENTO o VEHICULO_PLACA
                const isSpecialLabel = key === "CONDUCTOR_DOCUMENTO" || key === "VEHICULO_PLACA";

                // Aplicar estilos especiales si es necesario
                if (isSpecialLabel) {
                    doc.fontSize(parsedFontSize + 3) // Aumentar el tamaño de la fuente en +3
                        .font('Helvetica-Bold'); // Usar fuente en negrita
                } else {
                    doc.fontSize(parsedFontSize)
                        .font('Helvetica'); // Fuente normal
                }

                if (Array.isArray(value)) {
                    value.forEach((item, index) => {
                        doc.text(
                            String(item).toUpperCase(),
                            leftInPoints,
                            topInPoints + index * (parsedFontSize + 4) // Incrementar verticalmente si hay múltiples valores
                        );
                    });
                } else {
                    doc.text(
                        value ? String(value).toUpperCase() : '',
                        leftInPoints,
                        topInPoints
                    );
                }
            });

            // Finalizar el documento
            doc.end();
        } catch (error) {
            reject(error);
        }
    });
}