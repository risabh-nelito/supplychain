
'use strict';
var bcSdk =require("../sdk/query")
exports.range=()=>{
        return new Promise((resolve, reject) => {
   bcSdk.range()
   .then(results =>{
    console.log("result from sdk==========================================>",results)
    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


