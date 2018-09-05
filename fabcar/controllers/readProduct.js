
'use strict';
var bcSdk =require("../sdk/query")
exports.readProduct=(key)=>{
      return new Promise((resolve, reject) => {
   bcSdk.readProduct(key)
   .then(results =>{
    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


