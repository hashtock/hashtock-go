package bank

import (
    "errors"

    "appengine"
    "appengine/datastore"

    "profiles"
)

// Returns info about given HashTag
func HashTagInfo(ctx appengine.Context, hash string) (entry BankEntry, err error) {
    key := Key(ctx, hash)
    err = datastore.Get(ctx, key, &entry)
    return
}

// Returns info about all known HashTags
func HashTagInfoAll(ctx appengine.Context) (entries []BankEntry, err error) {
    q := datastore.NewQuery("BankEntry").Order("-Value")
    _, err = q.GetAll(ctx, &entries)
    return
}

// Executes selling shares by the bank
func SellToUser(ctx appengine.Context, hash string, amount float64) (err error) {
    var profile profiles.Profile
    var bank_share BankEntry

    if bank_share, err = HashTagInfo(ctx, hash); err != nil {
        return
    }

    if profile, err = profiles.CurrentProfile(ctx); err != nil {
        return
    }

    user_share := profile.GetShare(ctx, hash)

    if bank_share.InBank < amount {
        err = errors.New("Bank does not have enough shares")
        return
    }

    if profile.Founds < amount*bank_share.Value {
        err = errors.New("You does not have enough founds")
        return
    }

    profile.Founds -= amount * bank_share.Value
    user_share.Quantity += amount
    bank_share.InBank -= amount

    if err = profile.Save(ctx); err != nil {
        return
    }

    if err = user_share.Save(ctx); err != nil {
        return
    }

    if err = bank_share.Save(ctx); err != nil {
        return
    }

    return
}

// Executes buying shares by the bank
func BuyFromUser(ctx appengine.Context, hash string, amount float64) (err error) {
    var profile profiles.Profile
    var bank_share BankEntry

    if bank_share, err = HashTagInfo(ctx, hash); err != nil {
        return
    }

    if profile, err = profiles.CurrentProfile(ctx); err != nil {
        return
    }

    user_share := profile.GetShare(ctx, hash)

    if user_share.Quantity < amount {
        err = errors.New("You does not have so many shares")
        return
    }

    profile.Founds += amount * bank_share.Value
    user_share.Quantity -= amount
    bank_share.InBank += amount

    if err = profile.Save(ctx); err != nil {
        return
    }

    if err = user_share.Save(ctx); err != nil {
        return
    }

    if err = bank_share.Save(ctx); err != nil {
        return
    }

    return
}
