// Cron service
// Run as part of service test suite
package jobs_test

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/hashtock/hashtock-go/jobs"
    "github.com/hashtock/hashtock-go/models"
    "github.com/hashtock/hashtock-go/test_tools"
)

// TODO(tests): Those tests need to be simplified somehow :(
// And we need to test handling errors
func TestExecuteBankBuyOrders(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    user_req := new(http.Request)

    tag := models.HashTag{HashTag: "Tag1", Value: 2.00, InBank: 100}
    app.Put(tag)
    order := models.OrderBase{
        Action:    "buy",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  3,
    }
    models.PlaceOrder(user_req, app.User, order)

    // Initial state check
    profile, _ := models.GetOrRegisterProfile(user_req, app.User.UserID, "")
    assert.NotNil(t, profile)
    activeOrders, _ := models.GetActiveUserOrders(user_req, profile)
    completedOrders, _ := models.GetCompletedUserOrders(user_req, profile, "", "")
    shares, _ := models.GetProfileShares(user_req, profile)
    assert.Len(t, activeOrders, 1)
    assert.Len(t, completedOrders, 0)
    assert.Equal(t, profile.Founds, models.StartingFounds)
    assert.Len(t, shares, 0)

    // Execute bank orders
    jobs.ExecuteBankOrders()

    // New request to get real, not cached data
    user_req = new(http.Request)

    expected_share := models.TagShare{
        HashTag:  "Tag1",
        UserID:   app.User.UserID,
        Quantity: 3,
    }

    // State after
    profileAfter, _ := models.GetOrRegisterProfile(user_req, app.User.UserID, "")
    assert.NotNil(t, profileAfter)
    activeOrdersAfter, _ := models.GetActiveUserOrders(user_req, profile)
    completedOrdersAfter, _ := models.GetCompletedUserOrders(user_req, profile, "", "")
    sharesAfter, _ := models.GetProfileShares(user_req, profile)
    assert.Len(t, activeOrdersAfter, 0)
    assert.Len(t, completedOrdersAfter, 1)
    assert.Equal(t, profileAfter.Founds, models.StartingFounds-2*3)
    assert.Equal(t, sharesAfter, []models.TagShare{expected_share})
}

func TestExecuteBankSellOrders(t *testing.T) {
    app := test_tools.NewTestApp(t)
    defer app.Stop()

    user_req := new(http.Request)

    tag := models.HashTag{HashTag: "Tag1", Value: 2.00, InBank: 10}
    app.Put(tag)
    order := models.OrderBase{
        Action:    "sell",
        BankOrder: true,
        HashTag:   "Tag1",
        Quantity:  3,
    }
    models.PlaceOrder(user_req, app.User, order)
    user_share := models.TagShare{
        HashTag:  "Tag1",
        UserID:   app.User.UserID,
        Quantity: 4,
    }
    app.Put(user_share)

    // Initial state check
    profile, _ := models.GetOrRegisterProfile(user_req, app.User.UserID, "")
    assert.NotNil(t, profile)
    activeOrders, _ := models.GetActiveUserOrders(user_req, profile)
    completedOrders, _ := models.GetCompletedUserOrders(user_req, profile, "", "")
    shares, _ := models.GetProfileShares(user_req, profile)
    assert.Len(t, activeOrders, 1)
    assert.Len(t, completedOrders, 0)
    assert.Equal(t, models.StartingFounds, profile.Founds)
    assert.Equal(t, []models.TagShare{user_share}, shares)

    // Execute bank orders
    jobs.ExecuteBankOrders()

    // New request to get real, not cached data
    user_req = new(http.Request)

    expected_share := models.TagShare{
        HashTag:  "Tag1",
        UserID:   app.User.UserID,
        Quantity: 1,
    }
    expectedHashTagAfter := &models.HashTag{HashTag: "Tag1", Value: 2.00, InBank: 13}

    // State after
    profileAfter, _ := models.GetOrRegisterProfile(user_req, app.User.UserID, "")
    assert.NotNil(t, profileAfter)
    activeOrdersAfter, _ := models.GetActiveUserOrders(user_req, profileAfter)
    completedOrdersAfter, _ := models.GetCompletedUserOrders(user_req, profileAfter, "", "")
    sharesAfter, _ := models.GetProfileShares(user_req, profileAfter)
    hashTagAfter, _ := models.GetHashTag(user_req, "Tag1")
    assert.Len(t, activeOrdersAfter, 0)
    assert.Len(t, completedOrdersAfter, 1)
    assert.Equal(t, models.StartingFounds+2*3, profileAfter.Founds)
    assert.Equal(t, []models.TagShare{expected_share}, sharesAfter)
    assert.Equal(t, expectedHashTagAfter, hashTagAfter)
}
