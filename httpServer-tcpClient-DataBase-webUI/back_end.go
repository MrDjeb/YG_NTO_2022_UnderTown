package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Patient struct {
	Id          int
	Name        string
	Surname     string
	Patronymic  string
	Email       string
	PhoneNumber string
	RoomId      string
}

type Doctor struct {
	Id          int
	Name        string
	Surname     string
	Patronymic  string
	Email       string
	PhoneNumber string
}
type Room struct {
	RoomId       string
	Temperature  string
	Humidity     string
	Illumination string
	Energy       string
	Humans       string
	Door         string
	Heater       string
	Vent         string
}

func cherr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cherrH(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func newHash() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, 32)
	for i := range s {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	h := md5.New()
	io.WriteString(h, string(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

var t_min, t_max, t_start time.Time
var cAlarm chan bool
var isAlarm bool

func gettingAlarm(conn net.Conn, c chan bool) {
	timer := time.NewTicker(time.Millisecond * time.Duration(100))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			put := make([]byte, 4096)
			fmt.Fprintf(conn, "Where is my data?")
			_, err := conn.Read(put)
			cherr(err)

			str := string(put)[:5]
			//log.Println(str)
			fl, _ := strconv.ParseFloat(str, 32)
			//cherr(err)
			in := int(fl * 1000)
			//fmt.Println(in)
			if in == 0 {
				t_min = time.Now()
			} else if in >= 1500 {
				t_max = time.Now()
			}
			log.Println(math.Abs(float64(t_max.Sub(t_min))) / 1000000000)
			if int(math.Abs(float64(t_max.Sub(t_min)))/1000000000) > 10 {
				//log.Println(")DSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS")
				c <- true ///////////////////////
			}
		}

	}
}

func startAlarm(c chan bool) {
	conn, err := net.Dial("tcp", "192.168.1.46:2323")
	cherr(err)
	fmt.Fprintf(conn, "Where is my data?")
	//log.Print(bufio.NewReader(conn).ReadString('!'))
	put := make([]byte, 4096)
	_, err = conn.Read(put)
	cherr(err)
	t_min = time.Now()
	t_max = time.Now()
	gettingAlarm(conn, c)
}

const (
	URL  = "192.168.137.1/"
	PORT = ":39999"
	CORE = "./ui/html/"
)

var (
	dbDoctor  = DBDoctor{DBConnect()}
	dbPatient = DBPatient{DBConnect()}
	dbRoom    = DBRoom{DBConnect()}
	//dbCard = DBCard{DBConnect()}
	//dbLogCard = DBLogCard{DBConnect()}
)

func h_alert(w http.ResponseWriter, r *http.Request, c chan bool) {
	//var result map[string]interface{}
	//json.NewDecoder(r.Body).Decode(&result)
	//log.Println("server: ", "Request: ", result)
	select {
	case rez := <-c:
		//auth := smtp.PlainAuth("", "r.sadovec0403", "coltmvnubhxwidav", "smtp.yandex.ru")
		//cherr(smtp.SendMail("smtp.yandex.ru:25", auth, "from@ya.ru", []string{"to@ya.ru"}, []byte("Текст письма.")))

		w.Header().Set("Content-Type", "application/json")
		ans := map[string]interface{}{
			"IVL":              false,
			"enter":            "nill",
			"imposter_inside":  "nill",
			"imposter_outside": "nill",
		}
		ans["IVL"] = rez
		log.Println("ALLLLLLLLLLLLLLLLLLLLLLLLLEEEEEEEEEEEEEEEEEEEEEEEEEERRRRRRRRRRRRRTTTTTTTTTT:  ", ans)
		data, err := json.Marshal(ans)
		cherr(err)
		w.Write(data)
	default:
		w.Header().Set("Content-Type", "application/json")
		ans := map[string]interface{}{
			"IVL":              false,
			"enter":            "nill",
			"imposter_inside":  "nill",
			"imposter_outside": "nill",
		}
		ans["IVL"] = false
		//log.Println("server: ", "AnswerAlert:  ", ans)
		data, err := json.Marshal(ans)
		cherr(err)
		w.Write(data)
	}
}

func h_esp(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{}
	json.NewDecoder(r.Body).Decode(&result)
	//log.Println(r.Body)
	log.Println("server: ", "RequestEsp: ", result)

	var dataRoom Room
	dataRoom.RoomId, _ = result["rmId"].(string)
	dataRoom.Temperature, _ = result["temp"].(string)
	dataRoom.Humidity, _ = result["hum"].(string)
	dataRoom.Illumination, _ = result["il"].(string)
	dataRoom.Energy = "2.5"
	dataRoom.Humans, _ = result["mv"].(string)
	dataRoom.Door, _ = result["dr"].(string)
	dataRoom.Heater, _ = result["heat"].(string)
	dataRoom.Vent, _ = result["vent"].(string)

	log.Println("INFO ROOM", dataRoom)
	dbRoom.update(dataRoom)
	ans := map[string]interface{}{
		"~Card":  "~~~",
		"~IDesp": "nill",
		"~Info":  "nill",
		"~Pass":  "nill",
	}
	//log.Println("server: ", "Answer:  ", ans)
	data, err := json.Marshal(ans)
	cherr(err)
	w.Write(data)
}

func h_login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles(CORE + "autorization.html")
		cherrH(err, w)
		err = ts.Execute(w, nil)
		cherrH(err, w)
	} else {
		if r.FormValue("username") == "" {
			ts, err := template.ParseFiles(CORE + "autorization_Error.html")
			cherrH(err, w)
			err = ts.Execute(w, nil)
			cherrH(err, w)
		} else {
			okD := dbDoctor.isLog(r.FormValue("username"))
			okP := dbPatient.isLog(r.FormValue("username"))
			if okD && !okP {
				http.Redirect(w, r, "/doctor?id="+r.FormValue("username"), http.StatusSeeOther)
			} else if !okD && okP {
				http.Redirect(w, r, "/patient?id="+r.FormValue("username"), http.StatusSeeOther)
			} else if okD && okP {
				http.Error(w, "DataBase ERROR!!!", 500)
			} else {
				ts, err := template.ParseFiles(CORE + "autorization_Error.html")
				cherrH(err, w)
				err = ts.Execute(w, nil)
				cherrH(err, w)
			}
		}
	}
}

func h_reg(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ts, err := template.ParseFiles(CORE + "registration.html")
		cherrH(err, w)
		err = ts.Execute(w, nil)
		cherrH(err, w)
	} else {
		dataPat := Patient{1, r.FormValue("name"), r.FormValue("surname"), r.FormValue("patronymic"), r.FormValue("email"), r.FormValue("phone"), r.FormValue("roomId")}
		//log.Println(r.FormValue("name"), r.FormValue("surname"), r.FormValue("patronymic"), r.FormValue("email"), r.FormValue("phone"), r.FormValue("roomId"))
		if dataPat.Name == "" || dataPat.Surname == "" || dataPat.Patronymic == "" || dataPat.Email == "" || dataPat.PhoneNumber == "" || dataPat.RoomId == "" {
			//log.Println("null", dbPatient)
			ts, err := template.ParseFiles(CORE + "registration_Error.html")
			cherrH(err, w)
			err = ts.Execute(w, nil)
			cherrH(err, w)
		} else {
			if dbPatient.isLog(dataPat.PhoneNumber) {
				//log.Println("islog", dbPatient)
				ts, err := template.ParseFiles(CORE + "registration_ErrorExist.html")
				cherrH(err, w)
				err = ts.Execute(w, nil)
				cherrH(err, w)
			} else {
				log.Println("ok", dbPatient)
				dbPatient.add(dataPat)
				http.Redirect(w, r, "/patient?id="+dataPat.PhoneNumber, http.StatusSeeOther)
			}
		}
	}
}

func h_doctor(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(CORE + "doctor.html")
	cherr(err)
	idDoc := r.URL.Query().Get("id")
	idPat := r.FormValue("name_pat")
	dataPat := dbPatient.getData(idPat)
	dataRoom := dbRoom.getData(dataPat.RoomId)
	dataPhones := dbPatient.getListId()

	tempDoc := struct {
		In_Name      string
		Name         string
		Surname      string
		Patronymic   string
		Phones       []string
		RoomId       string
		Temperature  string
		Humidity     string
		Illumination string
		Energy       string
		Weight       string
		Humans       string
		Door         string
		Heater       string
		Vent         string
		Phone        string
		Mail         string
		UrlId        string
	}{
		In_Name: dbDoctor.getIn_Name(idDoc), Name: "___", Surname: "___", Patronymic: "___",
		Phones: dataPhones, RoomId: "___", Temperature: "___",
		Humidity: "___", Illumination: "___", Energy: "___", Weight: "___",
		Humans: "___", Door: "___", Heater: "___", Vent: "___", Phone: "___", Mail: "___", UrlId: "?id=" + idDoc}

	tempDoc.Name = dataPat.Name
	tempDoc.Surname = dataPat.Surname
	tempDoc.Patronymic = dataPat.Patronymic
	tempDoc.Mail = dataPat.Email
	tempDoc.Phone = dataPat.PhoneNumber
	tempDoc.RoomId = dataPat.RoomId

	tempDoc.Temperature = dataRoom.Temperature
	tempDoc.Humidity = dataRoom.Humidity
	tempDoc.Illumination = dataRoom.Illumination
	tempDoc.Energy = dataRoom.Energy
	tempDoc.Humans = dataRoom.Humans
	tempDoc.Door = dataRoom.Door
	tempDoc.Heater = dataRoom.Heater
	tempDoc.Vent = dataRoom.Vent
	r.FormValue("patInfo")
	//log.Println("Form: ", r.FormValue("name_pat"))
	//log.Println("Doctor ROOM", dataRoom)
	t.Execute(w, tempDoc)
}

func h_patient(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(CORE + "patient.html")
	cherr(err)
	idPat := r.URL.Query().Get("id")
	dataPat := dbPatient.getData(idPat)
	dataRoom := dbRoom.getData(dataPat.RoomId)
	dataPhones := dbPatient.getListId()

	tempDoc := struct {
		In_Name      string
		Name         string
		Surname      string
		Patronymic   string
		Phones       []string
		RoomId       string
		Temperature  string
		Humidity     string
		Illumination string
		Energy       string
		Weight       string
		Humans       string
		Door         string
		Heater       string
		Vent         string
		Phone        string
		Mail         string
		UrlId        string
	}{
		In_Name: dbDoctor.getIn_Name(idPat), Name: "___", Surname: "___", Patronymic: "___",
		Phones: dataPhones, RoomId: "___", Temperature: "___",
		Humidity: "___", Illumination: "___", Energy: "___", Weight: "___",
		Humans: "___", Door: "___", Heater: "___", Vent: "___", Phone: "___", Mail: "___", UrlId: "?id=" + idPat}

	tempDoc.Name = dataPat.Name
	tempDoc.Surname = dataPat.Surname
	tempDoc.Patronymic = dataPat.Patronymic
	tempDoc.Mail = dataPat.Email
	tempDoc.Phone = dataPat.PhoneNumber
	tempDoc.RoomId = dataPat.RoomId

	tempDoc.Temperature = dataRoom.Temperature
	tempDoc.Humidity = dataRoom.Humidity
	tempDoc.Illumination = dataRoom.Illumination
	tempDoc.Energy = dataRoom.Energy
	tempDoc.Humans = dataRoom.Humans
	tempDoc.Door = dataRoom.Door
	tempDoc.Heater = dataRoom.Heater
	tempDoc.Vent = dataRoom.Vent
	t.Execute(w, tempDoc)
}

func main() {
	ch := make(chan bool)
	go startAlarm(ch)

	log.SetFlags(log.Lmicroseconds)

	mux := http.NewServeMux()
	mux.HandleFunc(URL+"esp", h_esp)
	mux.HandleFunc(URL+"sing_in", h_login)
	mux.HandleFunc(URL+"sing_up", h_reg)
	mux.HandleFunc(URL+"alert", func(w http.ResponseWriter, r *http.Request) { h_alert(w, r, ch) })
	mux.HandleFunc(URL+"doctor", h_doctor)
	mux.HandleFunc(URL+"patient", h_patient)

	dbDoctor.delete()
	dbPatient.delete()
	dbRoom.delete()
	dbDoctor.create()
	dbPatient.create()
	dbRoom.create()
	dbDoctor.add(Doctor{1, "Иванов", "Иван", "Иванович", "ооо@kk.chgm ", "7999"})
	dbPatient.add(Patient{1, "Илья", "Дановский", "Валентинович", "ewfhdfnh@njedjd.ujr", "7111", "01"})
	dbRoom.add(Room{"01", "rr2", "rr3", "rr4", "rr5", "rr6", "rr7", "rr8", "rr9"})

	//DBPatient.add(Patient{Id: 23, Name: "Fff", Surname: "ff", Patronymic: "gfngh", Email: "gfgh", PhoneNumber: "343434"})

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("server: ", "Server is runing...")
	log.Fatal(http.ListenAndServe(URL[:len(URL)-1]+PORT, mux))
}
