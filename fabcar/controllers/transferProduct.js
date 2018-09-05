
'use strict';
var bcSdk =require("../sdk/invoke")
exports.transfer=(id, jsonobj,quantity,decission,requestid,requestedFrom,newOwner)=>{
          return new Promise((resolve, reject) => {
            var Id= id;
            var Jsonobj= jsonobj
            var Quantity=quantity
            var Requestid=requestid

            var transferObj={
              "id":Id,
              "jsonobj":Jsonobj,
              "quantity":Quantity,
              "decision":decission,
              "requestid":Requestid,
              "requestedFrom":requestedFrom,
              "newOwner":newOwner
            }
          
   bcSdk.transfer(transferObj)
   .then(results =>{
    console.log("result from sdk==========================================>",results)
    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


