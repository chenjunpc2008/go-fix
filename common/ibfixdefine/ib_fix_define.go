package ibfixdefine

const (
    // MsgTypeHeartbeat Heartbeat
    MsgTypeHeartbeat = "0"

    // MsgTypeTestRequest Test Request
    MsgTypeTestRequest = "1"

    // MsgTypeResendRequest Resend Request
    MsgTypeResendRequest = "2"

    // MsgTypeSessionLevelReject Session Level Reject
    MsgTypeSessionLevelReject = "3"

    // MsgTypeSequenceReset Sequence Reset and Gap Fill
    MsgTypeSequenceReset = "4"

    // MsgTypeLogout Logout
    MsgTypeLogout = "5"

    // MsgTypeExecutionReport Execution Report
    MsgTypeExecutionReport = "8"

    // MsgTypeOrderCancelReject Order Cancel Reject
    MsgTypeOrderCancelReject = "9"

    // MsgTypeLogon logon
    MsgTypeLogon = "A"

    // MsgTypeBulletinMessage News/Bulletin Message
    MsgTypeBulletinMessage = "B"

    // MsgTypeNewOrderSingle New Order Single
    MsgTypeNewOrderSingle = "D"

    // MsgTypeOrderCancel Order Cancel Request
    MsgTypeOrderCancel = "F"

    // MsgTypeOrderReplace Cancel/Replace Request
    MsgTypeOrderReplace = "G"

    // MsgTypeOrderStatus Order Status Request
    MsgTypeOrderStatus = "H"

    // MsgTypeAllocation Allocation
    MsgTypeAllocation = "J"

    // MsgTypeAllocationAck Allocation Ack
    MsgTypeAllocationAck = "P"

    // MsgTypeMultilegCancelReplaceRequest Multileg Cancel Replace Request
    MsgTypeMultilegCancelReplaceRequest = "AC"
)

const (
    // tags for admin message

    // Tag8 8 BeginString
    Tag8 = 8

    // Tag35 35 MsgType
    Tag35 = 35

    // Tag43  PossDupFlag
    Tag43 = 43

    // Tag49 49 SenderCompID
    Tag49 = 49

    // Tag56 56 TargetCompID Default is “IB.”
    Tag56 = 56

    // Tag34 34 MsgSeqNum
    Tag34 = 34

    // Tag36 36 NewSeqNo
    Tag36 = 36

    // Tag45 45 RefSeqNo
    Tag45 = 45

    // Tag52 52 SendingTime UTC Time Stamp
    // Time of message transmission (always expressed in UTC (Universal Time
    // Coordinated, also known as “GMT”)
    Tag52 = 52

    // Tag58 58 Text
    Tag58 = 58

    // Tag98 98 EncryptMethod
    // Valid Value = 0 (none, encryption is currently not supported)
    Tag98 = 98

    // Tag102 CxlRejReason
    Tag102 = 102

    // Tag108 108 HeartBtIn
    Tag108 = 108

    // Tag112 112 TestReqID Y A timestamp string is suggested for TestReqID
    Tag112 = 112

    // Tag123 123 GapFillFlag
    // Indicates that the Sequence Reset message is
    // replacing administrative or application messages which
    // will not be resent.
    // Y = Gap Fill message, MsgSeqNum field valid
    // N = Sequence Reset, ignore MsgSeqNum
    Tag123 = 123

    // Tag141 141 ResetSeqNumFlag
    Tag141 = 141
)

const (
    // tags for app message

    // Tag1 1 Account
    Tag1 = 1

    // Tag6 AvgPx. Calculated average price of all fills on this order.
    Tag6 = 6

    // Tag7 7 BeginSeqNo
    Tag7 = 7

    // Tag11 ClOrdID Unique identifier for Order as assigned by order routing firm. Uniqueness must be
    //       guaranteed within a single trading day. Firms which electronically submit multi-day orders should
    //       consider embedding a date within the ClOrdID field to assure uniqueness across days.
    Tag11 = 11

    // Tag15  Currency. This tag is required if tag 100=SMART or if ISIN/CUSIP is used
    Tag15 = 15

    // Tag14 CumQty. Total number of shares filled.
    Tag14 = 14

    // Tag16 16 EndSeqNo
    Tag16 = 16

    // Tag17 ExecID
    Tag17 = 17

    // Tag18 ExecInst Used to create Relative, Market Peg, Trailing Stop, and VWAP orders.
    Tag18 = 18

    // Tag19 ExecRefID
    Tag19 = 19

    // Tag20 ExecTransType
    Tag20 = 20

    // Tag21 21 HandlInst
    Tag21 = 21

    // Tag22 IDSource
    Tag22 = 22

    // Tag30 LastMkt. Market of execution for last fill.
    Tag30 = 30

    // Tag31 LastPx. Price of this (last) fill.
    Tag31 = 31

    // Tag32 LastShares
    Tag32 = 32

    // Tag37 OrderID. Unique identifier for Order as assigned by broker. Uniqueness must be guaranteed
    //       within a single trading day. Firms which accept multi-day orders should consider
    //       embedding a date within the OrderID field to assure uniqueness across days.
    Tag37 = 37

    // Tag38 38 OrderQty. Number of shares ordered
    Tag38 = 38

    // Tag39 OrdStatus Identifies the current status of the order
    Tag39 = 39

    // Tag40 40 OrdType
    Tag40 = 40

    // Tag41 41 OrigClOrdID. ClOrdID of the previous order (NOT the initial order of the day) as assigned by the
    //       institution, used to identify the previous order in cancel and cancel/replace requests.
    Tag41 = 41

    // Tag44 44 Price. Price per share. Should not be present for market orders.
    Tag44 = 44

    // Tag47 Rule80A
    Tag47 = 47

    // Tag48 SecurityID
    Tag48 = 48

    // Tag54 54 Side
    Tag54 = 54

    // Tag55 55 Symbol
    Tag55 = 55

    // Tag59 TimeInForce
    Tag59 = 59

    // Tag60 TransactTime
    Tag60 = 60

    // Tag61 Urgency. Urgency of message
    Tag61 = 61

    // Tag63 SettlemntTyp
    Tag63 = 63

    // Tag65 Symbol Suffix
    Tag65 = 65

    // Tag77 OpenClose
    Tag77 = 77

    // Tag99  StopPx. Required for stop and stop limit orders
    Tag99 = 99

    // Tag100 100 ExDestination
    Tag100 = 100

    // Tag103 OrdRejReason
    Tag103 = 103

    // Tag109 ClientID
    Tag109 = 109

    // Tag111 Used to create Hidden or Iceberg Orders
    Tag111 = 111

    // Tag114 LocateReqd. Indicates whether the broker is to locate the stock in conjunction with a short sale order.
    Tag114 = 114

    // Tag126 ExpireTime. Time/Date of order expiration
    Tag126 = 126

    // Tag148 Headline. The text of the Bulletin will be contained in this tag
    Tag148 = 148

    // Tag150 ExecType Identifies the type of Execution Report
    Tag150 = 150

    // Tag151 Leaves Qty The amount of shares open for further execution.
    Tag151 = 151

    // Tag167 SecurityType
    Tag167 = 167

    // Tag168 EffectiveTime. Time the details within the message should take effect
    Tag168 = 168

    // Tag200 MaturityMonthYear. Month and Year of the maturity for SecurityType=FUT or SecurityType=OPT
    Tag200 = 200

    // Tag201 PutOrCall
    Tag201 = 201

    // Tag202 202 Strike Price
    Tag202 = 202

    // Tag204 CustomerOrFirm
    Tag204 = 204

    // Tag205 Strike Price
    Tag205 = 205

    // Tag207 SecurityExchange
    Tag207 = 207

    // Tag211 Peg Difference. Used to create Relative, Market Peg, & Trailing Stop Orders.
    Tag211 = 211

    // Tag231 ContractMultiplier. Specifies the ratio or multiply factor to convert from contracts to shares
    Tag231 = 231

    // Tag378  ExecRestatementReason
    Tag378 = 378

    // Tag388 DiscretionInst. Code to identify the price a DiscretionOffset is related to and should be mathematically added to.
    Tag388 = 388

    // Tag389 DiscretionOffset. Amount (signed) added to the “related to” price specified via DiscretionInst.
    Tag389 = 389

    // Tag432 ExpireDate. Date of order expiration
    Tag432 = 432

    // Tag434 CxlRejResponseTo. dentifies the type of request that a cancel replace request is in response to.
    Tag434 = 434

    // Tag439 ClearingFirm
    Tag439 = 439

    // Tag440 ClearingAccount. Supplemental accounting information forwarded to clearing house/ firm
    Tag440 = 440

    // Tag5700 LocateBroker.
    Tag5700 = 5700

    // Tag6010 Order Reference
    Tag6010 = 6010

    // Tag6013 ComboLegInfo
    Tag6013 = 6013

    // Tag6035 IBKR Local Symbol
    Tag6035 = 6035

    // Tag6037 OptionOrigin
    Tag6037 = 6037

    // Tag6058 TradingClass
    Tag6058 = 6058

    // Tag6086 ShortSaleRule
    Tag6086 = 6086

    // Tag6143 DailyNewID. ID number associated with a particular bulletin
    Tag6143 = 6143

    // Tag6205 ForceOnlyRTH
    Tag6205 = 6205

    // Tag6207 Cust Account. A unique value required for any and all ND omnibus sub-account
    Tag6207 = 6207

    // Tag6273 PipExchanges
    Tag6273 = 6273

    // Tag6257 NoBarriers
    Tag6257 = 6257

    // Tag6258 BarrierPrice
    Tag6258 = 6258

    // Tag6259 BarrierStopPrice
    Tag6259 = 6259

    // Tag6260 BarriertrailingAmt
    Tag6260 = 6260

    // Tag6261 BarrierPriceDelimter
    Tag6261 = 6261

    // Tag6262 BarrierLimitPrice
    Tag6262 = 6262

    // Tag6269 BarrierTrailingAmtUnit
    Tag6269 = 6269
)

const (
    // InvalidHbtInt invalid HeartBtInt
    InvalidHbtInt = "invalid HeartBtInt(108)"
)
