//@ts-nocheck
console.log('index.js run! ');

/** @type {string} classes */


/** @type {[Set<string>, Set<string>]} */
let [selectedWorkouts, selectedLocations] = [new Set(), new Set()];
function saveSetToLocalStorage(key, set) {
   localStorage.setItem(key, JSON.stringify([...set]));
}
function getSetFromLocalStorage(key) {
   return new Set(JSON.parse(localStorage.getItem(key) || "[]"));
}

document.addEventListener("DOMContentLoaded", () => {
   console.log(`DOMContentLoaded`);

   selectedWorkouts = getSetFromLocalStorage("selectedWorkouts");
   selectedLocations = getSetFromLocalStorage("selectedLocations");

   let tableEl = document.querySelector("table");
   if (!tableEl) {
      throw new Error();
   }
   let emptyTableMessageEl = document.querySelector("#emptyTableMessage");
   if (!emptyTableMessageEl) {
      throw new Error();
   }

   console.log(`selectedWorkouts: `, selectedWorkouts);
   console.log(`selectedLocations: `, selectedLocations);
   if (selectedWorkouts.size == 0 || selectedLocations.size == 0) {
      tableEl.setAttribute("hidden", true);
      emptyTableMessageEl.removeAttribute("hidden");
   } else {
      tableEl.removeAttribute("hidden");
      emptyTableMessageEl.setAttribute("hidden", true);
   }
   updateTable(tableEl)

   let tableInfoEl = document.querySelector("#tableInfo");
   tableInfoEl.innerHTML = `<i>showing ${selectedLocations.size}/${locations.length} locations and ${selectedWorkouts.size}/${workouts.length} workouts</i>`

});

/** 
 * @param {[]string} classes 
 * @param {HTMLTableElement} tableEl 
 */
function updateTable(tableEl) {
   const tbody = tableEl.tBodies[0];
   tbody.innerHTML = "";

   console.log(`selectedLocations: `, selectedLocations);
   for (const _class of classes) {
      /** @type {[string, string, string, string, string]} */
      let [date, time, workout, instructor, location] = [_class[0], _class[1], _class[2], _class[3], _class[4]];
      if (!selectedWorkouts.has(workout) || !selectedLocations.has(location)) {
         continue;
      }
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