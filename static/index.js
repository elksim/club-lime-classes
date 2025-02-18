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

   if (selectedWorkouts.size == 0) {
      selectedWorkouts = new Set([...workouts])
      saveSetToLocalStorage("selectedWorkouts", selectedWorkouts);
   }
   if (selectedLocations.size == 0) {
      selectedLocations = new Set([...locations])
      saveSetToLocalStorage("selectedLocations", selectedLocations);
   }
   console.log(`selectedWorkouts: `, selectedWorkouts);
   console.log(`selectedLocations: `, selectedLocations);

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