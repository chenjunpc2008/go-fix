package ibfixdefine

// 150 ExecType
const (
    // ExecType_New 0 = New
    ExecType_New = "0"

    // ExecType_PartialFill 1 = Partial Fill
    ExecType_PartialFill = "1"

    // ExecType_Fill 2 = Fill
    ExecType_Fill = "2"

    // ExecType_Canceled 4 = Canceled
    ExecType_Canceled = "4"

    // ExecType_Replace 5 = Replace
    ExecType_Replace = "5"

    // ExecType_PendingCancel 6 = Pending Cancel/Replace
    ExecType_PendingCancel = "6"

    // ExecType_PendingNew A = Pending New
    ExecType_PendingNew = "A"

    // ExecType_Rejected 8 = Rejected
    ExecType_Rejected = "8"

    // ExecType_Expired C = Expired
    ExecType_Expired = "C"

    // ExecType_Restated D = Restated
    ExecType_Restated = "D"

    // ExecType_PendingReplace E = Pending Replace
    ExecType_PendingReplace = "E"
)

// 39 OrdStatus
const (
    // OrdStatus_New 0 = New
    OrdStatus_New = "0"

    // OrdStatus_PartiallyFilled 1 = Partially filled
    OrdStatus_PartiallyFilled = "1"

    // OrdStatus_Filled 2 = Filled
    OrdStatus_Filled = "2"

    // OrdStatus_Canceled 4 = Canceled
    OrdStatus_Canceled = "4"

    // OrdStatus_Replaced5 = Replaced
    OrdStatus_Replaced = "5"

    // OrdStatus_PendingCancel 6 = Pending Cancel
    OrdStatus_PendingCancel = "6"

    // OrdStatus_Rejected 8 = Rejected
    OrdStatus_Rejected = "8"

    // OrdStatus_PendingNew A = Pending New
    OrdStatus_PendingNew = "A"

    // OrdStatus_Expired C = Expired
    OrdStatus_Expired = "C"

    // OrdStatus_PendingReplace E = Pending Replace (Returned if tag 8=FIX.4.3)
    OrdStatus_PendingReplace = "E"
)

// 40 OrdType
const (
    // OrderType_Market 1 = Market
    OrderType_Market = "1"

    // OrderType_Limit 2 = Limit
    OrderType_Limit = "2"

    // OrderType_Stop 3 = Stop
    OrderType_Stop = "3"

    // OrderType_StopLimit 4 = Stop limit
    OrderType_StopLimit = "4"

    // OrderType_MarketOnClose 5 = Market on close
    OrderType_MarketOnClose = "5"

    // OrderType_LimitOnClose B = Limit on close
    OrderType_LimitOnClose = "B"

    // OrderType_MarketIfTouchedJ = Market If Touched
    OrderType_MarketIfTouched = "J"

    // OrderType_LimitIfTouched LT = Limit if touched
    OrderType_LimitIfTouched = "LT"

    // OrderType_Pegged P = Pegged (requires execInst = L, R, M, P or O)
    OrderType_Pegged = "P"

    // OrderType_TSLTSL = Trailing Stop Limit
    OrderType_TSL = "TSL"

    // OrderType_MIDPX MIDPX = MidPrice
    OrderType_MIDPX = "MIDPX"
)

// 15 Currency
const (
    // Currency_USD USD = US Dollar
    Currency_USD = "USD"

    // Currency_AUD AUD = Australian Dollar
    Currency_AUD = "AUD"

    // Currency_CAD CAD = Canadian Dollar
    Currency_CAD = "CAD"

    // Currency_CHF CHF = Swiss Franc
    Currency_CHF = "CHF"

    // Currency_EUR EUR = Euros
    Currency_EUR = "EUR"

    // Currency_GBP GBP = British Pound
    Currency_GBP = "GBP"

    // Currency_HKD HKD = Hong Kong $
    Currency_HKD = "HKD"

    // Currency_JPY JPY = Japanese Yen
    Currency_JPY = "JPY"
)

// 20 ExecTransType
const (
    // ExecTransTyp_eNew 0 = New
    ExecTransType_New = "0"

    // ExecTransType_Cancel 1 = Cancel
    ExecTransType_Cancel = "1"

    // ExecTransType_Correct 2 = Correct
    ExecTransType_Correct = "2"

    // ExecTransType_Status 3 = Status
    ExecTransType_Status = "3"
)

// 47 Rule80A
const (
    // Rule80A_A A = Agency single order
    Rule80A_A = "A"

    // Rule80A_J J = Program Order, index arb, for individual customer
    Rule80A_J = "J"

    // Rule80A_K K = Program Order, non-index arb, for individual customer
    Rule80A_K = "K"

    // Rule80A_I I = Individual Investor, single order
    Rule80A_I = "I"

    // Rule80A_P P = Principal
    Rule80A_P = "P"

    // Rule80A_U U = Program Order, index arb for other agency
    Rule80A_U = "U"

    // Rule80A_Y Y = Program Order, non-index arb, for other agency
    Rule80A_Y = "Y"

    // Rule80A_M M = Program Order, index arb for other member
    Rule80A_M = "M"

    // N = Program Order, non-index arb for other member
    Rule80A_N = 9

    // W = All other orders as agent for other member
    Rule80A_W = 10
)

// 54 Side
const (
    // Side_Buy 1 = Buy
    Side_Buy = "1"

    // Side_Sell 2 = Sell
    Side_Sell = "2"

    // Side_BuyMinus 3 = BuyMinus (interpreted as buy)
    Side_BuyMinus = "3"

    // Side_SellPlus 4 = SellPlus (interpreted as sell)
    Side_SellPlus = "4"

    // Side_SellShort 5 = Sell Short
    Side_SellShort = "5"

    // Side_SellShortExempt 6 = Sell Short Exempt
    Side_SellShortExempt = "6"
)

// 59 TimeInForce
const (
    // TimeInForce_Day 0 = Day
    TimeInForce_Day = "0"

    // TimeInForce_GTC 1 = GTC
    TimeInForce_GTC = "1"

    // TimeInForce_OPG 2 = OPG
    TimeInForce_OPG = "2"

    // TimeInForce_IOC 3 = IOC
    TimeInForce_IOC = "3"

    // TimeInForce_FOK 4 = Fill or Kill (FOK)
    TimeInForce_FOK = "4"

    // TimeInForce_GTD 6 = GTD (If used, EITHER tag 432 or tag 126 can be used. NOT both.)
    TimeInForce_GTD = "6"

    // TimeInForce_AtTheClosing 7 = At the Closing
    TimeInForce_AtTheClosing = "7"

    // TimeInForce_Auction 8 = Auction
    TimeInForce_Auction = "8"
)

// 61 Urgency
const (
    // Urgency_Normal 0 = Normal
    Urgency_Normal = "0"

    // Urgency_Flash 1 = Flash
    Urgency_Flash = "1"

    // Urgency_Background 2 = Background
    Urgency_Background = "2"
)

// 63 SettlmntTyp
const (
    // SettlmntTyp 0 = Regular
    SettlmntTyp_Regular = "Regular"
)

// 77 OpenClose
const (
    // OpenClose_Open O=Open
    OpenClose_Open = "O"

    // OpenClose_Close C=Close
    OpenClose_Close = "C"
)

// 87 AllocStatus
const (
    // AllocStatus_Accepted 0 = Accepted (Successfully Processed)
    AllocStatus_Accepted = 0

    // AllocStatus_Rejected 1 = Rejected
    AllocStatus_Rejected = 1

    // AllocStatus_Received 3 = Received
    AllocStatus_Received = 3
)

// 88 RejectCode AllocRejCode
const (
    // RejectCode_UnknownAccount 0 = unknown account(s)
    RejectCode_UnknownAccount = 0

    // RejectCode_IncorrectQuantity 1 = incorrect quantity
    RejectCode_IncorrectQuantity = 1

    // RejectCode_IncorrectAveragePrice 2 = incorrect average price
    RejectCode_IncorrectAveragePrice = 2

    // RejectCode_UnknownExecutingBrokerMnemonic 3 = unknown executing broker mnemonic
    RejectCode_UnknownExecutingBrokerMnemonic = 3

    // RejectCode_CommissionDifference 4 = commission difference
    RejectCode_CommissionDifference = 4

    // RejectCode_UnknownOrderID 5 = unknown OrderID
    RejectCode_UnknownOrderID = 5

    // RejectCode_UnknownListID 6 = unknown ListID
    RejectCode_UnknownListID = 6

    // 7 = other
    RejectCode_Other = 7
)

// 167 SecurityType
const (
    // SecurityType_CS CS = Common stock, can also be sent as STK.
    SecurityType_CS = "CS"

    // SecurityType_FUT FUT = Future
    SecurityType_FUT = "FUT"

    // SecurityType_OPT OPT = Option, except for options on futures.
    SecurityType_OPT = "OPT"

    // SecurityType_FOP FOP = Options on Futures
    SecurityType_FOP = "FOP"

    // SecurityType_WAR WAR = Warrant
    SecurityType_WAR = "WAR"

    // SecurityType_MLEG MLEG = Multi-leg component, can also be sent as MULTILEG.
    SecurityType_MLEG = "MLEG"

    // SecurityType_ CASH = Foreign exchange.
    SecurityType_CASH = "CASH"

    // SecurityType_BOND BOND = Bonds
    SecurityType_BOND = "BOND"

    // SecurityType_CFD CFD = Contract for difference
    SecurityType_CFD = "CFD"

    // SecurityType_CMDTY CMDTY = Spot
    SecurityType_CMDTY = "CMDTY"

    // SecurityType_FUND FUND = Mutual fund
    SecurityType_FUND = "FUND"
)

// 201 PutOrCall
const (
    // PutOrCall_Put 0 = Put
    PutOrCall_Put = 0

    // PutOrCall_Call 1 = Call
    PutOrCall_Call = 1
)

// 204 CustomerOrFirm
const (
    // CustomerOrFirm_Customer 0 = Customer
    CustomerOrFirm_Customer = 1

    // CustomerOrFirm_Firm 1 = Firm
    CustomerOrFirm_Firm = 2
)

// 102 CxlRejReason
const (
    // CxlRejReason_TooLateToCancel 0 = Too late to cancel
    CxlRejReason_TooLateToCancel = 0

    // CxlRejReason_UnknownOrder 1 = Unknown order
    CxlRejReason_UnknownOrder = 1

    // CxlRejReason_BrokerOption 2 = Broker Option
    CxlRejReason_BrokerOption = 2

    // CxlRejReason_OrderAlreadyInPending 3 = Order already in Pending Cancel or Pending Replace status
    CxlRejReason_OrderAlreadyInPending = 3
)

// 103 OrdRejReason
const (
    // OrdRejReason_BrokerOption 0 = Broker Option(default if omitted)
    OrdRejReason_BrokerOption = 0

    // OrdRejReason_OrderExceedsLimit 3 = Order exceeds limit
    OrdRejReason_OrderExceedsLimit = 3
)

// 434 CxlRejResponseTo
const (
    // CxlRejResponseTo_OrderCancelRequest 1 = Order Cancel Request
    CxlRejResponseTo_OrderCancelRequest = 1

    // CxlRejResponseTo_OrderReplaceRequest 2 = Order Cancel/Replace Request
    CxlRejResponseTo_OrderReplaceRequest = 2
)

// 388 DiscretionInst
const (
    // DiscretionInst_RelatedToDisplayedPrice 0 = Related to displayed price
    DiscretionInst_RelatedToDisplayedPrice = "0"
)

// 18 ExecInst
const (
    // ExecInst_MarketPeg P = Market Peg
    ExecInst_MarketPeg = "P"

    // ExecInst_PrimaryPeg R = Primary peg (primary market – buy at  bid/sell at        // offer)
    ExecInst_PrimaryPeg = "R"

    // ExecInst_MidpointPeg M = Midpoint peg
    ExecInst_MidpointPeg = "M"

    // ExecInst_AllOrNone G = All or None
    ExecInst_AllOrNone = "G"

    // ExecInst_PegToVWAP W = Peg to VWAP
    ExecInst_PegToVWAP = "W"

    // ExecInst_TrailingStopPeg a = Trailing Stop Peg(FIX 4.4)
    ExecInst_TrailingStopPeg = "a"

    // ExecInst_PegToStock s = Peg to stock (IBKR Custom value)
    ExecInst_PegToStock = "s"

    // ExecInst_AlgoOrders e = Used for IBKR Algo orders
    ExecInst_AlgoOrders = "e"
)

// 71 AllocTransType
const (
    // AllocTransType_New 0 = New
    AllocTransType_New = "0"

    // AllocTransType_Replace 1 = Replace
    AllocTransType_Replace = "1"

    // AllocTransType_Cancel 2 = Cancel
    AllocTransType_Cancel = "2"

    // AllocTransType_Preliminary 3 = Preliminary (without MiscFees and NetMoney)
    AllocTransType_Preliminary = "3"

    // AllocTransType_Calculated 4 = Calculated (includes MiscFees and NetMoney)
    AllocTransType_Calculated = "4"

    // AllocTransType_CalculatedWithoutPreliminary 5 = Calculated without Preliminary (sent unsolicited by broker, includes MiscFees and NetMoney)
    AllocTransType_CalculatedWithoutPreliminary = "5"
)

// 6086 ShortSaleRule
const (
    ShortSaleRule_1 = 1
    ShortSaleRule_2 = 2
)
