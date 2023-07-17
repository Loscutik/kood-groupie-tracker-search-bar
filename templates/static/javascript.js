const minRange = document.getElementById('minRange');
const maxRange = document.getElementById('maxRange');

minRange.addEventListener('input', updateSlider);
maxRange.addEventListener('input', updateSlider);

function updateSlider() {
  const min = parseInt(minRange.value);
  const max = parseInt(maxRange.value);

  if (min > max) {
    minRange.value = max;
  } else if (max < min) {
    maxRange.value = min;
  }
}

var slider1 = document.getElementById("minRange");
var slider2 = document.getElementById("maxRange");
var output1 = document.getElementById("CreationYearFrom");
var output2 = document.getElementById("CreationYearTo");
output1.innerHTML = slider1.value; // Display the default slider value
slider1.oninput = function() {
  output1.innerHTML = this.value;
}
slider2.oninput = function() {
  output2.innerHTML = this.value;
}

const selectElement = document.getElementById("frday");
const spanElement = document.getElementById("frdate");

selectElement.addEventListener("change", function() {
  const selectedOption = selectElement.options[selectElement.selectedIndex];
  spanElement.textContent = selectedOption.textContent;
});

const selectElement2 = document.getElementById("frmont");
const spanElement2 = document.getElementById("frmonth");

selectElement2.addEventListener("change", function() {
  const selectedOption2 = selectElement2.options[selectElement2.selectedIndex];
  spanElement2.textContent = selectedOption2.textContent;
});

const selectElement3 = document.getElementById("today");
const spanElement3 = document.getElementById("todate");

selectElement3.addEventListener("change", function() {
  const selectedOption3 = selectElement3.options[selectElement3.selectedIndex];
  spanElement3.textContent = selectedOption3.textContent;
});

const selectElement4 = document.getElementById("tomont");
const spanElement4 = document.getElementById("tomonth");

selectElement4.addEventListener("change", function() {
  const selectedOption4 = selectElement4.options[selectElement4.selectedIndex];
  spanElement4.textContent = selectedOption4.textContent;
});

function filter() {
  if (document.getElementById("filtr").style.display === "block") {
      document.getElementById("filtr").style.display = "none";
      window.scrollTo(0, 0);
  } else {
      document.getElementById("filtr").style.display = "block";
      window.scrollTo(0, 0);
  }
}

function rangeValue() {
  const value = document.querySelectorAll(".outYear")
  const input = document.querySelectorAll(".inYear")
  for (let index = 0; index < value.length; index++) {
      value[index].textContent = input[index].value
      input[index].addEventListener("input", (event) => {
          value[index].textContent = event.target.value
      })
  }
}

var slider3 = document.getElementById("minnYear");
var output3 = document.getElementById("fryear");
var slider4 = document.getElementById("maxxYear");
var output4 = document.getElementById("toyear");
slider3.oninput = function() {
  output3.innerHTML = this.value;
}
slider4.oninput = function() {
  output4.innerHTML = this.value;
}

let numb = document.getElementById("frday");
let mont = document.getElementById("frmont");
if (mont.selectedIndex === 2) {
  numb.remove(document.querySelector('#frday')[30])
}


  const searchinput = document.querySelector("#searchbar");
  const dropdwn = document.getElementById("dropdown");
  function checkfocus() {
      document.getElementById("dropdown").style.display = 'block';
  }
  function checkfocusout() {
    // document.getElementById("dropdown").style.display = 'none'; <<<< EDIT TO REMOVE 
  }
  searchinput.addEventListener("input", updateValue);

  let url = 'http://localhost:8080/json';
  let urlloc = 'http://localhost:8080/jsonloc'; // For Locations
  let jsonres = []
  fetch(url)
  .then(response => {
    // console.log(response)
    if (response.ok) {
      return response.json();
    } else {
      throw new Error('Network response was not OK.');
    }
  })
  .then(data => {
    // Process the JSON data
    // console.log(data);
    jsonres.push(data);
  })
  .catch(error => {
    console.error('Error:', error);
  });

  //console.log("data:", jsonres)
  let jsonresloc = []
  fetch(urlloc)
  .then(response2 => {
    if (response2.ok) {
      return response2.json();
    } else {
      throw new Error('Network response was not OK.');
    }
  })
  .then(data2 => {
    jsonresloc.push(data2);
  })
  .catch(error => {
    console.error('Error:', error);
  });

  console.log("data loc:", jsonresloc)
  // Here should be the function which places the JSON elements according to the typed input
  function updateValue(e) { 
    dropdwn.innerHTML = ""
    let drop = e.target.value
    let objvalue = [] // Array of strings
    let objvalueloc = []
    let typed = ''
    let drp = []
    let groupname = ''
    //dropdwn.textContent = e.target.value;
    for (v = 0; v < drop.length; v++) {
      typed = typed.concat(... drop[v]) 
    }
    // console.log("Typed", typed) /// <<<<< TYPED here is the current text in the search field
    for (i = 0; i < jsonres.length; i++) {
      for (j = 0; j < jsonres[i].length; j++) {
        //  jsonres[i][j] - one index group data, object with key-value pairs
        objvalue = Object.values(jsonres[i][j])
        // console.log(Object.keys(groupobj)) // ['id', 'image', 'name', 'members', 'creationDate', 'firstAlbum', 'locations', 'concertDates', 'relations']
        // console.log(Object.values(groupobj)) 
        for (f = 0; f < objvalue.length; f++) {
          let ob = objvalue[f].toString().toLowerCase()
          if ((ob.includes(typed.toLowerCase())) && ((Object.keys(jsonres[i][j])[f] == 'name') || (Object.keys(jsonres[i][j])[f] == 'members') || (Object.keys(jsonres[i][j])[f] == 'creationDate') || (Object.keys(jsonres[i][j])[f] == 'firstAlbum'))) {
            switch (Object.keys(jsonres[i][j])[f]) {
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
                drp.push(objvalue[f] + ' - ' + Object.keys(jsonres[i][j])[f]); 
                break
            }
          }
        }
      }
    }
    // Locations
    for (let k = 0; k < jsonresloc.length; k++) {
      //console.log("logloghere1", Object.values(jsonresloc[k])[0]) //l
      objvalueloc = Object.values(jsonresloc[k])[0]
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
                break;
          }
        }
      }
    } // End of locations
     console.log('drp', drp)
    for (let y = 0; y < drp.length; y++) {
      let obj = document.createElement("li");
      obj.setAttribute('class', 'listclass');
      // obj.innerHTML = drp[y];
      if (drp[y].includes(' - Artist')) {
        let ar = " - Artist"
        let drpp = drp[y].split(ar).join('')
        obj.innerHTML = '<a href="http://localhost:8080/?Name=' + drpp + '" style="text-decoration:none;color:black;">' + drp[y] + '</a>';
        dropdwn.appendChild(obj);
      }

      if (drp[y].includes(' - Member')) { // NEED TO DO filter
        obj.innerHTML = drp[y];
        // found members ->> "Phil Collins" -> "Phil Collins" and "Genesis" 
        // console.log(groupname, " - Group name;\n", drp[y]) // = Group name
        dropdwn.appendChild(obj);
      }

      if (drp[y].includes(' - Location')) {
        obj.innerHTML = drp[y];
        let arl = " - Location"
        let drppl = (drp[y].split(arl).join('')).replace(', ', '%2C+')
        obj.innerHTML = '<a href="http://localhost:8080/?Location=' + drppl + '" style="text-decoration:none;color:black;">' + drp[y] + '</a>';
        dropdwn.appendChild(obj);
      }

      if (drp[y].includes(' - Creation year')) {
        let yr = " - Creation year"
        let drppy = drp[y].split(yr).join('')
        obj.innerHTML = '<a href="http://localhost:8080/?CreationDateFrom=' + drppy + '&CreationDateTo=' + drppy + '" style="text-decoration:none;color:black;">' + drp[y] + '</a>';
        dropdwn.appendChild(obj);
      }

      if (drp[y].includes(' - First album')) {
        let yf = " - First album"
        let drppya = (drp[y].split(yf).join('')).split('-').join('')
        obj.innerHTML = '<a href="http://localhost:8080/?FirstAlbumFromDay=' + drppya[0] + drppya[1] + '&FirstAlbumFromMonth=' + drppya[2] + drppya[3] + '&FirstAlbumFromYear=' + drppya[4] + drppya[5] + drppya[6] + drppya[7] + '&FirstAlbumToDay=' + drppya[0] + drppya[1] + '&FirstAlbumToMonth=' + drppya[2] + drppya[3] + '&FirstAlbumToYear=' + drppya[4] + drppya[5] + drppya[6] + drppya[7] + '" style="text-decoration:none;color:black;">' + drp[y] + '</a>';
        dropdwn.appendChild(obj);
      }
    }
    if (drp.length == 396) {
      dropdwn.textContent = ''
    }
  }

// NEED TO DO:
// - Member
// - Remove doubles!!!