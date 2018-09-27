package services

import (
    "log"
    "io/ioutil"
    "net/http"
    "time"
    "regexp"
    "errors"
    "strings"
    "strconv"
    "github.com/tarm/serial"
    "github.com/influxdata/influxdb/client/v2"
    "github.com/oskarszura/smarthome/utils"
)

const (
    pkgPattern = "<[0-9]+\\.[0-9]+\\|-?[0-9]+\\|[0-1]\\|[0-9]+\\.[0-9]+\\>"
)

type Agent struct {
    Id          string
    Name        string
    Url         string
}

var (
    Agents              []Agent
    tmpNotifyTime       time.Time
    motionNotifyTime    time.Time
    gasNotifyTime       time.Time
)

func getPackageData(stream string) (string, error) {
    pkgRegExp, _ := regexp.Compile(pkgPattern)
    dataPackage := pkgRegExp.FindString(stream)

    if dataPackage == "" {
        return "", errors.New("Data stream doesn't contain valid package (" + stream + ")")
    }

    return strings.Split(strings.Replace(dataPackage, "<", "", -1), ">")[0], nil
}

func getTemperature(data string) string {
    return strings.Split(data, "|")[0]
}

func getMotion(data string) string {
    return strings.Split(data, "|")[1]
}

func getGas(data string) string {
    return strings.Split(data, "|")[2]
}

func getSound(data string) string {
    return strings.Split(data, "|")[3]
}

func writePackage(port *serial.Port) {
    _, err := port.Write([]byte("CMD001"))
    if err != nil {
        log.Println("services: ", err)
    }
}

func addAgent(id string, name string, url string) {
    log.Println("services: adding home agent '" + name + "'")

    agent := Agent{
        Id: id,
        Name: name,
        Url: url,
    }

    Agents = append(Agents, agent)
}

func (a Agent) fetchPackage() {
    response, err := http.Get(a.Url)
    defer response.Body.Close()

    contents, err := ioutil.ReadAll(response.Body)

    if err != nil {
        log.Println("services:  agent '" + a.Name + "'", err)
        return
    }

    unwrappedData, err := getPackageData(string(contents))

    temperature := getTemperature(unwrappedData)
    motion := getMotion(unwrappedData)
    gas := getGas(unwrappedData)
    sound := getSound(unwrappedData)

    if utils.IsAlerts == true {
        if t, err := strconv.ParseFloat(temperature, 32); err == nil {
            if t > 30 {
                now := time.Now()

                if now.Sub(tmpNotifyTime).Hours() >= 1 {
                    tmpNotifyTime = now
                    SendEmail("[" + now.UTC().String() + "] temperature = " + temperature)
                }
            }
        }

        if motion != "0" {
            now := time.Now()

            if now.Sub(motionNotifyTime).Hours() >= 1 {
                motionNotifyTime = now
                SendEmail("[" + now.UTC().String() + "] motion detected")
            }
        }

        if gas != "0" {
            now := time.Now()

            if now.Sub(gasNotifyTime).Hours() >= 1 {
                gasNotifyTime = now
                SendEmail("[" + now.UTC().String() + "] gas alert")
            }
        }
    }

    pt, _ := client.NewPoint(
        a.Name,
        map[string]string{ "home": a.Name },
        map[string]interface{}{
            "temperature": temperature,
            "presence": motion,
            "gas": gas,
            "sound": sound,
            "agent": a.Name,
        },
        time.Now(),
    )

    InfluxBp.AddPoint(pt)
    err = InfluxClient.Write(InfluxBp)
}

func setupAgents() {
    agentsCnf, err := ioutil.ReadFile("hardware/agents.config")

    if err != nil {
        log.Print("services", err)
    }

    agentsConf := strings.Split(string(agentsCnf), "\n")

    for _, c := range agentsConf {
        agentConf := strings.Split(c, ":")

        if (len(agentConf) == 3) {
            id := agentConf[0]
            name := agentConf[1]
            ip := agentConf[2]
            apiUrl := "http://" + ip + "/api"

            addAgent(id, name, apiUrl)
        }
    }
}

func runCommunicationLoop() {
    for range time.Tick(time.Second * 10) {
        if InfluxConnected == false {
            log.Println("services: cannot fetch packages, Influx is down")
            return
        }

        for i := 0; i < len(Agents); i++ {
            a := Agents[i]
            a.fetchPackage()
        }
    }
}

func RunHomeService() {
    setupAgents()
    runCommunicationLoop()
}


