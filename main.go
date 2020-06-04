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
		fmt.Println("[interrupted]")
		if err != nil {
			fmt.Printf("Monitor.Stop failed: %v\n", err)
		}
	}()

	cb := models.NewClipboard()
	monitorErr := make(chan error, 1)
	go func() {
		fmt.Println("[begin monitoring]")
		err := monitor.Monitoring()
		monitorErr <- err
	}()
	loop:
		for {
			select{
			case <- written:
				fmt.Println("[written]")
				stringable, err := cb.IsStringable();
				if err != nil {
					fmt.Println(err)
					continue
				}
				if !stringable {
					continue
				}
				str, err := cb.GetAsString()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(str)
				}
			case <- quit:
				fmt.Println("[quit]")
				break loop
			case err := <- monitorErr:
				fmt.Printf("MonitoringError: %v\n", err)
			}
		}

	// key hook
	// queue
	// ui
	// paste
	fmt.Println("[process finished]")
	// str, err := clipboard.ReadAll()
	// if err == nil {
	// 	fmt.Println(str)
	// } else {
	// 	fmt.Println(err)
	// }
}