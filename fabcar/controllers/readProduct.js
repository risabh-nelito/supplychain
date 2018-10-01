
'use strict';
var bcSdk =require("../sdk/query")
exports.readProduct=(key)=>{
  var array=[];
      return new Promise((resolve, reject) => {
   bcSdk.readProduct(key)
   .then(results =>{
     if(results.status==401){
     return resolve({ "status":results.status, "message":"no request raised yet" })
     }
     for(let i=0;i<results.message.Decision.length;i++){
     if(results.message.Decision[i].Status=="initiated"){
          array.push(results.message.Decision[i])
     }else{
       console.log("dint match condition")
     }
    }
    resolve({ "status":results.status, "message": array })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}


