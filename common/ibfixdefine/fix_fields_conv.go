package ibfixdefine

import "fmt"

/*
values to message to fix format values
*/

const (
    /*Handlinst 21 HandlInst
      2 = Automated execution
      order,
      public, Broker
      intervention
      OK.
    */
    Handlinst = "2"

    // SecurityExchange_ISLAND ISLAND -- Nasdaq
    SecurityExchange_ISLAND = "ISLAND"
)

/*
FixBool to bool
*/
func FixBool(b bool) string {
    if b {
        return "Y"
    }
    return "N"
}

/*
FixCurrency 15 Currency
USD = US Dollar
AUD = Australian Dollar
CAD = Canadian Dollar
CHF = Swiss Franc
EUR = Euros
GBP = British Pound
HKD = Hong Kong $
JPY = Japanese Yen
*/
func FixCurrency(src string) (currency string, err error) {
    switch src {
    case "USD", "AUD", "CAD", "CHF", "EUR", "GBP", "HKD", "JPY":
        currency = src

    default:
        return "", fmt.Errorf("illegal Currency val:%s", src)
    }

    return currency, nil
}

/*
FixSide 54 side
1 = Buy
2 = Sell
3 = BuyMinus (interpreted as buy)
4 = SellPlus (interpreted as sell)
5 = Sell Short
6 = Sell Short Exempt
*/
func FixSide(src string) (side string, err error) {
    switch src {
    case "1", "2", "3", "4", "5", "6":
        side = src

    default:
        return "", fmt.Errorf("illegal Side val:%s", src)
    }

    return side, nil
}

/*
FixOrdType 40 OrdType
1 = Market
2 = Limit
3 = Stop
4 = Stop limit
5 = Market on close
B = Limit on close
J = Market If Touched
LT = Limit if touched
P = Pegged (requires execInst = L, R, M, P or O)
TSL = Trailing Stop Limit
MIDPX = MidPrice
*/
func FixOrdType(src string) (ordType string, err error) {
    switch src {
    case "1", "2", "3", "4", "5", "B", "J", "LT", "P", "TSL", "MIDPX":
        ordType = src

    default:
        return "", fmt.Errorf("illegal OrdType val:%s", src)
    }

    return ordType, nil
}

/*
FixExecInst 18 ExecInst
P = Market Peg
R = Primary peg (primary market – buy at bid/sell at offer)
M = Midpoint peg
G = All or None
W = Peg to VWAP
a = Trailing Stop Peg
(FIX 4.4)
s = Peg to stock (IBKR Custom
value)
e = Used for IBKR Algo orders
*/
func FixExecInst(src string) (execInst string, err error) {
    switch src {
    case "P", "R", "M", "G", "W", "a", "s", "e":
        execInst = src

    default:
        return "", fmt.Errorf("illegal ExecInst val:%s", src)
    }

    return execInst, nil
}

/*
FixTimeInForce 59 TimeInForce
0 = Day
1 = GTC
2 = OPG
3 = IOC
4 = Fill or Kill (FOK)
6 = GTD (If used, EITHER tag 432 or tag 126 can be used. NOT both.)
7 = At the Closing
8 = Auction
*/
func FixTimeInForce(src string) (timeInForce string, err error) {
    switch src {
    case "0", "1", "2", "3", "4", "6", "7", "8":
        timeInForce = src

    default:
        return "", fmt.Errorf("illegal TimeInForce val:%s", src)
    }

    return timeInForce, nil
}

/*
FixCustomerOrFirm 204 CustomerOrFirm
0 = Customer
1 = Firm
*/
func FixCustomerOrFirm(c int32) (execInst string, err error) {
    switch c {
    case 0:
        execInst = "0"

    case 1:
        execInst = "1"

    default:
        return "", fmt.Errorf("illegal CustomerOrFirm enum:%d", c)
    }

    return execInst, nil
}

/*
FixOpenClose 77 OpenClose
O=Open
C=Close
*/
func FixOpenClose(src string) (openClose string, err error) {
    switch src {
    case "O", "C":
        openClose = src

    default:
        return "", fmt.Errorf("illegal OpenClose val:%s", src)
    }

    return openClose, err
}

/*
FixSecurityType 167 SecurityType
CS = Common stock, can also
be sent as STK.
FUT = Future
OPT = Option, except for
options on futures.
FOP = Options on Futures
WAR = Warrant
MLEG = Multi-leg component,
can also be sent as
MULTILEG.
CASH = Foreign exchange.
BOND = Bonds
CFD = Contract for difference
CMDTY = Spot
FUND = Mutual fund
*/
func FixSecurityType(src string) (securityType string, err error) {
    switch src {
    case "CS", "FUT", "OPT", "FOP", "WAR", "MLEG", "CASH", "BOND", "CFD", "CMDTY", "FUND":
        securityType = src

    default:
        return "", fmt.Errorf("illegal SecurityType val:%s", src)
    }

    return securityType, nil
}

/*
FixPutOrCall 201 PutOrCall
0 = Put
1 = Call
*/
func FixPutOrCall(p int32) (putOrCall string, err error) {
    switch p {
    case 0:
        putOrCall = "0"

    case 1:
        putOrCall = "1"

    default:
        return "", fmt.Errorf("illegal PutOrCall enum:%d", p)
    }

    return putOrCall, nil
}

/*
FixDiscretionInst 388 DiscretionInst
0 = Related to displayed price
*/
func FixDiscretionInst(src string) (discretionInst string, err error) {
    switch src {
    case "0":
        discretionInst = src

    default:
        return "", fmt.Errorf("illegal DiscretionInst val:%s", src)
    }

    return discretionInst, nil
}

/*
FixShortSaleRule 6086 ShortSaleRule
The value “1” may only be used by an IBKR
customer who uses IBKR as its executing broker
only and does not use IBKR as its clearing broker
(a “Non-Clearing Customer”). The value “1” may
only be used by a Non-Clearing Customer who
has arranged with its clearing broker or custody
agent (as designated in Tag 439 on the order) to
borrow the shares that are required to be
delivered in connection with the short sale order.

The value “2” may only be used by a Non-Clearing
Customer who has arranged with a clearing
broker or custody agent other than its usual
clearing broker or custody agent to borrow the
shares that are required to be delivered in
connection with the short sale order, i.e. a party
other than the clearing firm designated in Tag 439
on the order. If a Non-Clearing customer uses the
value “2” for this Tag 6086, the Non-Clearing
Customer must also use in combination Tags 114
(with the value “N”) and 5700 (with the clearing
broker’s or custody agent’s MPID).

An IBKR customer who executes and clears
through IBKR (a “Cleared Customer”) may not use
this Tag 6086.
*/
func FixShortSaleRule(s int32) (shortSaleRule string, err error) {
    switch s {
    case 1:
        shortSaleRule = "1"

    case 2:
        shortSaleRule = "2"

    default:
        return "", fmt.Errorf("illegal ShortSaleRule enum:%d", s)
    }

    return shortSaleRule, nil
}

/*
FixSettlemntTyp 63 SettlmntTyp
0 = Regular
*/
func FixSettlemntTyp(src string) (settlemntTyp string, err error) {
    switch src {
    case "0":
        settlemntTyp = src

    default:
        return "", fmt.Errorf("illegal SettlemntTyp val:%s", src)
    }

    return settlemntTyp, nil
}

/*
FixRule80A 47 Rule80A (OrderCapacity)
A = Agency single order
J = Program Order, index arb, for individual  customer
K = Program Order, non-index arb, for  individual customer
I = Individual Investor,  single order
P = Principal
U = Program Order, index arb for other agency
Y = Program Order, non-index arb, for other agency
M = Program Order, index arb for other member
N = Program Order, non-index arb for other member
W = All other orders as agent for other member
*/
func FixRule80A(src string) (rule80A string, err error) {
    switch src {
    case "A", "J", "K", "I", "P", "U", "Y", "M", "N", "W":
        rule80A = src

    default:
        return "", fmt.Errorf("illegal Rule80A val:%s", src)
    }

    return rule80A, nil
}
