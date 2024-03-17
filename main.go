package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const (
    numChairs   = 5
    closingTime = 30 * time.Second
)

type BarberShop struct {
    waitingRoom   chan struct{}
    maxCustomers  int
    customersLeft int
    wg            sync.WaitGroup
}

func NewBarberShop(numBarbers, maxCustomers int) *BarberShop {
    bs := &BarberShop{
        waitingRoom:   make(chan struct{}, numChairs),
        maxCustomers:  maxCustomers,
        customersLeft: maxCustomers,
    }

    for i := 0; i < numBarbers; i++ {
        bs.wg.Add(1)
        go bs.barber()
    }

    return bs
}

func (bs *BarberShop) barber() {
    defer bs.wg.Done()
    for {
        select {
        case <-time.After(time.Duration(rand.Intn(3)) * time.Second):
            select {
            case <-bs.waitingRoom:
                fmt.Println("Barber is cutting hair")
                bs.customersLeft--
            default:
                fmt.Println("Barber is sleeping")
            }
        }

        if bs.customersLeft == 0 {
            break
        }
    }
}

func (bs *BarberShop) openShop() {
    timer := time.NewTimer(closingTime)
    <-timer.C
    bs.wg.Wait()
    fmt.Println("All barbers are done for the day. Shop is closed.")
}

func (bs *BarberShop) run() {
    go bs.openShop()

    for i := 0; i < bs.maxCustomers; i++ {
        time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
        select {
        case bs.waitingRoom <- struct{}{}:
            fmt.Println("Customer enters the shop and waits")
        default:
            fmt.Println("Customer enters the shop and leaves because it's full")
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    numBarbers := 2
    maxCustomers := 10
    bs := NewBarberShop(numBarbers, maxCustomers)
    bs.run()
}
