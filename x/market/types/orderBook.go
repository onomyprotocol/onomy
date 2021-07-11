package types

import (
	"strconv"
)

func NewOrderBook(AmountDenom string, ExchRateDenom string) OrderBook {
	return OrderBook{
		OrderIDTrack: 1,
		AmountDenom:  AmountDenom,
		ExchRateDenom: ExchRateDenom,
	}
}

func (book OrderBook) InsertAt(i int, order Order) OrderBook {
	if len(book.Orders) == i {
		book.Orders = append(book.Orders, &order)
		return book
	}

	book.Orders = append(book.Orders[:i+1], book.Orders[i:]...)
	book.Orders[i] = &order

	return book
}

func (book OrderBook) BondAskAmount(bondAskBal, bondBidBal, bidAmount float64) float64 {
	return bondAskBal * bidAmount / bondBidBal
}

func (book OrderBook) Reconcile(bondAsk, bondBid float64) {
	var f func(book OrderBook)
	f = func(book OrderBook) {
		// Iterate over the bookAsk sequence processing bond
		// purchases until bookAsk record with exchrate less than
		// the bond price is reached
		for i := 0; i < len(book.GetOrders()); i++ {
			ask := book.GetOrders()[i]

			// Skip bids
			if ask.GetId()/2 == 0 {
				continue
			}

			exchrate, err := strconv.ParseFloat(ask.ExchRate, 64)
			if err != nil {
				// TODO: what should we do in case of error?
				return
			}

			// Case 1
			// The order at the Head of bookAsk sequence has an
			// exchange rate greater than or equal to the ask bond
			// exchange rate
			if exchrate >= bondAsk/bondBid {
				bondAsk -= float64(ask.GetAmount())
				bondBid += float64(ask.GetAmount())
				book.Orders = book.GetOrders()[:i]
				f(book)

				break
			}

			// Case 2
			// Head of bookAsk exchange rate less than ask bond
			// exchange rate
			if exchrate < bondAsk/bondBid {
				var g func(book OrderBook)
				g = func(book OrderBook) {
					for j := 0; j < len(book.GetOrders()); j++ {
						bid := book.GetOrders()[j]

						// Skip asks
						if bid.GetId()/2 != 0 {
							continue
						}

						// Case 2.1
						// Head of bookBid exchange rate
						// greater than or equal to the
						// updated bid bond exchange rate
						if exchrate >= bondAsk/bondBid {
							bondBid -= float64(bid.GetAmount())
							bondAsk += float64(bid.GetAmount())
							book.Orders = book.GetOrders()[:j]
							g(book)

							break
						}

						// Case 2.2
						// Head of bookBid exchange rate
						// less than the updated bid bond
						// exchange rate
						//
						// Processing Complete
						// Update bonds and books states
						if exchrate < bondBid/bondAsk {
							return
						}
					}
				}
			}
		}
	}
}

func (book OrderBook) ProcessOrder(order Order) OrderBook {
	if len(book.Orders) <= 0 {
		return book
	}

	return book
}

