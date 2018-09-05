
'use strict';
var bcSdk =require("../sdk/invoke")
exports.createUser=(Product)=>{
    return new Promise((resolve, reject) => {
   bcSdk.createProduct(Product)
   .then(results =>{
    console.log("result from sdk==========================================>",results)
    resolve({ "status": 200, "message": "request sent Successfully" })
    })
     
  })


}


