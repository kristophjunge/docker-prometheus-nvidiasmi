package main

import (
    "io"
    "net/http"
    "encoding/xml"
    "os/exec"
    //"strconv"
    "log"
    "os"
    "fmt"
    "strconv"
)

var indexHtml = string(`
<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Nvidia SMI Exporter</title>
    </head>
    <body>
        <h1>Prometheus Nvidia SMI Exporter</h1>
        <p><a href="/metrics">Metrics</a></p>
    </body>
</html>
`)


type NvidiaSmiLog struct {
    DriverVersion string `xml:"driver_version"`
    AttachedGPUs int `xml:"attached_gpus"`
    GPUs []struct {
        ProductName string `xml:"product_name"`
        ProductBrand string `xml:"product_brand"`
        FanSpeed string `xml:"fan_speed"`
        PCI struct {
            PCIBus string `xml:"pci_bus"`
        } `xml:"pci"`
        FbMemoryUsage struct {
            Total string `xml:"total"`
            Used string `xml:"used"`
            Free string `xml:"free"`
        } `xml:"fb_memory_usage"`
        Utilization struct {
            GPUUtil string `xml:"gpu_util"`
            MemoryUtil string `xml:"memory_util"`
        } `xml:"utilization"`
        Temperature struct {
            GPUTemp string `xml:"gpu_temp"`
            GPUTempMaxThreshold string `xml:"gpu_temp_max_threshold"`
            GPUTempSlowThreshold string `xml:"gpu_temp_slow_threshold"`
        } `xml:"temperature"`
        PowerReadings struct {
            PowerDraw string `xml:"power_draw"`
            PowerLimit string `xml:"power_limit"`
        } `xml:"power_readings"`
        Clocks struct {
            GraphicsClock string `xml:"graphics_clock"`
            SmClock string `xml:"sm_clock"`
            MemClock string `xml:"mem_clock"`
            VideoClock string `xml:"video_clock"`
        } `xml:"clocks"`
        MaxClocks struct {
            GraphicsClock string `xml:"graphics_clock"`
            SmClock string `xml:"sm_clock"`
            MemClock string `xml:"mem_clock"`
            VideoClock string `xml:"video_clock"`
        } `xml:"max_clocks"`
    } `xml:"gpu"`
}

func metrics(w http.ResponseWriter, r *http.Request) {

    log.Print("Serving /metrics")

    // Get current dir
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    // Call system command
    app := "/bin/cat"
    arg0 := dir + "/test.xml";
    cmd := exec.Command(app, arg0)
    stdout, err := cmd.Output()
    if err != nil {
        println(err.Error())
        return
    }
    //print(string(stdout))

    // Parse XML
    var xmlData NvidiaSmiLog
    xml.Unmarshal(stdout, &xmlData)
    fmt.Println(xmlData)

    // Output
    io.WriteString(w, "nvidiasmi_driver_version" + " " + xmlData.DriverVersion + "\n")
    io.WriteString(w, "nvidiasmi_attached_gpus" + " " + strconv.Itoa(xmlData.AttachedGPUs) + "\n")
    for _, GPU := range xmlData.GPUs {
        io.WriteString(w, "nvidiasmi_fan_speed{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.FanSpeed + "\n")
        io.WriteString(w, "nvidiasmi_memory_usage_total{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.FbMemoryUsage.Total + "\n")
        io.WriteString(w, "nvidiasmi_memory_usage_used{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.FbMemoryUsage.Used + "\n")
        io.WriteString(w, "nvidiasmi_memory_usage_free{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.FbMemoryUsage.Free + "\n")
        io.WriteString(w, "nvidiasmi_utilization_gpu{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Utilization.GPUUtil + "\n")
        io.WriteString(w, "nvidiasmi_utilization_memory{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Utilization.MemoryUtil + "\n")
        io.WriteString(w, "nvidiasmi_temp_gpu{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Temperature.GPUTemp + "\n")
        io.WriteString(w, "nvidiasmi_temp_gpu_max{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Temperature.GPUTempMaxThreshold + "\n")
        io.WriteString(w, "nvidiasmi_temp_gpu_slow{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Temperature.GPUTempSlowThreshold + "\n")
        io.WriteString(w, "nvidiasmi_power_draw{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.PowerReadings.PowerDraw + "\n")
        io.WriteString(w, "nvidiasmi_power_limit{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.PowerReadings.PowerLimit + "\n")
        io.WriteString(w, "nvidiasmi_clock_graphics{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Clocks.GraphicsClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_graphics_max{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.MaxClocks.GraphicsClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_sm{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Clocks.SmClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_sm_max{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.MaxClocks.SmClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_mem{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Clocks.MemClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_mem_max{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.MaxClocks.MemClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_video{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.Clocks.VideoClock + "\n")
        io.WriteString(w, "nvidiasmi_clock_video_max{gpu=\"" + GPU.PCI.PCIBus + "\"}" + " " + GPU.MaxClocks.VideoClock + "\n")
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /index")

    io.WriteString(w, indexHtml)
}

func main() {
    log.Print("Prometheus Nvidia SMI Exporter running!")

    http.HandleFunc("/", index)
    http.HandleFunc("/metrics", metrics)
    http.ListenAndServe(":8000", nil)
}
