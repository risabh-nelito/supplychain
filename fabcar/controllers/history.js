
'use strict';
var bcSdk =require("../sdk/query")
exports.history=(id)=>{
          return new Promise((resolve, reject) => {
            var history={
              "id":id
            }
          
   bcSdk.getHistoryForProduct(history)
   .then(results =>{
    console.log("result from sdk==========================================>",results)
    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


