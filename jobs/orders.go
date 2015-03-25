package jobs

import (
    "log"

    "github.com/hashtock/hashtock-go/models"
)

// Trigger execution of bank orders
// TODO(access): This endpoind should only by available to admin or cron
func ExecuteBankOrders() {
    activeOrders, err := models.GetAllActiveBankOrders(nil)
    if err != nil {
        log.Println("job:ExecuteBankOrders: Could not fetch active bank orders. Err:", err)
        return
    }

    //TODO(error): Handle errors somehow
    for _, order := range activeOrders {
        if err := models.ExecuteBankOrder(nil, order); err != nil {
            log.Printf("job:ExecuteBankOrders: Could not execute bank order %v. Err: %v", order.UUID, err)
        }
    }

    if len(activeOrders) > 0 {
        log.Printf("job:ExecuteBankOrders: %v bank orders executed", len(activeOrders))
    } else {
        log.Println("job:ExecuteBankOrders: No bank orders to execute")
    }
}
