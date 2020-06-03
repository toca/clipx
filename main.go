package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "github.com/atotto/clipboard"
	"clipx/models"

)


func main(){
	// clipboard 監視
	written := make(chan bool, 16)
	quit := make(chan bool, 1)
	monitor := models.NewMonitor(written, quit)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		err := monitor.Stop()
		fmt.Println("interrupted!!")
		if err != nil {
			fmt.Printf("Monitor.Stop failed: %v\n", err)
		}
	}()

	clipboard := models.NewClipboard()
	monitorErr := make(chan error, 1)
	go func() {
		err := monitor.Monitoring()
		monitorErr <- err
	}()
	loop:
		for {
			select{
			case <- written:
				fmt.Println("written")
				if res, err := clipboard.IsStringable(); err == nil{
					fmt.Printf("IsStringable:%v\n", res)
				}
				// str, err := clipboard.ReadAll()
				// if err == nil {
				// 	fmt.Println(str)
				// } else {
				// 	fmt.Println(err)
				// }
			case <- quit:
				fmt.Println("quit")
				break loop
			case err := <- monitorErr:
				fmt.Printf("MonitoringError: %v\n", err)
			}
		}

	// key hook
	// queue
	// ui
	// paste
	fmt.Println("finish!")
	// str, err := clipboard.ReadAll()
	// if err == nil {
	// 	fmt.Println(str)
	// } else {
	// 	fmt.Println(err)
	// }
}