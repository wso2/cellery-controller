import ballerina/http;

endpoint http:Client client {
    
};

service<http:Service> serviceName {
    newResource (endpoint caller, http:Request request) {
    }
}

// Control plane poller. Pull policies from control plane
// A task will periodically poll the control plane, and get the policies 
// for the different endpoints. How do we uniquely identify endpoints? Unique IDs or names?

// The developer will use this VICKEndpoint to call out, and the resiliency etc policies will be
// enforced by that endpoint.



function planTrip(string lastname, time date) {
    transaction with type = "compensation" t1 {
        string hotelRID = reserveHotel(lastName, date);
        string flightRID = reserveFlight(lastName, date);
    }

    if(hotelRID == "-1" || flightRID == "-1") {
        compensate t1;  
    } else {
        commit t1;
    }
}

// returns reservation ID
@compensation {
    oncompensate = cancelHotel
}
function reserverHotel(@compensatable {paramName: "lName"} string lastName, 
                    @compensatable time date) 
                    returns (@ compensatable {paramName: "hotelRID"} string) {

    transaction with type = "acid" {
        @compensatabe{paramName: "bar"} string foo = randomInt + date;
        chargeCreditCard();
        writeDataToDB();
    }

}

function cancelHotel(string lName, time date, string hotelRID, string bar) {
    transaction with type = "acid" {
        reverseCreditCardCharge();
        undoDBChanges();
    }
}