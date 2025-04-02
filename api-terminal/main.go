package main

import (
	"api-terminal/conn"
	"api-terminal/repository"
	"api-terminal/service"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	cn, errN := conn.NewConn()
	if errN != nil {
		log.Fatal(errN)
	}

	repo := repository.NewRepository(cn)
	conductor := service.NewConductorService(repo)
	detalle := service.NewDetalleRelacionService(repo)
	relacion := service.NewRelacionService(repo)
	vehiculo := service.NewVehiculoService(repo)
	plantilla := service.NewPlantillaService()
	empresaService := service.NewEmpresaService(repo)
	historialService := service.NewHistorialService(repo)
	preferenceService := service.NewPreferenceService()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON("VERSION 1.0")
	})
	// Endpoint 1: GET /api/vehiculo/:patente
	app.Get("/api/vehiculo/:patente", func(c *fiber.Ctx) error {
		patente := c.Params("patente")

		v, err := vehiculo.ObtenerVehiculoPorPatente(patente)
		if err != nil {
			log.Errorf("/api/vehiculo/:patente", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"patente": patente,
				"data":    v})
		}

		return c.Status(200).JSON(fiber.Map{
			"error":   "",
			"patente": patente,
			"data":    v})
	})

	// Endpoint 2: GET /api/conductor/:documento
	app.Get("/api/conductor/:documento", func(c *fiber.Ctx) error {
		documento := c.Params("documento")
		// Lógica para buscar conductor por documento
		d, err := conductor.ObtenerConductorPorDocumento(documento)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":     err.Error(),
				"documento": documento,
				"data":      d})
		}

		return c.JSON(fiber.Map{
			"error":     "",
			"documento": documento,
			"data":      d})
	})

	// Endpoint 3: GET /api/vehiculo-conductor/:nroAuto
	app.Get("/api/vehiculo-conductor/:nroAuto", func(c *fiber.Ctx) error {
		nroAuto := c.Params("nroAuto")
		// Lógica para buscar vehículo y conductor por NroAuto
		a, err := vehiculo.ObtenerVehiculoConConductorPorNroAuto(nroAuto)
		if err != nil {
			log.Errorf("/api/vehiculo-conductor/:nroAuto", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"nroAuto": nroAuto,
				"data":    a})
		}

		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok",
			"nroAuto": nroAuto,
			"data":    a})
	})

	// Endpoint 4: GET /api/ultima-relacion
	app.Get("/api/ultima-relacion", func(c *fiber.Ctx) error {
		// Lógica para obtener la última relación
		r, err := relacion.ObtenerUltimaRelacion()
		if err != nil {
			log.Errorf("/api/ultima-relacion", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error": err.Error(),
					"data":  r})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok",
			"data":    r})
	})

	// Endpoint 5: POST /api/conductor
	app.Post("/api/conductor/create-update", func(c *fiber.Ctx) error {
		var conductorBody repository.Conductor
		if err := c.BodyParser(&conductorBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		err := conductor.CrearOActualizarConductor(conductorBody)
		if err != nil {
			log.Errorf("/api/conductor/create-update", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error()})
		}

		// Lógica para crear o actualizar conductor
		return c.JSON(fiber.Map{"message": "Conductor creado o actualizado", "error": ""})
	})

	// Endpoint 6: POST /api/detalle-relacion
	app.Post("/api/detalle-relacion", func(c *fiber.Ctx) error {
		var dr repository.DetalleRelacion
		if err := c.BodyParser(&dr); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		err := detalle.AgregarDetalleRelacion(dr)
		if err != nil {
			log.Errorf("/api/detalle-relacion", err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "fallo al agregar detalle relacion"})
		}

		// Lógica para agregar detalle de relación
		return c.JSON(fiber.Map{
			"error": ""})
	})

	// Endpoint 7: POST /api/relacion
	app.Post("/api/relacion", func(c *fiber.Ctx) error {
		var r repository.Relacion
		if err := c.BodyParser(&r); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		err := relacion.AgregarRelacion(r)
		if err != nil {
			log.Errorf("/api/relacion", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"message": "fallo al agregar relacion",
				"error":   err.Error(),
			})
		}

		// Lógica para agregar relación
		return c.JSON(fiber.Map{"message": "Relación agregada", "error": ""})
	})

	app.Post("/api/guardar-plantilla", func(c *fiber.Ctx) error {
		var im repository.Impresion
		if err := c.BodyParser(&im); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}
		// Lógica para obtener la última relación
		err := plantilla.GuardarPlantilla(im.Titulo, im.Fuente, im.Datos)
		if err != nil {
			log.Errorf("/api/guardar-plantilla", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok"})
	})

	app.Get("/api/buscar-plantilla/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
		}

		r, err := plantilla.BuscarPlantilla(id)
		if err != nil {
			log.Errorf("/api/buscar-plantilla/:id", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error": err.Error(),
					"data":  r})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok",
			"data":    r})
	})

	app.Get("/api/listar-plantillas", func(c *fiber.Ctx) error {
		r, err := plantilla.ObtenerTodasLasPlantillas()
		if err != nil {
			log.Errorf("/api/listar-plantillas", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error": err.Error(),
					"data":  r})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok",
			"data":    r})
	})

	app.Post("/api/vehiculo", func(c *fiber.Ctx) error {
		var r repository.Vehiculo
		if err := c.BodyParser(&r); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		err := vehiculo.AgregarVehiculo(r)
		if err != nil {
			log.Errorf("/api/vehiculo", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"message": "fallo al agregar el vehiculo",
				"error":   err.Error(),
			})
		}

		// Lógica para agregar relación
		return c.JSON(fiber.Map{"message": "Relación agregada", "error": ""})
	})

	app.Get("/api/historial", func(c *fiber.Ctx) error {
		// Obtener parámetros de la URL
		page := c.QueryInt("page", 1)           // Valor predeterminado: 1
		pageSize := c.QueryInt("page_size", 10) // Valor predeterminado: 10

		// Obtener relaciones
		r, totalRecords, totalPages, err := relacion.ObtenerRelaciones(page, pageSize)
		if err != nil {
			log.Errorf("/api/historial", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error obteniendo el historial",
				"data":    nil,
			})
		}

		// Retornar respuesta exitosa
		return c.JSON(fiber.Map{
			"error":        "",
			"message":      "ok",
			"data":         r,
			"totalRecords": totalRecords,
			"totalPages":   totalPages,
		})
	})

	app.Put("/api/plantilla/:id", func(c *fiber.Ctx) error {
		// Obtener el ID de los parámetros de la URL
		id, err := c.ParamsInt("id", 0)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID inválido o faltante",
				"error":   "El campo ID es obligatorio y debe ser un número entero válido",
			})
		}

		// Parsear el cuerpo de la solicitud
		var requestData repository.Impresion
		if err := c.BodyParser(&requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Error al parsear el JSON",
				"error":   err.Error(),
			})
		}

		// Validar que los campos requeridos estén presentes
		if requestData.Titulo == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "El campo 'titulo' es obligatorio",
				"error":   "El título no puede estar vacío",
			})
		}
		if len(requestData.Datos) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "El campo 'datos' es obligatorio",
				"error":   "Los datos no pueden estar vacíos",
			})
		}

		// Llamar al servicio para actualizar la plantilla
		err = plantilla.ActualizarPlantilla(id, requestData.Titulo, requestData.Fuente, requestData.Datos)
		if err != nil {
			log.Errorf("/api/plantilla/:id", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al actualizar la plantilla",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Plantilla actualizada correctamente",
			"error":   "",
		})
	})

	app.Delete("/api/plantilla/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id", 0)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID inválido o faltante",
				"error":   "El campo ID es obligatorio y debe ser un número entero válido",
			})
		}

		// Llamar al servicio para eliminar la plantilla
		err = plantilla.EliminarPlantilla(id)
		if err != nil {
			log.Errorf("/api/plantilla/:id", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al eliminar la plantilla",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Plantilla eliminada correctamente",
			"error":   "",
		})
	})

	app.Delete("/api/conductor", func(c *fiber.Ctx) error {
		doc := c.Query("doc", "")
		if doc == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID inválido o faltante",
				"error":   "El campo ID es obligatorio y debe ser un número entero válido",
			})
		}

		err := conductor.EliminarConductor(doc)
		if err != nil {
			log.Errorf("/api/conductor", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al eliminar el conductor",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Conductor eliminada correctamente",
			"error":   "",
		})
	})

	app.Get("/api/conductores/paginacion", func(c *fiber.Ctx) error {
		// Obtener parámetros de la URL
		page := c.QueryInt("page", 1)                     // Valor predeterminado: 0
		page_size := c.QueryInt("page_size", 10)          // Valor predeterminado: 10
		isConductor := c.QueryBool("is_conductor", false) // Valor predeterminado: false
		documento := c.Query("documento", "")
		nombre := c.Query("nombre", "")
		// Obtener relaciones
		r, totalRecords, totalPages, err := conductor.ListarConductoresPaginados(page, page_size, isConductor, documento, nombre)
		if err != nil {
			log.Errorf("/api/conductores/paginacion", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error obteniendo los conductores",
				"data":    nil,
			})
		}

		// Retornar respuesta exitosa
		return c.JSON(fiber.Map{
			"error":        "",
			"message":      "ok",
			"data":         r,
			"totalRecords": totalRecords,
			"totalPages":   totalPages,
		})
	})
	app.Delete("/api/vehiculo", func(c *fiber.Ctx) error {
		patente := c.Query("patente", "")
		if patente == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "patente inválido o faltante",
				"error":   "El campo patente es obligatorio y debe ser un número entero válido",
			})
		}

		// Llamar al servicio para eliminar la plantilla
		err := vehiculo.EliminarVehiculo(patente)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al eliminar el vehiculo",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "vehiculo eliminado correctamente",
			"error":   "",
		})
	})

	app.Get("/api/vehiculos/paginacion", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)           // Valor predeterminado: 1
		pageSize := c.QueryInt("page_size", 10) // Valor predeterminado: 10
		patente := c.Query("patente", "")
		vehiculos, totalRecords, totalPages, err := vehiculo.ListarVehiculosPaginados(page, pageSize, patente)
		if err != nil {
			log.Errorf("/api/vehiculos/paginacion", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error obteniendo los vehiculos",
				"data":    nil,
			})
		}

		// Retornar respuesta exitosa
		return c.JSON(fiber.Map{
			"error":        "",
			"message":      "ok",
			"data":         vehiculos,
			"totalRecords": totalRecords,
			"totalPages":   totalPages,
		})
	})

	// Endpoint 5: POST /api/conductor
	app.Post("/api/vehiculo/create-update", func(c *fiber.Ctx) error {
		var vehiculoBody repository.Vehiculo
		if err := c.BodyParser(&vehiculoBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		err := vehiculo.CrearOActualizarVehiculo(vehiculoBody)
		if err != nil {
			log.Errorf("/api/conductor/vehiculo-create-update", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error()})
		}

		// Lógica para crear o actualizar Vehiculo
		return c.JSON(fiber.Map{"message": "Vehiculo creado o actualizado", "error": ""})
	})

	// Ruta para crear o actualizar una empresa
	app.Post("/api/empresa", func(c *fiber.Ctx) error {
		var empresa repository.Empresa

		// Parsear el cuerpo de la solicitud
		if err := c.BodyParser(&empresa); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Datos inválidos",
				"error":   "No se pudo parsear el cuerpo de la solicitud",
			})
		}

		// Llamar al servicio para crear o actualizar la empresa
		err := empresaService.CrearOActualizarEmpresa(empresa)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al crear o actualizar la empresa",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Empresa creada o actualizada correctamente",
			"error":   "",
		})
	})

	// Ruta para obtener todas las empresas
	app.Get("/api/empresas", func(c *fiber.Ctx) error {
		// Llamar al servicio para obtener todas las empresas
		empresas, err := empresaService.ObtenerEmpresas()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al obtener las empresas",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Empresas obtenidas correctamente",
			"data":    empresas,
			"error":   "",
		})
	})

	// Ruta para obtener una empresa por ID
	app.Get("/api/empresa", func(c *fiber.Ctx) error {
		// Obtener el parámetro 'nauto' de la consulta
		nauto := c.Query("nauto", "")
		if nauto == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID de empresa inválido o faltante",
				"error":   "El campo nauto es obligatorio",
			})
		}

		// Llamar al servicio para obtener la empresa por ID
		empresa, err := empresaService.ObtenerEmpresaPorID(nauto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al obtener la empresa",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Empresa obtenida correctamente",
			"data":    empresa,
			"error":   "",
		})
	})

	// Ruta para eliminar una empresa
	app.Delete("/api/empresa", func(c *fiber.Ctx) error {
		// Obtener el parámetro 'nauto' de la consulta
		nauto := c.Query("nauto", "")
		if nauto == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID de empresa inválido o faltante",
				"error":   "El campo nauto es obligatorio",
			})
		}

		// Llamar al servicio para eliminar la empresa
		err := empresaService.EliminarEmpresa(nauto)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Fallo al eliminar la empresa",
				"error":   err.Error(),
			})
		}

		// Respuesta exitosa
		return c.JSON(fiber.Map{
			"message": "Empresa eliminada correctamente",
			"error":   "",
		})
	})

	app.Delete("/api/relacion/borrar-fechas", func(c *fiber.Ctx) error {
		inicio := c.Query("inicio", "")
		if inicio == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "",
				"error":   "La fecha inicio esta vacia",
			})
		}

		fin := c.Query("fin", "")
		if inicio == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "",
				"error":   "La fecha fin esta vacia",
			})
		}

		err := relacion.EliminarRelacionesPorFecha(inicio, fin)
		if err != nil {
			log.Errorf("/api/relacion/borrar-fechas", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error":   err.Error(),
					"message": ""})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok"})
	})

	app.Post("/api/pdf/upload", func(c *fiber.Ctx) error {
		type PrintRequest struct {
			PDF         string `json:"pdf"`          // PDF en Base64
			PrinterName string `json:"printer_name"` // Nombre de la impresora
		}
		var req PrintRequest

		// Parsear la solicitud
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Solicitud inválida",
				"message": err.Error(),
			})
		}

		// Decodificar el PDF desde Base64
		pdfData, err := base64.StdEncoding.DecodeString(req.PDF)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Error al decodificar el PDF",
				"message": err.Error(),
			})
		}

		// Detectar el sistema operativo y enviar el PDF a la impresora
		switch os := runtime.GOOS; os {
		case "windows":
			err = printPDFWindows(pdfData, req.PrinterName)
		case "darwin", "linux":
			err = printPDFUnix(pdfData, req.PrinterName)
		default:
			err = fiber.NewError(fiber.StatusInternalServerError, "Sistema operativo no soportado")
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Error al imprimir el PDF",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"error":   "",
			"message": "PDF enviado a la impresora correctamente",
		})
	})

	app.Get("/api/impresoras", func(c *fiber.Ctx) error {
		impresoras, actual, err := obtenerImpresoras()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "No se pudieron obtener las impresoras",
				"message": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"impresoras": impresoras,
			"actual":     actual,
		})
	})

	app.Get("/api/historial/paginacion", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)           // Valor predeterminado: 1
		pageSize := c.QueryInt("page_size", 10) // Valor predeterminado: 10

		historial, err := historialService.ObtenerHistorialPaginado(page, pageSize)
		if err != nil {
			log.Errorf("/api/historial/paginacion: %s", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Error obteniendo el historial",
				"data":    nil,
			})
		}

		// Convertir cada elemento de historial.Data a un objeto JSON
		for i := range historial.Data {
			// 🔹 Afirmación de tipo: convertir historial.Data[i] a string
			datosStr, ok := historial.Data[i].(string)
			if !ok {
				log.Errorf("ERROR: historial.Data[%d] no es de tipo string", i)
				continue
			}

			// 🔹 Limpiar y formatear el JSON
			jsonData, err := limpiarJSON(datosStr)
			if err != nil {
				log.Errorf("ERROR paginacion historial: %s", err.Error())
				continue
			}

			// 🔹 Asignar el objeto JSON a historial.Data[i]
			historial.Data[i] = jsonData
		}

		// 🔹 Devolver la respuesta con los datos formateados
		return c.JSON(fiber.Map{
			"data":         historial.Data,
			"error":        "",
			"message":      "ok",
			"totalPages":   historial.TotalPages,
			"totalRecords": historial.TotalRecords,
		})
	})

	app.Delete("/api/historial/borrar-fechas", func(c *fiber.Ctx) error {
		inicio := c.Query("inicio", "")
		if inicio == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "",
				"error":   "La fecha inicio esta vacia",
			})
		}

		fin := c.Query("fin", "")
		if inicio == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "",
				"error":   "La fecha fin esta vacia",
			})
		}

		err := relacion.EliminarRelacionesPorFecha(inicio, fin)
		if err != nil {
			log.Errorf("/api/relacion/borrar-fechas", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error":   err.Error(),
					"message": ""})
		}
		err = historialService.EliminarPorfechas(inicio, fin)
		if err != nil {
			log.Errorf("/api/historial/relacion/borrar-fechas", err.Error())
			return c.Status(500).JSON(
				fiber.Map{
					"error":   err.Error(),
					"message": ""})
		}
		return c.JSON(fiber.Map{
			"error":   "",
			"message": "ok"})
	})

	app.Post("/api/historial/guardar/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") // Obtiene el valor del parámetro `id`
		if strings.TrimSpace(id) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID es vacío"})
		}

		// Enviar respuesta inmediata y ejecutar en segundo plano
		go func() {
			data, err := relacion.ObtenerRelacionPorID(id)
			if err != nil {
				log.Errorf("/api/historial/guardar: fallo al ObtenerRelacionPorID: %s", err.Error())
				return
			}

			jsonBytes, err := json.Marshal(data)
			if err != nil {
				log.Errorf("/api/historial/guardar: fallo al Marshal: %s", err.Error())
				return
			}

			err = historialService.AgregarHistorial(string(jsonBytes))
			if err != nil {
				log.Errorf("/api/historial/guardar: fallo al agregar el historial: %s", err.Error())
				return
			}
		}()

		// Responder de inmediato sin esperar la goroutine
		return c.JSON(fiber.Map{"message": "Historial agregado", "error": ""})
	})

	app.Post("/api/preferencias", func(c *fiber.Ctx) error {
		var prefs map[string]interface{}
		if err := c.BodyParser(&prefs); err != nil {
			log.Errorf("Error al parsear el cuerpo de la solicitud: %v", err)
			return c.Status(400).JSON(fiber.Map{
				"error": "Cuerpo de la solicitud inválido",
				"data":  nil,
			})
		}

		err := preferenceService.SavePreferencias(prefs)
		if err != nil {
			log.Errorf("Error al guardar preferencias: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Error al guardar preferencias",
				"data":  nil,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"error": "",
			"data":  "Preferencias guardadas correctamente",
		})
	})

	app.Get("/api/preferencias", func(c *fiber.Ctx) error {
		prefs, err := preferenceService.GetPreferencias()
		if err != nil {
			log.Errorf("Error al obtener preferencias: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Error al obtener preferencias",
				"data":  nil,
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"error": "",
			"data":  prefs,
		})
	})
	// Iniciar el servidor en localhost:8080
	log.Info("API running on http://localhost:8082")
	log.Fatal(app.Listen(":8082"))
}

func obtenerImpresoras() ([]string, string, error) {
	var cmd *exec.Cmd
	var currentPrinter string

	// Obtener la lista de impresoras
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "printer", "get", "name")
	case "darwin":
		cmd = exec.Command("lpstat", "-a")
	default:
		return nil, "", fmt.Errorf("sistema operativo no soportado")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, "", err
	}

	// Procesar la salida para obtener las impresoras
	lines := strings.Split(string(out), "\n")
	var impresoras []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" && name != "Name" {
			if runtime.GOOS == "windows" {
				impresoras = append(impresoras, name)
			} else if runtime.GOOS == "darwin" {
				impresoras = append(impresoras, strings.Fields(name)[0])
			}
		}
	}

	// Obtener la impresora predeterminada
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "printer", "where", "Default=true", "get", "Name")
	case "darwin":
		cmd = exec.Command("lpstat", "-d")
	default:
		return nil, "", fmt.Errorf("sistema operativo no soportado")
	}

	out, err = cmd.Output()
	if err != nil {
		return nil, "", err
	}

	// Limpiar la salida para obtener el nombre de la impresora predeterminada
	currentPrinter = strings.TrimSpace(string(out))
	if runtime.GOOS == "windows" {
		currentPrinter = strings.ReplaceAll(currentPrinter, "Name", "")
		currentPrinter = strings.TrimSpace(currentPrinter)
	} else if runtime.GOOS == "darwin" {
		currentPrinter = strings.ReplaceAll(currentPrinter, "system default destination:", "")
		currentPrinter = strings.TrimSpace(currentPrinter)
	}

	return impresoras, currentPrinter, nil
}

func printPDFWindows(pdfData []byte, printerName string) error {
	fmt.Println("Iniciando proceso de impresión...")

	// Crear un archivo temporal para almacenar el PDF
	tempDir := os.TempDir()
	//fmt.Printf("Directorio temporal: %s\n", tempDir)

	tempFile, err := ioutil.TempFile(tempDir, "*.pdf")
	if err != nil {
		fmt.Printf("Error al crear el archivo temporal: %v\n", err)
		return fmt.Errorf("no se pudo crear el archivo temporal: %v", err)
	}
	tempFilePath := tempFile.Name()
	//fmt.Printf("Archivo temporal creado en: %s\n", tempFilePath)

	// Cerrar antes de escribir manualmente
	tempFile.Close()

	// Escribir los datos del PDF en el archivo temporal
	fmt.Println("Escribiendo datos del PDF en el archivo temporal...")
	if err := ioutil.WriteFile(tempFilePath, pdfData, 0644); err != nil {
		fmt.Printf("Error al escribir en el archivo temporal: %v\n", err)
		return fmt.Errorf("no se pudo escribir en el archivo temporal: %v", err)
	}

	// Verificar que el archivo realmente existe antes de continuar
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: el archivo temporal no existe: %s\n", tempFilePath)
		return fmt.Errorf("el archivo temporal no existe: %s", tempFilePath)
	}
	//fmt.Println("Archivo temporal verificado correctamente.")

	// Ruta completa al ejecutable de SumatraPDF (ajústala si es necesario)
	sumatraPDFPath := "./SumatraPDF.exe"
	//fmt.Printf("Usando SumatraPDF en: %s\n", sumatraPDFPath)

	// Obtener la ruta absoluta del archivo temporal
	absPDFPath, err := filepath.Abs(tempFilePath)
	if err != nil {
		fmt.Printf("Error al obtener la ruta absoluta del archivo PDF: %v\n", err)
		return fmt.Errorf("error al obtener la ruta absoluta del archivo PDF: %v", err)
	}
	//fmt.Printf("Ruta absoluta del PDF: %s\n", absPDFPath)

	// Construir el comando para imprimir el PDF con ajustes de orientación
	//cmd := exec.Command(sumatraPDFPath, "-print-to", printerName, "-print-settings", "portrait,fit", absPDFPath)
	cmd := exec.Command(sumatraPDFPath, "-print-to", printerName, "-print-settings", "portrait,noscale,A4", absPDFPath)
	//fmt.Printf("Ejecutando comando: %s %s %s %s %s %s\n", sumatraPDFPath, "-print-to", printerName, "-print-settings", "portrait,fit", absPDFPath)

	// Ejecutar el comando y capturar la salida
	output, err := cmd.CombinedOutput()
	//fmt.Printf("Salida del comando SumatraPDF:\n%s\n", string(output))

	if err != nil {
		fmt.Printf("Error al ejecutar SumatraPDF: %v\n", err)
		return fmt.Errorf("error al ejecutar SumatraPDF: %v, salida: %s", err, string(output))
	}

	// Esperar a que el proceso de impresión termine
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error al esperar la finalización del proceso de impresión: %v\n", err)
		return fmt.Errorf("error al esperar la finalización del proceso de impresión: %v", err)
	}
	//fmt.Println("Impresión completada con éxito.")

	go func(path string) {
		time.Sleep(1 * time.Minute)
		if err := os.Remove(path); err != nil {
			fmt.Printf("Error al eliminar archivo temporal: %v\n", err)
		} else {
			fmt.Printf("Archivo temporal eliminado: %s\n", path)
		}
	}(tempFilePath)

	fmt.Println("Proceso de impresión finalizado.")
	return nil
}

func printPDFUnix(pdfData []byte, printerName string) error {
	// Crear el comando lpr
	cmd := exec.Command("lpr", "-P", printerName)

	// Obtener el pipe de entrada del comando
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error al obtener StdinPipe: %v", err)
	}

	// Capturar la salida estándar y de error
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// Iniciar el comando
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error al iniciar el comando: %v", err)
	}

	// Escribir los datos del PDF en el pipe de entrada
	if _, err := stdin.Write(pdfData); err != nil {
		return fmt.Errorf("error al escribir en el pipe: %v", err)
	}

	// Cerrar el pipe de entrada
	if err := stdin.Close(); err != nil {
		return fmt.Errorf("error al cerrar el pipe: %v", err)
	}

	// Esperar a que el comando finalice
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error al esperar el comando: %v, stderr: %s", err, errb.String())
	}

	// Imprimir la salida estándar (si existe)
	if outb.Len() > 0 {
		fmt.Println("Salida estándar:", outb.String())
	}

	return nil
}

func limpiarJSON(datos string) (map[string]interface{}, error) {
	// 🔹 1. Desescapar la cadena JSON
	var jsonData map[string]interface{}
	datosLimpios := strings.ReplaceAll(datos, `\"`, `"`)

	// 🔹 2. Remover comillas extra si envuelve todo el JSON
	if strings.HasPrefix(datosLimpios, `"`) && strings.HasSuffix(datosLimpios, `"`) {
		datosLimpios = datosLimpios[1 : len(datosLimpios)-1]
	}

	// 🔹 3. Intentar convertirlo a JSON
	if err := json.Unmarshal([]byte(datosLimpios), &jsonData); err != nil {
		return nil, fmt.Errorf("error al deserializar JSON: %w", err)
	}

	return jsonData, nil
}
