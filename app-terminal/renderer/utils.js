// utils.js

/**
 * Función para imprimir una plantilla con etiquetas y datos.
 * @param {Array} labels - Las etiquetas con sus posiciones.
 * @param {Object} data - Los datos que se insertarán en las etiquetas.
 * @param fontSize
 * @param impresoraActual
 */
export const printTemplate = async (labels, data, fontSize, impresoraActual) => {
    window.electronAPI.imprimirPDF(labels, data, fontSize, impresoraActual);
};

export const printTemplateTest = (labels, data, fontSize) => {
    if (fontSize === "") {
        fontSize = "14px";
    } else {
        fontSize = `${fontSize}px`;
    }

    const printContainer = document.createElement("div");
    printContainer.className = "relative w-[900px] h-[1000px]";

    labels.forEach((label) => {
        const key = label.text.replace(":", "").trim();
        const value = data[key];

        if (Array.isArray(value)) {
            value.forEach((item, index) => {
                const valueElement = document.createElement("div");
                valueElement.className = "absolute text-xs";
                valueElement.style.left = label.left;
                valueElement.style.top = `${
                    parseInt(label.top.replace("px", "")) + index * 20
                }px`;
                valueElement.style.fontSize = fontSize;
                valueElement.textContent = item;
                printContainer.appendChild(valueElement);
            });
        } else {
            const valueElement = document.createElement("div");
            valueElement.className = "absolute text-xs";
            valueElement.style.left = label.left;
            valueElement.style.top = label.top;
            valueElement.style.fontSize = fontSize;
            valueElement.textContent = value || "";

            if (key === "VEHICULO_PLACA") {
                valueElement.style.fontWeight = "bold";
            }

            printContainer.appendChild(valueElement);
        }
    });

    const printWindow = window.open("", "_blank", "width=1100,height=500,resizable=yes,scrollbars=yes");
    printWindow.document.write(`
      <html>
        <head>
          <title>Plantilla para Impresión</title>
          <style>
            body { font-family: Arial, sans-serif; margin: 0; padding: 0; }
            .absolute { position: absolute; }
          </style>
        </head>
        <body>
          ${printContainer.innerHTML}
        </body>
      </html>
    `);
    printWindow.document.close();
    printWindow.print();
};

const toggleInputs = (disable) => {
    const inputs = document.querySelectorAll('input, textarea');
    inputs.forEach(input => {
        input.disabled = disable;
    });
};

const mostrarAlerta = async (tipo, titulo) => {
    toggleInputs(true);

    await Swal.fire({
        title: titulo,
        icon: tipo,
        showConfirmButton: true,
        allowOutsideClick: false,
        didOpen: () => {
            const confirmButton = Swal.getConfirmButton();
            confirmButton.focus();

            // Forzar que el modal siempre esté adelante
            const container = document.querySelector('.swal2-container');
            if (container) {
                container.style.zIndex = "9999";  // Asegurar el z-index alto
                container.style.position = "fixed";  // Evitar que se mueva
            }

            // Cerrar el modal con Enter
            const keyHandler = (event) => {
                if (event.key === 'Enter') {
                    Swal.close();
                    document.removeEventListener('keydown', keyHandler);
                }
            };
            document.addEventListener('keydown', keyHandler);
        }
    });

    toggleInputs(false);
};

export const alertaWarning = async (titulo) => await mostrarAlerta('warning', titulo);
export const alertaSuccess = async (titulo) => await mostrarAlerta('success', titulo);
export const alertaError = async (titulo) => await mostrarAlerta('error', titulo);