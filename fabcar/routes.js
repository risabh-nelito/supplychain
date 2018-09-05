'use strict';
const cors = require('cors'); 
var path = require('path'); 
var express = require('express'); 
var createUser=require("./controllers/createUser")
var readProduct=require("./controllers/readProduct")
var range=require("./controllers/range")
var transferProduct= require("./controllers/transferProduct")
var request= require("./controllers/request")
var history= require("./controllers/history")

module.exports = router => {
    
    //test submition service
    router.post('/mocklogin',cors(), (req, res) => {
              var email= req.body.email;
              if(email=="man@supply.com"){
                res.status(200).json({message:"manufacturer"})
              }else if(email=="dist@supply.com"){
               res.status(200).json({message:"distributor"})
              }else if(email=="retail@supply.com"){
               res.status(200).json({message:"retailer"})
              }else if(email=="end@supply.com"){
                res.status(200).json({message:"endUser"})
              }else{
                res.status(400).json({message:"please enter a valid email id"})
    }})

    router.post("/createProduct",cors(),(req,res)=>{
        var max=100
        var min=10
      var random=  Math.floor(Math.random()*(max-min+1)+min);
        var Id= random.toString()
        console.log(Id)
        var jsonobj =req.body.jsonobj
        var str= JSON.stringify(jsonobj)
        var quantity= req.body.quantity
        var Product={"Id":Id,
        "Jsonobj":str,
        "quantity":quantity
        }
   createUser.createUser(Product).then(results=>{
       console.log("results in routes.js",results)
         return res.status(results.status).json({"message":results.message,
        "id":Id})  
   }).catch(err=>{
       console.log(err)
   })

})
router.post("/readProduct",cors(),(req,res)=>{
    var key=req.body.key
    if(!key){
        res.status(401).json({message:"please enter a valid key"})
    }
    readProduct.readProduct(key).then(results=>{
        console.log("results=======>",results)
      res.status(results.status).json({message:results.message})
    }).catch(err=>{
     console.log(err)
    })
})

router.get("/range",cors(),(req,res)=>{
  range.range().then(results=>{
      console.log(results)
      res.status(results.status).json({message:results.message})
  }).catch(err=>{
      console.log(err)
  })
})

router.post("/transferProduct",cors(),(req,res)=>{
    var id, jsonobj,quantity,requestid,decission,requestedFrom,newOwner
    id= req.body.id;
    jsonobj= req.body.jsonObj;
    quantity= req.body.quantity;
    requestid= req.body.requestid;
    decission=req.body.decission;
    requestedFrom= req.body.requestedFrom;
    newOwner=req.body.newOwner;
    
    if(!id|| !jsonobj||!quantity||!requestid||!decission){
        res.status(401).json({ message:"some fields missing" })
    }
    transferProduct.transfer(id, jsonobj,quantity,decission,requestid,requestedFrom,newOwner).then(results=>{
     console.log(results)
     res.status(results.status).json({ message:results.message })
    }).catch(err=>{
        console.log("error in catch block routes.js",err)
    })

})

router.post("/request",cors(),(req,res)=>{
    var id= req.body.id
    var requestid = randomString(5, '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ');
    var requestedFrom = req.body.requestedFrom
    var quantity= req.body.quantity
    var requester= req.body.requester
    request.request(id,requestid,requestedFrom,quantity,requester).then(results=>{
        res.send({status:200,message:results.message})
    }).catch(err=>{
        console.log(err)
    })
})

router.post("/gethistoryforProduct",cors(),(req,res)=>{
    var id= (req.body.key)
    console.log(id)
        history.history(id).then(results=>{
        res.send({status:200,message:result.message})
    }).catch(err=>{
        console.log(err)
    })
})
function randomString(length, chars) {
    var result = '';
    for (var i = length; i > 0; --i) result += chars[Math.floor(Math.random() * chars.length)];
    return result;
}

}
