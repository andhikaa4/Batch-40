// document.write("hello world")
// return alert("hello")
// var, let, const

// type data
//let nama = "Andhika"
//let umur = 22

//console.log(`nama saya ${nama} umur saya ${umur}`);
//console.log("nama saya", nama, "umur saya", umur);
//console.log("nama saya" + " " + nama + " " + "umur saya" + " " + umur);

//let nilai = 60

//if (nilai>70){
    //console.log("A");
//} else if (nilai >= 60)  {
   // console.log("B");
//} else {
//console.log("C");
//}

//bilangan (20, 30)

//function bilangan(bilanganSatu, bilanganDua){

    //result = bilanganSatu + bilanganDua
    //return result;
//}

function submitData() {
    let name = document.getElementById("input-name").value
    let email = document.getElementById("input-email").value
    let phone = document.getElementById("input-pn").value
    let subjek = document.getElementById("input-select").value
    let pesan = document.getElementById("input-message").value

    if (name == ""){
        return alert("Nama Harus diisi")
    }
    if (email == ""){
        return alert("email Harus diisi")
    }
    if (phone == ""){
        return alert("Nomor Telepon Harus diisi")
    }
    if (subjek == ""){
        return alert("Subjek Harus diisi")
    }
    if (pesan == ""){
        return alert("Pesan Harus diisi")
    }

    console.log(name);
    console.log(email);
    console.log(phone);
    console.log(subjek);
    console.log(pesan);

    let a = document.createElement('a')
    a.href=`mailto: ${email}?subject=${subjek} &body=Hallo Nama Saya ${name} Harap Hubungi Saya Secepatnya. ${pesan}`
    a.click()
} 
