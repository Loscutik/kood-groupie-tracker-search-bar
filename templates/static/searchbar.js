const searchinput = document.querySelector("#searchbar");
const dropdwn = document.getElementById("dropdown");
function checkfocus() {
  document.getElementById("dropdown").style.display = 'block';
}
searchinput.addEventListener("input", updateValue);

let url = 'http://localhost:8080/json';
//  let urlloc = 'http://localhost:8080/jsonloc'; // For Locations
let jsonres = []
fetch(url, { method: 'POST', })
  .then(response => {
    //  console.log(response)
    if (response.ok) {
      return response.json();
    } else {
      throw new Error('Network response was not OK.');
    }
  })
  .then(data => {
    // Process the JSON data
    jsonres = data;
  })
  .catch(error => {
    console.error('Error:', error);
  });

// console.log("data:", jsonres)
//  let jsonresloc = []
//  fetch(urlloc, { method: 'POST', })
//   .then(response2 => {
//     if (response2.ok) {
//       return response2.json();
//     } else {
//       throw new Error('Network response was not OK.');
//     }
//    })
//   .then(data2 => {
//     jsonresloc = data2;
//     console.log("data loc:", jsonresloc)
//    })
//    .catch(error => {
//      console.error('Error:', error);
//    });


// Here should be the function which places the JSON elements according to the typed input
function updateValue(e) {
  dropdwn.innerHTML = ""
  let drop = e.target.value
  let objvalue = [] // Array of strings
 // let objvalueloc = []
  let typed = ''
  let drp = []
  let groupname = ''
  //dropdwn.textContent = e.target.value;
  for (v = 0; v < drop.length; v++) {
    typed = typed.concat(...drop[v])
  }

  //  console.log("Typed", typed) /// <<<<< TYPED here is the current text in the search field
  for (j = 0; j < jsonres.length; j++) {
    //   jsonres[j] - one index group data, object with key-value pairs
    objvalue = Object.values(jsonres[j])
    //  console.log(Object.keys(groupobj)) // ['id', 'image', 'name', 'members', 'creationDate', 'firstAlbum', 'locations', 'concertDates', 'relations']
    //  console.log(Object.values(groupobj)) 
    for (f = 0; f < objvalue.length; f++) {
      let ob = objvalue[f].toString().toLowerCase()
      if (ob.includes(typed.toLowerCase())) { // && ((Object.keys(jsonres[j])[f] == 'name') || (Object.keys(jsonres[j])[f] == 'members') || (Object.keys(jsonres[j])[f] == 'creationDate') || (Object.keys(jsonres[j])[f] == 'firstAlbum'))) {
        switch (Object.keys(jsonres[j])[f]) {
          case 'name':
            drp.push(objvalue[f] + ' - Artist');
            break
          case 'members':
            for (t = 0; t < objvalue[f].length; t++) {
              if (objvalue[f][t].toString().toLowerCase().includes(typed.toLowerCase())) {
                groupname = objvalue[2];
                drp.push(objvalue[f][t] + ' - Member');
              }
            }
            break
          case 'creationDate':
            drp.push(objvalue[f] + ' - Creation year');
            break
          case 'firstAlbum':
            drp.push(objvalue[f] + ' - First album');
            break
          default:
            // drp.push(objvalue[f] + ' - ' + Object.keys(jsonres[j])[f]);
            break
        }
      }

      if (Object.keys(jsonres[j])[f] === 'locations') {
        for (let r = 0; r < (objvalue[f].locations.length); r++) {
          if ((objvalue[f].locations)[r].toString().toLowerCase().includes(typed.toLowerCase())) {
            let obloco = (objvalue[f].locations)[r].replace('-', ', ').replace('_', ' ')
            let arr1 = obloco.split(" ");
            for (var e = 0; e < arr1.length; e++) {
              arr1[e] = arr1[e].charAt(0).toUpperCase() + arr1[e].slice(1);
            }
            let obloco2 = arr1.join(" ");
            drp.push(obloco2 + ' - Location');
          }
        }
      }
    }
  }

  /*
    // Locations
     objvalueloc = Object.values(jsonresloc)[0]
     for (let o = 0; o < (objvalueloc.length); o++) {
      //console.log("logloghere2", objvalueloc[o]) // data of one id
      for (let r = 0; r < (objvalueloc[o].locations.length); r++) {
        if ((objvalueloc[o].locations)[r].toString().toLowerCase().includes(typed.toLowerCase())) {
          let obloco = (objvalueloc[o].locations)[r].replace('-', ', ').replace('_', ' ')
          let arr1 = obloco.split(" ");
          for (var e = 0; e < arr1.length; e++) {
            arr1[e] = arr1[e].charAt(0).toUpperCase() + arr1[e].slice(1);
          }
          let obloco2 = arr1.join(" ");
          drp.push(obloco2 + ' - Location');
          //break;
        }
      }
    } // End of locations
  */
  console.log('drp', drp)
  let drpn = [...new Set(drp)];
  for (let y = 0; y < drpn.length; y++) {
    let obj = document.createElement("li");
    obj.setAttribute('class', 'listclass');
    if (drpn[y].includes(' - Artist')) {
      let ar = " - Artist"
      let drpp = drpn[y].split(ar)[0]
      obj.innerHTML = '<a href="http://localhost:8080/?Name=' + drpp + '" style="text-decoration:none;color:black;">' + drpn[y] + '</a>';
      dropdwn.appendChild(obj);
    }

    if (drpn[y].includes(' - Member')) { 
      obj.innerHTML = drpn[y];
      let arm = " - Member"
      let drppm = drpn[y].split(arm)[0]
      obj.innerHTML = '<a href="http://localhost:8080/?NameOfMember=' + drppm + '" style="text-decoration:none;color:black;">' + drpn[y] + '</a>';
      dropdwn.appendChild(obj);
    }

    if (drpn[y].includes(' - Location')) {
      obj.innerHTML = drpn[y];
      let arl = " - Location"
      let drppl = (drpn[y].split(arl)[0]).replace(', ', '%2C+')
      obj.innerHTML = '<a href="http://localhost:8080/?Location=' + drppl + '" style="text-decoration:none;color:black;">' + drpn[y] + '</a>';
      dropdwn.appendChild(obj);
    }

    if (drpn[y].includes(' - Creation year')) {
      let yr = " - Creation year"
      let drppy = drpn[y].split(yr)[0]
      obj.innerHTML = '<a href="http://localhost:8080/?CreationDateFrom=' + drppy + '&CreationDateTo=' + drppy + '" style="text-decoration:none;color:black;">' + drpn[y] + '</a>';
      dropdwn.appendChild(obj);
    }

    if (drpn[y].includes(' - First album')) {
      let yf = " - First album"
      let drppya = (drpn[y].split(yf)[0]).split('-')
      obj.innerHTML = '<a href="http://localhost:8080/?FirstAlbumFromDay=' + drppya[0] + '&FirstAlbumFromMonth=' + drppya[1] + '&FirstAlbumFromYear=' + drppya[2] + '&FirstAlbumToDay=' + drppya[0] + '&FirstAlbumToMonth=' + drppya[1] + '&FirstAlbumToYear=' + drppya[2] + '" style="text-decoration:none;color:black;">' + drpn[y] + '</a>';
      dropdwn.appendChild(obj);
    }
  }
  if (drp.length == 680) {
    dropdwn.textContent = ''
  }
}
