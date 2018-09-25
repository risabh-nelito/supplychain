
'use strict';
var bcSdk =require("../sdk/query")
exports.range=(productsOf)=>{
  var array=[];
        return new Promise((resolve, reject) => {
   bcSdk.range()
   .then(results =>{
    // for(let i=0;i<results.length;i++){
    //   if(results[i].)
    // }    
    console.log("result from sdk==========================================>",results)
    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


