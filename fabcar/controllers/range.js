
'use strict';
var bcSdk =require("../sdk/query")
exports.range=(productsOf)=>{
  var array=[];
        return new Promise((resolve, reject) => {
   bcSdk.range()
   .then(results =>{
     if(results.message[0].Record==undefined){
       return resolve({"status":200,"message":"products not created/requested yet"})
     }
       
    for(let i=0;i<results.message.length;i++){
      console.log(results.message[i].Record.Dispatchedto)
      if(results.message[i].Key.length<=3 && results.message[i].Record.Owner==productsOf){
        array.push(results.message[i])
      }else if(results.message[i].Key.length<=3&&results.message[i].Record.Dispatchedto!==null){
       if(results.message[i].Record.Dispatchedto.length==1){
          array.push(results.message[i])
       }else if(results.message[i].Record.Dispatchedto.length==2){
        array.push(results.message[i])
       }else if(results.message[i].Record.Dispatchedto.length==3){
        array.push(results.message[i])
       }
       }else{
         console.log("condition dint match")
       }
       console.log("length of array======================>",array.length) 
    }    
           
    resolve({ "status":results.status, "message": array })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}
exports.allrange=()=>{
        return new Promise((resolve, reject) => {
   bcSdk.range()
   .then(results =>{
     console.log("results======================>",results.message[0].Record.Owner) 

    resolve({ "status":results.status, "message": results.message })
    }).catch(err=>{
        console.log(err)
    })
     
  })


}





