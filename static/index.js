console.log('index.js run!');

// classes
// workouts
// locations

const chosenWorkouts = [...workouts];
const chosenLocations = [...locations];

document.addEventListener("DOMContentLoaded", () => {
  console.log(`DOMContentLoaded`);
  let tableEl = document.querySelector("table");
  if (!tableEl) {
    throw new Error("tableElement not found.");
  }
  updateTable(tableEl)
});


/** 
 * @param {[]string} classes 
 * @param {HTMLTableElement} tableEl 
 */

function updateTable(tableEl) {
  const tbody = tableEl.tBodies[0];
  tbody.innerHTML = "";

  for (const _class of classes) {
    let row = document.createElement("tr");
    for (const datum of _class) {
      let cell = document.createElement("td");
      cell.innerHTML = datum;
      row.appendChild(cell);
    }
    tbody.appendChild(row);
  }
}


function handleInputChange(event) {
  console.log("event: ", event);
}