package ibfixdefine

import "fmt"

/*
values to pb format values
*/

// PbOrdRejReason to proto
func PbOrdRejReason(s string) (ordRejReason int32, ordRejReasonText string, err error) {
    switch s {
    case "0":
        // 0 = Broker Option (default if omitted)
        ordRejReason = 0
        ordRejReasonText = "Broker Option"

    case "3":
        // 3 = Order exceeds limit
        ordRejReason = 3
        ordRejReasonText = "Order exceeds limit"

    default:
        return 0, s, fmt.Errorf("illegal OrdRejReason string:%s", s)
    }

    return ordRejReason, ordRejReasonText, nil
}

/*
PbPutOrCall 201 PutOrCall
0 = Put
1 = Call
*/
func PbPutOrCall(s string) (putOrCall int32, err error) {
    switch s {
    case "0":
        putOrCall = 0

    case "1":
        putOrCall = 1

    default:
        return 0, fmt.Errorf("illegal PutOrCall string:%s", s)
    }

    return putOrCall, nil
}

/*
PbCustomerOrFirm 204 CustomerOrFirm
0 = Customer
1 = Firm
*/
func PbCustomerOrFirm(s string) (customerOrFirm int32, err error) {
    switch s {
    case "0":
        customerOrFirm = 0

    case "1":
        customerOrFirm = 1

    default:
        return 0, fmt.Errorf("illegal CustomerOrFirm string:%s", s)
    }

    return customerOrFirm, nil
}

// PbCxlRejReason to proto
func PbCxlRejReason(s string) (cxlRejReasonCode int32, cxlRejReasonText string, err error) {
    switch s {
    case "0":
        //  0 = Too late to cancel
        cxlRejReasonCode = 0
        cxlRejReasonText = "Too late to cancel"

    case "1":
        // 1 = Unknown order
        cxlRejReasonCode = 1
        cxlRejReasonText = "Unknown order"

    case "2":
        // 2 = Broker Option
        cxlRejReasonCode = 2
        cxlRejReasonText = "Broker Option"

    case "3":
        // 3 = Order already in Pending Cancel or Pending Replace status
        cxlRejReasonCode = 3
        cxlRejReasonText = "Order already in Pending Cancel or Pending Replace status"

    default:
        return -1, s, fmt.Errorf("illegal cxlRejReasonText string:%s", s)
    }

    return cxlRejReasonCode, cxlRejReasonText, nil
}

// PbCxlRejResponseTo to proto
func PbCxlRejResponseTo(s string) (cxlRejResponseToCode int32, cxlRejResponseToText string, err error) {
    switch s {
    case "1":
        // 1 = Order Cancel Request
        cxlRejResponseToCode = 1
        cxlRejResponseToText = " Order Cancel Request"

    case "2":
        // 2 = Order Cancel/Replace Request
        cxlRejResponseToCode = 2
        cxlRejResponseToText = "Order Cancel/Replace Request"

    default:
        return 0, s, fmt.Errorf("illegal CxlRejResponseTo string:%s", s)
    }

    return cxlRejResponseToCode, cxlRejResponseToText, nil
}
