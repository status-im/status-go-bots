console.log("Hellow from static!");
console.log("version", JSON.stringify(web3.version));

var msg = "0xFFAFFA"; 

var from = web3.eth.accounts[0]
try {
    web3.personal.sign(msg, from, function(err, result) {
        console.log("result of that is", result, "error is", err);
    })
} catch(err) {
    console.log("catched an error while running the DApp", err);
}
/*
var text = terms
var msg = ethUtil.bufferToHex(new Buffer(text, 'utf8'))
var from = web3.eth.accounts[0]

var params = [from, msg]
*/


