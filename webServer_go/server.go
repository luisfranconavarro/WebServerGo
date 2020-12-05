package main

import (
	
	"fmt"
	"strconv"
	"io/ioutil"
	"net/http"
)

var Server1 Server
type Server struct{
	Materias map[string]map[string]float64
	Alumnos map[string]map[string]float64
	Iniciar bool
}

func (this* Server) Constructor(s string){
	if this.Iniciar == false{
		this.Materias = make(map[string]map[string]float64)
		this.Alumnos = make(map[string]map[string]float64)
		this.Iniicar = true
	}
}

func (this *Server) AgregarCalificacion(s []string){
	v, err := this.Materias[s[0]]

	if err == false{
		alumno := make(map[string]float64)
		f2, _ := strconv.ParseFloat(s[2], 8)
		alumno[s[1]] = f2
		this.Materias[s[0]] = alumno
	} else {
		_, err2 := v[s[1]]
		if err2 == false{
			alumno := make(map[string]float64)
			f2, _ := strconv.ParseFloat(s[2], 8)

			for auxAlumno, calificacion := range this.Materias[s[0]] {
				alumno[auxAlumno] = calificacion
			}

			alumno[s[1]] = f2	
			this.Materias[s[0]] = alumno	
		}
	}
	v2, err2 := this.Alumnos[s[1]]
	
	if err2 == false{
		clase := make(map[string]float64)
		
		f2, _ := strconv.ParseFloat(s[2], 8)
		clase[s[0]] = f2
		this.Alumnos[s[1]] = clase
	} else {
		_, err4 := v2[s[0]]
		if err4 == false{
			clase := make(map[string]float64)

			f2, _ := strconv.ParseFloat(s[2], 8)

			for auxClase, calificacion := range this.Alumnos[s[1]] {
				clase[auxClase] = calificacion
			}

			clase[s[0]] = f2
			this.Alumnos[s[1]] = clase
			
		}
	}
}

func (this *Server) PromedioAlumno(nombre string) float64{
	var promedio float64
	var i int64
	promedio = 0
	i = 0

	for _, calificacion := range this.Alumnos[nombre] {
		promedio += calificacion
		i ++
	}
	promedio /= float64(i)
	return promedio
}

func (this *Server) PromedioMateria(nombre string) float64{
	var promedio float64
	var i int64
	promedio = 0
	i = 0

	for _, calificacion := range this.Materias[nombre] {
		promedio += calificacion
		i ++
	}
	
	promedio /= float64(i)
	return promedio
}

func (this *Server) PromedioGeneral() float64{
	var promedio, promedioGeneral float64
	var i int64
	i = 0
	promedio, promedioGeneral = 0, 0

	for nombreAlumno := range this.Alumnos {
		promedio = 0
		for _, calificacion := range this.Alumnos[nombreAlumno] {
			promedio += float64(calificacion)
			i ++
		}
		
		promedioGeneral += promedio
		
	}
	promedioGeneral /= float64(i)
	if promedioGeneral > 0 {
		return promedioGeneral
	}
	return 0
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

func calificacion(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("agregar.html"),
	)
}

func agregarCal(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		var s []string

		s = append(s,req.FormValue("materia"))
		s = append(s,req.FormValue("alumno"))
		s = append(s,req.FormValue("calificacion"))
		Server1.AgregarCalificacion(s)

		fmt.Fprintf(
			res,
			cargarHtml("cali_agregada.html"),
		)
	}
}

func alumno(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("buscar_alumno.html"),
	)
}

func buscarAl(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("prom_alumno.html"),
			Server1.PromedioAlumno(req.FormValue("alumno")),
		)
	}
}

func materia(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("buscar_materia.html"),
	)
}

func buscarMa(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("prom_materia.html"),
			Server1.PromedioMateria(req.FormValue("materia")),
		)
	}
}

func promedioGral(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("general.html"),
	)
}

func mostrarPromGral(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("prom_general.html"),
			Server1.PromedioGeneral(),
		)
	}
}

func main() {
	Server1.Constructor("")
	http.HandleFunc("/agregar", calificacion)
	http.HandleFunc("/agregarCal", agregarCal)

	http.HandleFunc("/buscar_alumno", alumno)
	http.HandleFunc("/buscarAl", buscarAl)

	http.HandleFunc("/buscar_materia", materia)
	http.HandleFunc("/buscarMa", buscarMa)

	http.HandleFunc("/general", promedioGral)
	http.HandleFunc("/mostrarPromGral", mostrarPromGral)

	fmt.Println("Servidor encendido")
	http.ListenAndServe(":9999", nil)
}