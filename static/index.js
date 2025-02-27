console.log('index.js run! ');


/** @type {[Set<string>, Set<string>]} *///@ts-ignore
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

   /** @type {HTMLTableElement} */ //@ts-ignore
   let tableEl = document.querySelector("table");
   /** @type {HTMLTableSectionElement} */ //@ts-ignore
   let tbody = tableEl.querySelector("tbody");
   tbody.addEventListener('click', function (event) {
      handleTbodyClick(tbody, event)
   }
   );

   let emptyTableMessageEl = document.querySelector("#emptyTableMessage");
   if (!emptyTableMessageEl) {
      throw new Error();
   }

   if (selectedWorkouts.size == 0 || selectedLocations.size == 0) {
      tableEl.setAttribute("hidden", "");
      emptyTableMessageEl.removeAttribute("hidden");
   } else {
      tableEl.removeAttribute("hidden");
      emptyTableMessageEl.setAttribute("hidden", "");
   }
   updateTable(tableEl)

   /** @type {HTMLDivElement} */ //@ts-ignore
   let tableInfoEl = document.querySelector("#tableInfo");
   console.log(`selectedWorkouts xddd: `, selectedWorkouts);
   //@ts-ignore
   tableInfoEl.innerHTML = `<i>showing ${selectedLocations.size}/${locations.length} locations and ${selectedWorkouts.size}/${workouts.length} workouts</i>`
});


/** 
 * @param {HTMLTableElement} tableEl 
 */
function updateTable(tableEl) {
   const tbody = tableEl.tBodies[0];
   tbody.innerHTML = "";

   let prevDate = undefined;
   let curDateRowIndex = 0;

   //@ts-ignore
   for (const _class of classes) {
      /** @type {[string, string, string, string, string]} */
      let [date, time, workout, instructor, location] = [_class[0], _class[1], _class[2], _class[3], _class[4]];
      if (!selectedWorkouts.has(workout) || !selectedLocations.has(location)) {
         continue;
      }

      if (prevDate != date) {
         createDateRow(tbody, curDateRowIndex, date);
         curDateRowIndex += 1;
         prevDate = date;
      }

      let row = document.createElement("tr");
      for (const datum of [time, workout, instructor, location]) {
         let cell = document.createElement("td");
         let span = document.createElement("span");
         span.innerHTML = datum;
         cell.appendChild(span);
         row.appendChild(cell);
      }
      // row.addEventListener("click", () => {
      //    [highlightedProp, highlightedValue] = [undefined, undefined];
      //    for (const row of tableEl.children) {
      //       row.classList.remove("highlighted");
      //    }
      //    highlightedProp
      // })


      tbody.appendChild(row);
   }


}

function handleInputChange(event) {
   console.log("event: ", event);
}


/** 
 * @param {HTMLTableSectionElement} tbody
 * @param {number} index
 * @param {string} date
 * */
function createDateRow(tbody, index, date) {
   let row = document.createElement("tr");
   row.dataset["index"] = `${index}`;
   row.classList.add("dateRow");

   let rowData = document.createElement("td");
   rowData.colSpan = 100;

   let flexContainer = document.createElement("div");
   flexContainer.style.display = "flex";
   flexContainer.style.justifyContent = "space-between";
   flexContainer.style.width = "100%";

   let leftButton = document.createElement("div");
   leftButton.classList.add("focusYesterday");
   leftButton.style.userSelect = "none";
   leftButton.style.display = "flex";
   leftButton.innerHTML = "<=";
   leftButton.dataset["index"] = `${index}`;
   leftButton.style.border = "1px solid black";
   leftButton.style.padding = "10px";

   let centerDate = document.createElement("div");
   centerDate.style.display = "flex";
   centerDate.classList.add("dateRowCell");
   centerDate.innerHTML = date;
   centerDate.style.border = "1px solid black";

   let rightButton = document.createElement("div");
   rightButton.style.userSelect = "none";
   rightButton.classList.add("focusTomorrow");
   rightButton.style.display = "flex";
   rightButton.innerHTML = "=>";
   rightButton.dataset["index"] = `${index}`;
   rightButton.style.border = "1px solid black";
   rightButton.style.padding = "10px";

   // Append elements
   flexContainer.appendChild(leftButton);
   flexContainer.appendChild(centerDate);
   flexContainer.appendChild(rightButton);

   rowData.appendChild(flexContainer);
   row.appendChild(rowData);

   // Add the row with spacers
   let spacerRow = createSpacerRow();
   tbody.appendChild(spacerRow);
   tbody.appendChild(row);
   spacerRow = createSpacerRow();
   tbody.appendChild(spacerRow);
}

function createSpacerRow(colspan = 4) {
   const spacerRow = document.createElement("tr");
   spacerRow.classList.add("spacerRow");

   const spacerCell = document.createElement("td");
   spacerCell.colSpan = colspan;

   spacerRow.appendChild(spacerCell);
   return spacerRow;
}

/** @param {number} index */
function scrollToDate(index) {
   console.log(`scrollToDate called`);
   let elToScrollTo = document.querySelector(`tr.dateRow[data-index="${index}"]`);
   if (!elToScrollTo) {
      return
   }
   elToScrollTo.scrollIntoView();
}

/** 
 * @param {HTMLTableSectionElement} tbody
 * @param {MouseEvent} event
 */
function handleTbodyClick(tbody, event) {
   if (!event.target) {
      return
   }
   /** @type {HTMLElement} */ //@ts-ignore
   const eventTarget = event.target;

   console.log(`event.target: `, event.target);
   console.log(`typeof event.target.tagName: `, typeof eventTarget.tagName);
   if (eventTarget.tagName === "DIV" &&
      (eventTarget.classList.contains("focusYesterday") || eventTarget.classList.contains("focusTomorrow"))) {
      let indexToFocus = Number(eventTarget.dataset["index"]);
      eventTarget.classList.contains("focusYesterday") ? indexToFocus -= 1 : indexToFocus += 1;
      scrollToDate(indexToFocus)
   } else if (eventTarget.tagName !== "SPAN") {
      console.log(`xdd`);
      for (const row of tbody.children) {
         row.classList.remove("highlighted");
      }
   } else {
      /** @type {HTMLElement} */ //@ts-ignore
      let parentRow = event.target.parentElement.parentElement;
      console.log(`parentRow: `, parentRow);
      /** @type {string} */ //@ts-ignore
      let cellText = event.target.innerHTML;

      let tdIndex = 0;
      for (tdIndex; tdIndex < parentRow.children.length; tdIndex++) {
         let val = parentRow.children[tdIndex].children[0].innerHTML;
         if (val == cellText) {
            break
         }
      }
      console.log(`tdIndex: `, tdIndex);

      // for now clicking on time is the way to clear the highlighting..
      if (tdIndex == 0) {
         for (const row of tbody.children) {
            row.classList.remove("highlighted");
         }
         return;
      }

      for (const row of tbody.children) {
         if (row.classList.contains("dateRow") || row.classList.contains("spacerRow")) {
            continue
         }
         if (row.children[tdIndex].children[0].innerHTML == cellText) {
            row.classList.add("highlighted");
         } else {
            row.classList.remove("highlighted");
         }
      }
   }
}