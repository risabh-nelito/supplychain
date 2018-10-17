
'use strict';
var bcSdk =require("../sdk/invoke")
exports.request=(id,requestid,requestedFrom,quantity,requester)=>{
          return new Promise((resolve, reject) => {
            var Quantity=quantity.toString()
            var requestObj={
              "id":id,
              "quantity":Quantity,
              "requestedid":requestid,
              "requestedFrom":requestedFrom,
              "requester":requester
            }
          
   bcSdk.request(requestObj)
   .then(results =>{
    console.log("result from sdk==========================================>",results)
    resolve({ "status":results.status, "message": "request raised successfully your request id is:  "+requestid })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


